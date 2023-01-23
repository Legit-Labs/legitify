package outputer

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Legit-Labs/legitify/internal/enricher"
	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test"
	"github.com/stretchr/testify/require"
)

type writerMock struct {
	resultChannel chan<- []byte
}

func (m *writerMock) Write(data []byte) (int, error) {
	m.resultChannel <- data
	close(m.resultChannel)
	return len(data), nil
}

func TestOutputer(t *testing.T) {
	data := scheme_test.EnrichedDataSample()

	inputChannel := make(chan enricher.EnrichedData, len(data))
	outputer := NewOutputer(context.Background(), formatter.Json, scheme.TypeFlattened, false)

	// Setup a channel to get the output from the Writer mock
	resultChannel := make(chan []byte, 1)
	errChannel := make(chan error, 1)
	writerMock := &writerMock{resultChannel}
	waiter := outputer.Digest(inputChannel)

	go func() {
		waiter.Wait()
		err := outputer.Output(writerMock)
		errChannel <- err
	}()
	go func() {
		for _, d := range data {
			inputChannel <- d
		}
		close(inputChannel)
	}()

	output := <-resultChannel
	require.NotNil(t, output, "Expecting output")

	err := <-errChannel
	require.Nil(t, err, "Expecting no error")

	var reversed map[string]interface{}
	err = json.Unmarshal(output, &reversed)
	require.Nilf(t, err, "Error deserializing json: %v", err)
	require.NotEmptyf(t, reversed, "Error deserializing json: %v", err)
}
