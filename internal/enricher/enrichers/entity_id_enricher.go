package enrichers

import (
	"context"
	"strconv"

	"github.com/Legit-Labs/legitify/internal/analyzers"
)

const EntityId = "entityId"

func NewEntityIdEnricher(ctx context.Context) Enricher {
	return &entityIdEnricher{}
}

type entityIdEnricher struct {
}

func (e *entityIdEnricher) Enrich(data analyzers.AnalyzedData) (Enrichment, bool) {
	entityID := data.Entity.ID()
	return NewBasicEnrichment(strconv.FormatInt(entityID, 10)), true
}

func (e *entityIdEnricher) ShouldEnrich(requestedEnricher string) bool {
	return requestedEnricher == e.Name()
}

func (e *entityIdEnricher) Name() string {
	return EntityId
}
