package converter_test

import (
	"testing"

	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test"
	"github.com/stretchr/testify/require"
)

func bySeverityToByPolicy(bySeverity *scheme.BySeverity) *scheme.Flattened {
	result := scheme.NewFlattenedScheme()

	for _, severity := range bySeverity.AsOrderedMap().Keys() {
		subscheme := bySeverity.UnsafeGet(severity)
		result = scheme_test.CombineSchemes(result, subscheme)
	}

	return result
}

func TestBySeverityConverter(t *testing.T) {
	sample := scheme_test.SchemeSample()

	output, err := converter.Convert(scheme.TypeGroupBySeverity, sample)
	require.Nilf(t, err, "Error converting: %v", err)

	converted := output.(*scheme.BySeverity)
	for _, severity := range converted.AsOrderedMap().Keys() {
		subscheme := converted.UnsafeGet(severity)
		for _, policyName := range subscheme.AsOrderedMap().Keys() {
			outputData := subscheme.GetPolicyData(policyName)
			for range outputData.Violations {
				require.Equalf(t, severity, outputData.PolicyInfo.Severity, "Violation severity mismatch")
			}
		}
	}

	reversed := bySeverityToByPolicy(converted)

	require.Equalf(t, sample, reversed, "Expecting the same result for both directions: %v\n%v\n",
		sample, reversed)
}
