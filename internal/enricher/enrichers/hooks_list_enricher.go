package enrichers

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/map_utils"
	"github.com/iancoleman/orderedmap"
	"golang.org/x/net/context"
)

const HooksList = "hooksList"

func NewHooksListEnricher() hooksListEnricher {
	return hooksListEnricher{}
}

type hooksListEnricher struct {
}

func (e hooksListEnricher) Enrich(_ context.Context, data analyzers.AnalyzedData) (Enrichment, bool) {
	result, err := createHooksListEnrichment(data.ExtraData)
	if err != nil {
		log.Printf("failed to enrich hooks list: %v", err)
		return nil, false
	}
	return result, true
}

func (e hooksListEnricher) Parse(data interface{}) (Enrichment, error) {
	return NewGenericListEnrichmentFromInterface(data)
}

func createHooksListEnrichment(extraData interface{}) (GenericListEnrichment, error) {
	asMap, ok := extraData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hookslist extra data")
	}

	result := []orderedmap.OrderedMap{}
	for k := range asMap {
		var hooksEnrichment map[string]string

		err := json.Unmarshal([]byte(k), &hooksEnrichment)
		if err != nil {
			return nil, err
		}

		result = append(result, *map_utils.ToKeySortedMap(hooksEnrichment))
	}

	// order by url to maintain a determenistic order
	sort.Slice(result, func(i, j int) bool {
		urlI := map_utils.UnsafeGet[string](&result[i], "url")
		urlJ := map_utils.UnsafeGet[string](&result[j], "url")
		return strings.Compare(urlI, urlJ) < 0
	})

	return result, nil
}
