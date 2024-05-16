package enrichers

import (
	"encoding/json"
	"fmt"
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/map_utils"
	"github.com/iancoleman/orderedmap"
	"golang.org/x/net/context"
	"log"
)

const SecretsList = "secretsList"

func NewSecretsListEnricher() secretsListEnricher {
	return secretsListEnricher{}
}

type secretsListEnricher struct {
}

func (e secretsListEnricher) Enrich(_ context.Context, data analyzers.AnalyzedData) (Enrichment, bool) {
	result, err := createSecretListEnrichment(data.ExtraData)
	if err != nil {
		log.Printf("failed to enrich secrets list: %v", err)
		return nil, false
	}
	return result, true
}

func (e secretsListEnricher) Parse(data interface{}) (Enrichment, error) {
	return NewGenericListEnrichmentFromInterface(data)
}

func createSecretListEnrichment(extraData interface{}) (GenericListEnrichment, error) {
	asMap, ok := extraData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hookslist extra data")
	}

	result := []orderedmap.OrderedMap{}
	for k := range asMap {
		var secretsEnrichment map[string]string

		err := json.Unmarshal([]byte(k), &secretsEnrichment)
		if err != nil {
			return nil, err
		}

		result = append(result, *map_utils.ToKeySortedMap(secretsEnrichment))
	}

	return result, nil
}
