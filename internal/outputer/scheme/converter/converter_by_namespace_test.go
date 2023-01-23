package converter_test

import (
	"testing"

	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test"
	"github.com/stretchr/testify/require"
)

func byNamespaceToFlattened(byNamespace *scheme.ByNamespace) *scheme.Flattened {
	result := scheme.NewFlattenedScheme()

	for _, namespace := range byNamespace.Keys() {
		subscheme := byNamespace.UnsafeGet(namespace)
		result = scheme_test.CombineSchemes(result, subscheme)
	}

	return result
}

func TestByNamespaceConverter(t *testing.T) {
	sample := scheme_test.SchemeSample()

	output, err := converter.Convert(scheme.TypeGroupByNamespace, sample)
	require.Nilf(t, err, "Error converting: %v", err)

	converted := output.(*scheme.ByNamespace)
	for _, namespace := range converted.Keys() {
		subscheme := converted.UnsafeGet(namespace)
		for _, policyName := range subscheme.AsOrderedMap().Keys() {
			outputData := subscheme.GetPolicyData(policyName)
			require.Equalf(t, namespace, outputData.PolicyInfo.Namespace, "Violation namespace mismatch")
		}
	}

	reversed := byNamespaceToFlattened(converted)

	require.Equalf(t, sample, reversed, "Expecting the same result for both directions: %v\n%v\n",
		sample, reversed)
}
