package converter

import (
	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/iancoleman/orderedmap"
)

type grouper interface {
	Element(policyInfo scheme.PolicyInfo, violation scheme.Violation) string
	NewScheme() *orderedmap.OrderedMap
}

func ConvertToGroupBy(groupBy grouper, output scheme.FlattenedScheme) (interface{}, error) {
	byElement := groupBy.NewScheme()

	for _, policyName := range output.Keys() {
		outputData := output.GetPolicyData(policyName)
		for _, violation := range outputData.Violations {
			element := groupBy.Element(outputData.PolicyInfo, violation)

			if _, ok := byElement.Get(element); !ok {
				byElement.Set(element, scheme.NewFlattenedScheme())
			}
			byPolicy := utils.UnsafeGet(byElement, element).(scheme.FlattenedScheme)

			if _, ok := byPolicy.Get(policyName); !ok {
				byPolicy.Set(policyName, scheme.NewOutputData(outputData.PolicyInfo))
			}
			preAppend := byPolicy.GetPolicyData(policyName)

			postAppend := scheme.AppendViolations(preAppend, violation)
			byPolicy.Set(policyName, postAppend)
			byElement.Set(element, byPolicy)
		}
	}

	return byElement, nil
}
