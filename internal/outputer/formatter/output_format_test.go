package formatter_test

import (
	"log"
	"testing"

	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/formatter/formatter_test"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test.go"
	"github.com/stretchr/testify/require"
)

func TestOutputFormats(t *testing.T) {
	scheme := scheme_test.SchemeSample()
	mapped, err := scheme_test.StructToMap(scheme)
	require.Nilf(t, err, "Error converting struct to map: %v", err)

	for _, name := range formatter.OutputFormats() {
		output, err := formatter.Format(name, formatter.DefaultOutputIndent, scheme, true)

		require.Nilf(t, err, "Unexpected error for output format %s: %s", name, err)
		require.NotNil(t, output, "Expecting output for %s", name)

		var reversed interface{}
		switch name {
		case formatter.Human:
			log.Printf("Human-Readable output:\n%s", output)
			continue // Cannot test human formatter - by definition not machine readable

		case formatter.Markdown:
			log.Printf("Markdown output:\n%s", output)
			continue // Cannot test markdown formatter - by definition not machine readable

		case formatter.Json:
			reversed, err = formatter_test.DeserializeJson(output)
			require.Nilf(t, err, "Error deserializing json: %v", err)
			require.NotNil(t, output, "Error deserializing json")

		default:
			t.Fatalf("unexpected format: %s", name)
		}

		require.Equal(t, mapped, reversed, "Expecting output to be the same as the input")
	}
}
