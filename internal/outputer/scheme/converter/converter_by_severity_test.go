package converter_test

import (
	"testing"

	"github.com/Legit-Labs/legitify/internal/common/utils"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
	"github.com/Legit-Labs/legitify/internal/outputer/scheme/scheme_test.go"
	"github.com/stretchr/testify/require"
)

func bySeverityToByPolicy(bySeverity scheme.BySeverityScheme) scheme.FlattenedScheme {
	result := scheme.NewFlattenedScheme()

	for _, severity := range bySeverity.Keys() {
		subscheme := utils.UnsafeGet(bySeverity, severity).(scheme.FlattenedScheme)
		result = scheme_test.CombineSchemes(result, subscheme)
	}

	return result
}

func TestBySeverityConverter(t *testing.T) {
	sample := scheme_test.SchemeSample()

	output, err := converter.Convert(converter.GroupBySeverity, sample)
	require.Nilf(t, err, "Error converting: %v", err)

	converted := output.(scheme.BySeverityScheme)
	for _, severity := range converted.Keys() {
		subscheme := utils.UnsafeGet(converted, severity).(scheme.FlattenedScheme)
		for _, policyName := range subscheme.Keys() {
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
