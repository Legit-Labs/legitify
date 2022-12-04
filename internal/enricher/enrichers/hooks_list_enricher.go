package enrichers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/google/go-github/v44/github"
)

const HooksList = "violatedHooks"

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
		return nil, fmt.Errorf("invalid membersList extra data")
	}

	var result []github.Hook
	for k := range casted {
		var hook github.Hook
		err := json.Unmarshal([]byte(k), &hook)
		if err != nil {
			return nil, err
		}

		result = append(result, hook)
	}

	return &HooksListEnrichment{
		Hooks: result,
	}, nil
}

func (e *hooksListEnricher) ShouldEnrich(requestedEnricher string) bool {
	return requestedEnricher == e.Name()
}

func (e *hooksListEnricher) Name() string {
	return HooksList
}

type HooksListEnrichment struct {
	Hooks []github.Hook
}

func (se *HooksListEnrichment) HumanReadable(prepend string) string {
	sb := utils.NewPrependedStringBuilder(prepend)

	for i, hook := range se.Hooks {
		sb.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, *hook.Name, hook.GetURL()))
	}

	return sb.String()
}
