package pagination

import (
	"log"
	"reflect"
)

type MappedPager[T any, U any, O any, R any] struct {
	Opts     interface{}
	Fn       interface{}
	Mapper   func(T) []U
	optioner Optioner
}

func NewMapper[T any, U any, O any, R any](fn interface{}, opts interface{}, mapper func(T) []U, optioner Optioner) *MappedPager[T, U, O, R] {
	if fn == nil || mapper == nil {
		log.Panic("must provide a function and a mapper")
	}
	if opts == nil {
		opts = zeroOpts(fn)
	}
	return &MappedPager[T, U, O, R]{
		Fn:       fn,
		Opts:     opts,
		Mapper:   mapper,
		optioner: optioner,
	}
}

func (p *MappedPager[T, U, O, R]) Async(params ...interface{}) <-chan Result[U, R] {
	ch := make(chan Result[U, R])

	f, inputs := p.prepareFunc(params...)
	go func() {
		defer close(ch)
		for {
			result, resp, err := p.parseOutputs(f.Call(inputs))
			ch <- Result[U, R]{
				Err:       err,
				Resp:      resp,
				Collected: p.Mapper(result),
			}

			if err != nil || p.optioner.Done(resp) {
				break
			}

			p.optioner.Advance(resp, p.Opts)
		}
	}()

	return ch
}

func (p *MappedPager[T, U, O, R]) Sync(params ...interface{}) Result[U, R] {
	var results []U
	ch := p.Async(params...)

	for r := range ch {
		if r.Err != nil {
			return Result[U, R]{
				Err:       r.Err,
				Resp:      r.Resp,
				Collected: results,
			}
		}
		results = append(results, r.Collected...)
	}
	return Result[U, R]{
		Collected: results,
	}
}

func (p *MappedPager[T, U, O, R]) prepareFunc(params ...interface{}) (reflect.Value, []reflect.Value) {
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

func (p *MappedPager[T, U, O, R]) parseOutputs(outputs []reflect.Value) (T, R, error) {
	if len(outputs) != 3 {
		log.Panicf("incorrect number of return values")
	}
	res, ok := outputs[0].Interface().(T)
	if !ok {
		log.Panicf("unexpected result type (%T)", outputs[0].Interface())
	}
	resp, ok := outputs[1].Interface().(R)
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
