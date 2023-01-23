package enrichers

import (
	"context"

	"github.com/Legit-Labs/legitify/internal/analyzers"
)

const EntityName = "entityName"

type entityNameEnricher struct {
	basicEnricher
}

func NewEntityNameEnricher(ctx context.Context) Enricher {
	return entityNameEnricher{
		newBasicEnricher(enrichEntityName),
	}
}

func enrichEntityName(data analyzers.AnalyzedData) (string, bool) {
	return data.Entity.Name(), true
}
