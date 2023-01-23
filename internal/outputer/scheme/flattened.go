package scheme

import (
	"context"
	"encoding/json"

	"github.com/Legit-Labs/legitify/internal/analyzers"
	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/Legit-Labs/legitify/internal/context_utils"
	"github.com/iancoleman/orderedmap"
)

// Flattened maps policy-qualified-names to OutputData
// It is used as the default output, from which the converter selects the desired scheme
type Flattened orderedmap.OrderedMap

func NewFlattenedScheme() *Flattened {
	return ToFlattenedScheme(orderedmap.New())
}
func ToFlattenedScheme(m *orderedmap.OrderedMap) *Flattened {
	return (*Flattened)(m)
}

func (s *Flattened) ShallowClone() *Flattened {
	return ToFlattenedScheme(utils.ShallowCloneOrderedMap(s.AsOrderedMap()))
}

func (s *Flattened) AsOrderedMap() *orderedmap.OrderedMap {
	return (*orderedmap.OrderedMap)(s)
}
func (s *Flattened) GetPolicyData(policyName string) OutputData {
	return utils.UnsafeGet[OutputData](s.AsOrderedMap(), policyName)
}

func (s *Flattened) Sorted(lessFunc func(a *orderedmap.Pair, b *orderedmap.Pair) bool) *Flattened {
	output := s.ShallowClone()
	asMap := output.AsOrderedMap()

	asMap.Sort(lessFunc)
	for _, policyName := range asMap.Keys() {
		outputData := output.GetPolicyData(policyName)
		outputData = sortOutputData(outputData)
		output.AsOrderedMap().Set(policyName, outputData)
	}

	return output
}
func (s *Flattened) SortedBySeverity() *Flattened {
	return s.Sorted(policiesSortBySeverityLess)
}
func (s *Flattened) SortedByNamespace() *Flattened {
	return s.Sorted(policiesSortByNamespaceLess)
}

func (s *Flattened) OnlyFailedViolations() *Flattened {
	return s.FilteredByStatus(analyzers.PolicyFailed)
}

func (s *Flattened) FilteredByStatus(status analyzers.PolicyStatus) *Flattened {
	filter := func(violation Violation) bool {
		return violation.Status == status
	}
	return s.FilterByViolation(filter)
}

type ViolationFilter func(violation Violation) bool

func (s *Flattened) FilterByViolation(filter ViolationFilter) *Flattened {
	filteredScheme := NewFlattenedScheme()

	for _, policyName := range s.AsOrderedMap().Keys() {
		outputData := s.GetPolicyData(policyName)
		filteredViolations := []Violation{}
		any := false
		for _, violation := range outputData.Violations {
			if filter(violation) {
				filteredViolations = append(filteredViolations, violation)
				any = true
			}
		}
		if any {
			outputData.Violations = filteredViolations
			filteredScheme.AsOrderedMap().Set(policyName, outputData)
		}
	}

	return filteredScheme
}

func (s *Flattened) UnmarshalJSON(data []byte) error {
	// enable scorecard to support jsons with scorecard
	ctx := context_utils.NewContextWithScorecard(context.Background(), true, true)
	asMap := s.AsOrderedMap()

	err := json.Unmarshal(data, asMap)
	if err != nil {
		return nil
	}

	for _, policyName := range asMap.Keys() {
		outputDataMap := utils.UnsafeGet[orderedmap.OrderedMap](asMap, policyName)
		outputData, err := NewOutputDataFromMap(ctx, &outputDataMap)
		if err != nil {
			return err
		}
		asMap.Set(policyName, *outputData)
	}

	return nil
}
