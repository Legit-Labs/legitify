package enrichers

import (
	"github.com/Legit-Labs/legitify/internal/analyzers"
)

type Enricher interface {
	Enrich(data analyzers.AnalyzedData) (enrichment Enrichment, ok bool)
	Parse(data interface{}) (enrichment Enrichment, err error)
}

type Enrichment interface {
	HumanReadable(prepend string, linebreak string) string
}
