package converter_test

import (
	"testing"

	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test"
	"github.com/stretchr/testify/require"
)

func byResourceToByPolicy(byResource *scheme.ByResource) *scheme.Flattened {
	result := scheme.NewFlattenedScheme()

	for _, resourceLink := range byResource.Keys() {
		subscheme := byResource.UnsafeGet(resourceLink)
		result = scheme_test.CombineSchemes(result, subscheme)
	}

	return result
}

func TestByResourceConverter(t *testing.T) {
	sample := scheme_test.SchemeSample()

	output, err := converter.Convert(scheme.TypeGroupByResource, sample)
	require.Nilf(t, err, "Error converting: %v", err)

	converted := output.(*scheme.ByResource)
	for _, resourceLink := range converted.Keys() {
		subscheme := converted.UnsafeGet(resourceLink)
		for _, policyName := range subscheme.AsOrderedMap().Keys() {
			outputData := subscheme.GetPolicyData(policyName)
			for _, violation := range outputData.Violations {
				require.Equalf(t, resourceLink, violation.CanonicalLink, "Violation resource mismatch")
			}
		}
	}

	reversed := byResourceToByPolicy(converted)

	require.Equalf(t, sample, reversed, "Expecting the same result for both directions: %v\n%v\n",
		sample, reversed)
}
