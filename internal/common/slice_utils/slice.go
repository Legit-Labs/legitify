package slice_utils

func Map[T any, U any](slice []T, mapper func(T) U) []U {
	if slice == nil {
		return nil
	}
	mapped := make([]U, 0, len(slice))
	for _, v := range slice {
		mapped = append(mapped, mapper(v))
	}
	return mapped
}

func CastInterfaces[T any](slice []interface{}) []T {
	return Map(slice, func(i interface{}) T {
		return i.(T)
	})
}
