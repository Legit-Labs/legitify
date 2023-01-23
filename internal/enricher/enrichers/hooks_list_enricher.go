package enrichers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/iancoleman/orderedmap"
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
		log.Printf("failed to enrich hooks list: %v", err)
		return nil, false
	}
	return result, true
}

func (e *hooksListEnricher) Parse(data interface{}) (Enrichment, error) {
	return NewGenericListEnrichmentFromInterface(data)
}

func createHooksListEnrichment(extraData interface{}) (Enrichment, error) {
	asMap, ok := extraData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hookslist extra data")
	}

	var result []orderedmap.OrderedMap
	for k := range asMap {
		var hooksEnrichment map[string]string

		err := json.Unmarshal([]byte(k), &hooksEnrichment)
		if err != nil {
			return nil, err
		}

		result = append(result, *utils.ToKeySortedMap(hooksEnrichment))
	}

	// order by url to maintain a determenistic order
	sort.Slice(result, func(i, j int) bool {
		urlI := utils.UnsafeGet[string](&result[i], "url")
		urlJ := utils.UnsafeGet[string](&result[j], "url")
		return strings.Compare(urlI, urlJ) < 0
	})

	return GenericListEnrichment(result), nil
}
