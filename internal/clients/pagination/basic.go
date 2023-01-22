package pagination

type Basic[T any, O any, R any] struct {
	mapper *MappedPager[[]T, T, O, R]
}

func New[T any, O any, R any](fn interface{}, opts interface{}, optioner Optioner) *Basic[T, O, R] {
	mapper := func(t []T) []T { return t }
	return &Basic[T, O, R]{
		mapper: NewMapper[[]T, T, O, R](fn, opts, mapper, optioner),
	}
}

func (p *Basic[T, O, R]) Async(params ...interface{}) <-chan Result[T, R] {
	return p.mapper.Async(params...)
}

func (p *Basic[T, O, R]) Sync(params ...interface{}) Result[T, R] {
	return p.mapper.Sync(params...)
}
