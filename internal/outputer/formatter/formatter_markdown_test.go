package formatter_test

import (
	"testing"

	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test.go"
	"github.com/stretchr/testify/require"
)

func TestFormatMarkdown(t *testing.T) {
	sample := scheme_test.SchemeSample()

	bytes, err := formatter.Format(formatter.Markdown, formatter.DefaultOutputIndent, sample, true)
	require.Nilf(t, err, "Error formatting markdown: %v", err)
	require.NotNil(t, bytes, "Error formatting markdown")
	require.NotEmpty(t, bytes, "Error formatting markdown")
}
