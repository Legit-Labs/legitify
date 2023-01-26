package utils

import (
	"log"

	"github.com/iancoleman/orderedmap"
)

// UnsafeGet returns the value associated with the key without checking if the key exists.
// This is useful for loop iterations:
// The loop runs over the keys using Keys(), so it is safe to use Get(),
// Yet the original interface (with comma-ok) does not allow for inline-conversion,
// which forces the creation of an intermediate temporary variable.
func UnsafeGet(m *orderedmap.OrderedMap, key string) interface{} {
	v, _ := m.Get(key)
	return v
}

// Retry is a helper function that retries a function for a given number of times.
func Retry(op func() (shouldRetry bool, err error), max_attempts int, errString string) error {
	var err error
	var shouldRetry bool

	for i := 1; i <= max_attempts; i++ {
		shouldRetry, err = op()
		if err == nil {
			return nil
		}
		if shouldRetry {
			log.Printf("attempt %d/%d failed: %s with err: %s\n", i, max_attempts, errString, err)
		} else {
			log.Printf("failed: %s with err: %s\n", errString, err)
			return err
		}
	}
	log.Printf("all %d attempts failed (%s) with err: %s", max_attempts, errString, err)

	return err
}

func MapSlice[T any, U any](slice []T, mapper func(T) U) []U {
	result := make([]U, 0, len(slice))
	for _, s := range slice {
		result = append(result, mapper(s))
	}
	return result
}
