package formatter_test

import (
	"testing"

	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/formatter/formatter_test"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test.go"
	"github.com/stretchr/testify/require"
)

func TestFormatJson(t *testing.T) {
	sample := scheme_test.SchemeSample()

	bytes, err := formatter.Format(formatter.Json, formatter.DefaultOutputIndent, sample, true)
	require.Nilf(t, err, "Error formatting json: %v", err)
	require.NotNil(t, bytes, "Error formatting json")

	output, err := formatter_test.DeserializeJson(bytes)
	require.Nilf(t, err, "Error deserializing json: %v", err)
	require.NotNil(t, output, "Error deserializing json")

	mapped, err := scheme_test.StructToMap(sample)
	require.Nilf(t, err, "Error converting struct to map: %v", err)

	require.Equal(t, mapped, output)
}
