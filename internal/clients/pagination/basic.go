package pagination

type Basic[ApiRetT any, OptsT any, RespT any] struct {
	mapper *MappedPager[[]ApiRetT, ApiRetT, OptsT, RespT]
}

func New[ApiRetT any, OptsT any, RespT any](fn interface{}, opts interface{}, optioner Optioner) *Basic[ApiRetT, OptsT, RespT] {
	mapper := func(t []ApiRetT) []ApiRetT { return t }
	return &Basic[ApiRetT, OptsT, RespT]{
		mapper: NewMapper[[]ApiRetT, ApiRetT, OptsT, RespT](fn, opts, mapper, optioner),
	}
}

func (p *Basic[ApiRetT, OptsT, RespT]) Async(params ...interface{}) <-chan AsyncResult[ApiRetT, RespT] {
	return p.mapper.Async(params...)
}

func (p *Basic[ApiRetT, OptsT, RespT]) Sync(params ...interface{}) (SyncResult[ApiRetT, RespT], error) {
	return p.mapper.Sync(params...)
}
