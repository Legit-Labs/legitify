package utils

func MapSlice[T any, U any](slice []T, mapper func(T) U) []U {
	if slice == nil {
		return nil
	}
	mapped := make([]U, 0, len(slice))
	for _, v := range slice {
		mapped = append(mapped, mapper(v))
	}
	return mapped
}

func CastSliceOfInterface[T any](slice []interface{}) []T {
	return MapSlice(slice, func(i interface{}) T {
		return i.(T)
	})
}
