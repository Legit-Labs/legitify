package converter

import (
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

func newBySeverityConverter() outputConverter {
	return &bySeverityConverter{}
}

type bySeverityConverter struct {
}

func (*bySeverityConverter) Element(policyInfo scheme.PolicyInfo, violation scheme.Violation) string {
	return policyInfo.Severity
}
func (*bySeverityConverter) NewScheme() groupingScheme {
	return scheme.NewBySeverity()
}

func (c *bySeverityConverter) Convert(output *scheme.Flattened) (scheme.Scheme, error) {
	converted, err := ConvertToGroupBy(c, output)
	if err != nil {
		return nil, err
	}
	return (*scheme.BySeverity)(converted), nil
}
