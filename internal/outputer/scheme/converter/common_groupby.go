package converter

import (
	"github.com/Legit-Labs/legitify/internal/common/map_utils"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/iancoleman/orderedmap"
)

type grouper interface {
	Element(policyInfo scheme.PolicyInfo, violation scheme.Violation) string
	NewScheme() groupingScheme
}

type groupingScheme interface {
	AsOrderedMap() *orderedmap.OrderedMap
}

func ConvertToGroupBy(groupBy grouper, output *scheme.Flattened) (*orderedmap.OrderedMap, error) {
	byElement := groupBy.NewScheme().AsOrderedMap()

	for _, policyName := range output.AsOrderedMap().Keys() {
		outputData := output.GetPolicyData(policyName)
		for _, violation := range outputData.Violations {
			element := groupBy.Element(outputData.PolicyInfo, violation)

			if _, ok := byElement.Get(element); !ok {
				byElement.Set(element, scheme.NewFlattenedScheme())
			}
			byPolicy := map_utils.UnsafeGet[*scheme.Flattened](byElement, element)

			if _, ok := byPolicy.AsOrderedMap().Get(policyName); !ok {
				byPolicy.AsOrderedMap().Set(policyName, scheme.NewOutputData(outputData.PolicyInfo))
			}
			preAppend := byPolicy.GetPolicyData(policyName)

			postAppend := scheme.AppendViolations(preAppend, violation)
			byPolicy.AsOrderedMap().Set(policyName, postAppend)
			byElement.Set(element, byPolicy)
		}
	}

	return byElement, nil
}
