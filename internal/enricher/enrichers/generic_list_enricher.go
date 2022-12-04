package enrichers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/utils"
)

const GenericEnrichmentList = "genericList"

func NewGenericListEnricher(_ context.Context) Enricher {
	return &genericListEnricher{}
}

type genericListEnricher struct {
}

func (e *genericListEnricher) Enrich(data analyzers.AnalyzedData) (Enrichment, bool) {
	result, err := createGenericListEnrichment(data.ExtraData)
	if err != nil {
		return nil, false
	}
	return result, true
}

func createGenericListEnrichment(extraData interface{}) (Enrichment, error) {
	casted, ok := extraData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid membersList extra data")
	}
	var result []map[string]string

	for k := range casted {
		var genericEnrichment map[string]string

		err := json.Unmarshal([]byte(k), &genericEnrichment)
		if err != nil {
			return nil, err
		}

		result = append(result, genericEnrichment)
	}

	return &GenericListEnrichment{
		GenericEnrichments: result,
	}, nil
}

func (e *genericListEnricher) ShouldEnrich(requestedEnricher string) bool {
	return requestedEnricher == e.Name()
}

func (e *genericListEnricher) Name() string {
	return GenericEnrichmentList
}

type GenericListEnrichment struct {
	GenericEnrichments []map[string]string
}

func (se *GenericListEnrichment) HumanReadable(prepend string) string {
	sb := utils.NewPrependedStringBuilder(prepend)

	for i, enrichment := range se.GenericEnrichments {
		first := true
		for k, v := range enrichment {
			if first {
				sb.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, k, v))
				first = false
			} else {
				sb.WriteString(fmt.Sprintf("   %s: %s\n", k, v))
			}
		}
	}

	return sb.String()
}
