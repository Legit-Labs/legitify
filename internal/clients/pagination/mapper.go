package pagination

import (
	"log"
	"reflect"

	"github.com/Legit-Labs/legitify/internal/common/slice_utils"
)

const defaultChannelSize = 1024

type MappedPager[ApiRetT any, UserRetT any, RespT any] struct {
	Opts     interface{}
	Fn       interface{}
	Mapper   func(ApiRetT) []UserRetT
	optioner Optioner
}

func NewMapper[ApiRetT any, UserRetT any, RespT any](fn interface{}, opts interface{}, mapper func(ApiRetT) []UserRetT, optioner Optioner) *MappedPager[ApiRetT, UserRetT, RespT] {
	if fn == nil || mapper == nil {
		log.Panic("creating a pagination mapper requires both a function and a mapper")
	}
	if opts == nil {
		opts = zeroOpts(fn, optioner)
	}
	if reflect.ValueOf(opts).Kind() != reflect.Ptr {
		log.Panic("the options parameter must be of a pointer type")
	}
	return &MappedPager[ApiRetT, UserRetT, RespT]{
		Fn:       fn,
		Opts:     opts,
		Mapper:   mapper,
		optioner: optioner,
	}
}

func (p *MappedPager[ApiRetT, UserRetT, RespT]) Async(params ...interface{}) <-chan AsyncResult[UserRetT, RespT] {
	ch := make(chan AsyncResult[UserRetT, RespT], defaultChannelSize)

	apiCall := p.prepareFunc(params...)
	go func() {
		defer close(ch)
		for {
			result, resp, err := apiCall()

			ch <- newAsyncResult(p.Mapper(result), resp, err)
			if err != nil || p.optioner.Done(resp) {
				break
			}

			p.optioner.Advance(resp, p.Opts)
		}
	}()

	return ch
}

func (p *MappedPager[ApiRetT, UserRetT, RespT]) Sync(params ...interface{}) (SyncResult[UserRetT, RespT], error) {
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

// prepareFunc gets the params for the API call,
// verifies that they match the registered API call function,
// injects the options parameter to the params at the index proided by the optioner,
// and wraps the API call in a simple zero-args func that returns the explicit types.
// It assumes that all API calls return (data, resp, error) as their return values.
func (p *MappedPager[ApiRetT, UserRetT, RespT]) prepareFunc(params ...interface{}) func() (ApiRetT, RespT, error) {
	// params validation
	count, isVariadic := inputsCount(p.Fn)
	paramsCountWithOpts := len(params) + 1
	if isVariadic {
		if paramsCountWithOpts < count-1 {
			log.Panicf("incorrect number of parameters: %d != %d", count, paramsCountWithOpts)
		}
	} else if count != paramsCountWithOpts {
		log.Panicf("incorrect number of parameters: %d != %d", count, paramsCountWithOpts)
	}

	// add the options to the params
	optsIndex := p.optioner.OptionsIndex(count, isVariadic)
	paramsWithOpts := make([]interface{}, 0, paramsCountWithOpts)
	for i, v := range params {
		if i == optsIndex {
			paramsWithOpts = append(params, p.Opts)
		}
		paramsWithOpts = append(paramsWithOpts, v)
	}
	if optsIndex == paramsCountWithOpts-1 { // in case the options is the last arg
		paramsWithOpts = append(paramsWithOpts, p.Opts)
	}
	inputs := slice_utils.Map(paramsWithOpts, reflect.ValueOf)

	return func() (ApiRetT, RespT, error) {
		return p.parseOutputs(reflect.ValueOf(p.Fn).Call(inputs))
	}
}

func (p *MappedPager[ApiRetT, UserRetT, RespT]) parseOutputs(outputs []reflect.Value) (ApiRetT, RespT, error) {
	const (
		dataIndex         = iota
		respIndex         = iota
		errIndex          = iota
		totalReturnValues = iota
	)

	if len(outputs) != totalReturnValues {
		log.Panicf("incorrect number of return values")
	}

	res, ok := outputs[dataIndex].Interface().(ApiRetT)
	if !ok {
		log.Panicf("unexpected result type (%T)", outputs[0].Interface())
	}

	resp, ok := outputs[respIndex].Interface().(RespT)
	if !ok {
		log.Panicf("unexpected response type (%T)", resp)
	}

	errVal := outputs[errIndex].Interface()
	err, ok := outputs[errIndex].Interface().(error)
	if errVal != nil && !ok {
		log.Panicf("unexpected error type")
	}

	return res, resp, err
}

func isVariadic(fn interface{}) bool {
	return reflect.TypeOf(fn).IsVariadic()
}

func inputsCount(fn interface{}) (count int, variadic bool) {
	return reflect.TypeOf(fn).NumIn(), isVariadic(fn)
}

func zeroOpts(fn interface{}, optioner Optioner) interface{} {
	optsLocation := optioner.OptionsIndex(inputsCount(fn))
	return reflect.Zero(reflect.TypeOf(fn).In(optsLocation)).Interface()
}
