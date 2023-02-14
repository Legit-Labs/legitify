package enrichers

import (
	"strconv"

	"github.com/Legit-Labs/legitify/internal/analyzers"
)

const EntityId = "entityId"

type entityIdEnricher struct {
	basicEnricher
}

func NewEntityIdEnricher() entityIdEnricher {
	return entityIdEnricher{
		newBasicEnricher(enrichEntityId),
	}
}

func enrichEntityId(data analyzers.AnalyzedData) (string, bool) {
	entityID := data.Entity.ID()
	return strconv.FormatInt(entityID, 10), true
}
