package pagination

type AsyncResult[ApiRetT any, RespT any] struct {
	SyncResult[ApiRetT, RespT]
	Err error
}

func newAsyncResult[ApiRetT any, RespT any](collected []ApiRetT, resp RespT, err error) AsyncResult[ApiRetT, RespT] {
	return AsyncResult[ApiRetT, RespT]{
		SyncResult: newSyncResult(collected, resp),
		Err:        err,
	}
}

type SyncResult[ApiRetT any, RespT any] struct {
	Collected []ApiRetT
	Resp      RespT
}

func newSyncResult[ApiRetT any, RespT any](collected []ApiRetT, resp RespT) SyncResult[ApiRetT, RespT] {
	return SyncResult[ApiRetT, RespT]{
		Collected: collected,
		Resp:      resp,
	}
}

func newSyncCollection[ApiRetT any, RespT any](collected []ApiRetT) SyncResult[ApiRetT, RespT] {
	return SyncResult[ApiRetT, RespT]{
		Collected: collected,
	}
}
