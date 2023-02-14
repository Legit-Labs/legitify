package formatter_test

import (
	"log"
	"testing"

	"github.com/Legit-Labs/legitify/internal/outputer/formatter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test"
	"github.com/stretchr/testify/require"
)

func TestOutputFormats(t *testing.T) {
	scheme := scheme_test.SchemeSample()

	for _, name := range formatter.OutputFormats() {
		output, err := formatter.Format(name, formatter.DefaultOutputIndent, scheme, true)

		require.Nilf(t, err, "Unexpected error for output format %s: %s", name, err)
		require.NotNil(t, output, "Expecting output for %s", name)

		switch name {
		case formatter.Human:
			log.Printf("Human-Readable output:\n%s", output)
			continue // Cannot test human formatter - by definition not machine readable

		case formatter.Markdown:
			log.Printf("Markdown output:\n%s", output)
			continue // Cannot test markdown formatter - by definition not machine readable

		case formatter.Json:
			// json has dedicated tests
			continue

		default:
			t.Fatalf("unexpected format: %s", name)
		}
	}
}
