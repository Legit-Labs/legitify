package converter

import (
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/iancoleman/orderedmap"
)

func newBySeverityConverter() outputConverter {
	return &bySeverityConverter{}
}

type bySeverityConverter struct {
}

func (*bySeverityConverter) Element(policyInfo scheme.PolicyInfo, violation scheme.Violation) string {
	return policyInfo.Severity
}
func (*bySeverityConverter) NewScheme() *orderedmap.OrderedMap {
	return scheme.NewBySeverityScheme()
}

func (c *bySeverityConverter) Convert(output scheme.FlattenedScheme) (interface{}, error) {
	return ConvertToGroupBy(c, output)
}
