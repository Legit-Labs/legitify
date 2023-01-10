package enrichers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/utils"
)

const HooksList = "hooksList"

func NewHooksListEnricher(_ context.Context) Enricher {
	return &hooksListEnricher{}
}

type hooksListEnricher struct {
}

func (e *hooksListEnricher) Enrich(data analyzers.AnalyzedData) (Enrichment, bool) {
	result, err := createHooksListEnrichment(data.ExtraData)
	if err != nil {
		return nil, false
	}
	return result, true
}

func createHooksListEnrichment(extraData interface{}) (Enrichment, error) {
	casted, ok := extraData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hookslist extra data")
	}
	var result []map[string]string

	for k := range casted {
		var hooksEnrichment map[string]string

		err := json.Unmarshal([]byte(k), &hooksEnrichment)
		if err != nil {
			return nil, err
		}

		result = append(result, hooksEnrichment)
	}

	return &GenericListEnrichment{
		GenericEnrichments: result,
	}, nil
}

func (e *hooksListEnricher) Name() string {
	return HooksList
}

type GenericListEnrichment struct {
	GenericEnrichments []map[string]string
}

func (se *GenericListEnrichment) Name() string {
	return HooksList
}

func (se *GenericListEnrichment) HumanReadable(prepend string, linebreak string) string {
	sb := utils.NewPrependedStringBuilder(prepend)

	for i, enrichment := range se.GenericEnrichments {
		first := true
		for k, v := range enrichment {
			if first {
				sb.WriteString(fmt.Sprintf("%d. %s: %s%s", i+1, k, v, linebreak))
				first = false
			} else {
				sb.WriteString(fmt.Sprintf("   %s: %s%s", k, v, linebreak))
			}
		}
	}

	return "\n" + sb.String()
}
