package utils

import (
	"log"
)

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
