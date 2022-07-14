package converter_test

import (
	"testing"

	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test.go"
	"github.com/stretchr/testify/require"
)

func byNamespaceToFlattened(byType scheme.ByTypeScheme) scheme.FlattenedScheme {
	result := scheme.NewFlattenedScheme()

	for _, resourceType := range byType.Keys() {
		subscheme := utils.UnsafeGet(byType, resourceType).(scheme.FlattenedScheme)
		result = scheme_test.CombineSchemes(result, subscheme)
	}

	return result
}

func TestByNamespaceConverter(t *testing.T) {
	sample := scheme_test.SchemeSample()

	output, err := converter.Convert(converter.GroupByNamespace, sample)
	require.Nilf(t, err, "Error converting: %v", err)

	converted := output.(scheme.ByTypeScheme)
	for _, namespace := range converted.Keys() {
		subscheme := utils.UnsafeGet(converted, namespace).(scheme.FlattenedScheme)
		for _, policyName := range subscheme.Keys() {
			outputData := subscheme.GetPolicyData(policyName)
			require.Equalf(t, namespace, outputData.PolicyInfo.Namespace, "Violation namespace mismatch")
		}
	}

	reversed := byNamespaceToFlattened(converted)

	require.Equalf(t, sample, reversed, "Expecting the same result for both directions: %v\n%v\n",
		sample, reversed)
}
