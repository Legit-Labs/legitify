package enrichers

import (
	"github.com/Legit-Labs/legitify/internal/analyzers"
)

const EntityName = "entityName"

type entityNameEnricher struct {
	basicEnricher
}

func NewEntityNameEnricher() entityNameEnricher {
	return entityNameEnricher{
		newBasicEnricher(enrichEntityName),
	}
}

func enrichEntityName(data analyzers.AnalyzedData) (string, bool) {
	return data.Entity.Name(), true
}
