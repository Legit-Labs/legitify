package formatter_test

import (
	"testing"

	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test"
	"github.com/stretchr/testify/require"
)

func TestFormatCsv(t *testing.T) {
	sample := scheme_test.SchemeSample()

	for _, f := range []bool{true, false} {
		bytes, err := formatter.Format(formatter.Csv, formatter.DefaultOutputIndent, sample, f)
		require.Nilf(t, err, "Error formatting csv: %v", err)
		require.NotNil(t, bytes, "Error formatting csv")
		require.NotEmpty(t, bytes, "Error formatting csv")
	}
}
