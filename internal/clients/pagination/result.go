package pagination

type Result[T any, R any] struct {
	Err       error
	Resp      R
	Collected []T
}
