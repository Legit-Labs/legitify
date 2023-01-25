package pagination

import (
	"log"
	"reflect"
)

const defaultChannelSize = 1000

type MappedPager[ApiRetT any, UserRetT any, OptsT any, RespT any] struct {
	Opts     interface{}
	Fn       interface{}
	Mapper   func(ApiRetT) []UserRetT
	optioner Optioner
}

func NewMapper[ApiRetT any, UserRetT any, OptsT any, RespT any](fn interface{}, opts interface{}, mapper func(ApiRetT) []UserRetT, optioner Optioner) *MappedPager[ApiRetT, UserRetT, OptsT, RespT] {
	if fn == nil || mapper == nil {
		log.Panic("creating a pagination mapper requires both a function and a mapper")
	}
	if opts == nil {
		opts = zeroOpts(fn)
	}
	return &MappedPager[ApiRetT, UserRetT, OptsT, RespT]{
		Fn:       fn,
		Opts:     opts,
		Mapper:   mapper,
		optioner: optioner,
	}
}

func (p *MappedPager[ApiRetT, UserRetT, OptsT, RespT]) Async(params ...interface{}) <-chan AsyncResult[UserRetT, RespT] {
	ch := make(chan AsyncResult[UserRetT, RespT], defaultChannelSize)

	f, inputs := p.prepareFunc(params...)
	go func() {
		defer close(ch)
		for {
			result, resp, err := p.parseOutputs(f.Call(inputs))

			ch <- newAsyncResult(p.Mapper(result), resp, err)
			if err != nil || p.optioner.Done(resp) {
				break
			}

			p.optioner.Advance(resp, p.Opts)
		}
	}()

	return ch
}

func (p *MappedPager[ApiRetT, UserRetT, OptsT, RespT]) Sync(params ...interface{}) (SyncResult[UserRetT, RespT], error) {
	var results []UserRetT
	ch := p.Async(params...)

	for r := range ch {
		if r.Err != nil {
			return r.SyncResult, r.Err
		}
		results = append(results, r.Collected...)
	}
	return newSyncCollection[UserRetT, RespT](results), nil
}

func (p *MappedPager[ApiRetT, UserRetT, OptsT, RespT]) prepareFunc(params ...interface{}) (reflect.Value, []reflect.Value) {
	// XXX: happens to be true for both GH & GL clients as long as we don't use the extra options field
	// if we need to use it, fix here and in zeroOpts()
	params = append(params, p.Opts)

	if inputsCount(p.Fn) != len(params) {
		log.Panicf("incorrect number of parameters: %d != %d", inputsCount(p.Fn), len(params))
	}

	inputs := make([]reflect.Value, 0, len(params))
	for _, in := range params {
		inputs = append(inputs, reflect.ValueOf(in))
	}

	f := reflect.ValueOf(p.Fn)
	return f, inputs
}

func (p *MappedPager[ApiRetT, UserRetT, OptsT, RespT]) parseOutputs(outputs []reflect.Value) (ApiRetT, RespT, error) {
	if len(outputs) != 3 {
		log.Panicf("incorrect number of return values")
	}
	res, ok := outputs[0].Interface().(ApiRetT)
	if !ok {
		log.Panicf("unexpected result type (%T)", outputs[0].Interface())
	}
	resp, ok := outputs[1].Interface().(RespT)
	if !ok {
		log.Panicf("unexpected response type (%T)", resp)
	}
	errVal := outputs[2].Interface()
	err, ok := outputs[2].Interface().(error)
	if errVal != nil && !ok {
		log.Panicf("unexpected error type")
	}

	return res, resp, err
}

func isVariadic(fn interface{}) bool {
	return reflect.TypeOf(fn).IsVariadic()
}

func inputsCount(fn interface{}) int {
	count := reflect.TypeOf(fn).NumIn()
	if isVariadic(fn) {
		count--
	}
	return count
}

func zeroOpts(fn interface{}) interface{} {
	optsLocation := inputsCount(fn) - 1
	return reflect.Zero(reflect.TypeOf(fn).In(optsLocation)).Interface()
}
