package enrichers

import (
	"context"
	"fmt"

	"github.com/Legit-Labs/legitify/internal/analyzers"
)

type BasicEnrichment string

func (s BasicEnrichment) HumanReadable(_ string, _ string) string {
	return string(s)
}

func NewBasicEnrichment(val string) BasicEnrichment {
	return BasicEnrichment(val)
}

func NewBasicEnrichmentFromInterface(data interface{}) (BasicEnrichment, error) {
	if val, ok := data.(string); !ok {
		return "", fmt.Errorf("expecting a string, found %t", data)
	} else {
		return BasicEnrichment(val), nil
	}
}

type basicEnricherMethod func(analyzers.AnalyzedData) (string, bool)

type basicEnricher struct {
	EnrichWith basicEnricherMethod
}

func newBasicEnricher(w basicEnricherMethod) basicEnricher {
	return basicEnricher{
		EnrichWith: w,
	}
}

func (e basicEnricher) Enrich(_ context.Context, data analyzers.AnalyzedData) (Enrichment, bool) {
	v, ok := e.EnrichWith(data)
	if !ok {
		return nil, false
	}
	return NewBasicEnrichment(v), true
}

func (e basicEnricher) Parse(data interface{}) (Enrichment, error) {
	return NewBasicEnrichmentFromInterface(data)
}
