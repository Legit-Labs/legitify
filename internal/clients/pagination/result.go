package pagination

type Result[ApiRetT any, RespT any] struct {
	Err       error
	Resp      RespT
	Collected []ApiRetT
}
