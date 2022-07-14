package converter

import (
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/iancoleman/orderedmap"
)

func newByResourceConverter() outputConverter {
	return &byResourceConverter{}
}

type byResourceConverter struct {
}

func (*byResourceConverter) Element(policyInfo scheme.PolicyInfo, violation scheme.Violation) string {
	return violation.CanonicalLink
}
func (*byResourceConverter) NewScheme() *orderedmap.OrderedMap {
	return scheme.NewByResourceScheme()
}

func (c *byResourceConverter) Convert(output scheme.FlattenedScheme) (interface{}, error) {
	return ConvertToGroupBy(c, output)
}
