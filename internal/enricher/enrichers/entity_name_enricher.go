package enrichers

import (
	"context"

	"github.com/Legit-Labs/legitify/internal/analyzers"
)

const EntityName = "entityName"

func NewEntityNameEnricher(ctx context.Context) Enricher {
	return &entityNameEnricher{}
}

type entityNameEnricher struct {
}

func (e *entityNameEnricher) Enrich(data analyzers.AnalyzedData) (Enrichment, bool) {
	name := data.Entity.Name()
	return NewBasicEnrichment(name), true
}

func (e *entityNameEnricher) ShouldEnrich(requestedEnricher string) bool {
	return requestedEnricher == e.Name()
}

func (e *entityNameEnricher) Name() string {
	return EntityName
}
