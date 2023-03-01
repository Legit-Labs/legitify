package pagination

type Basic[ApiRetT any, RespT any] struct {
	mapper *MappedPager[[]ApiRetT, ApiRetT, RespT]
}

func New[ApiRetT any, RespT any](fn interface{}, opts interface{}, optioner Optioner) *Basic[ApiRetT, RespT] {
	mapper := func(t []ApiRetT) []ApiRetT { return t }
	return &Basic[ApiRetT, RespT]{
		mapper: NewMapper[[]ApiRetT, ApiRetT, RespT](fn, opts, mapper, optioner),
	}
}

func (p *Basic[ApiRetT, RespT]) Async(params ...interface{}) <-chan AsyncResult[ApiRetT, RespT] {
	return p.mapper.Async(params...)
}

func (p *Basic[ApiRetT, RespT]) Sync(params ...interface{}) (SyncResult[ApiRetT, RespT], error) {
	return p.mapper.Sync(params...)
}
