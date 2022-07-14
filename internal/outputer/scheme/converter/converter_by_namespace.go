package converter

import (
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
	"github.com/iancoleman/orderedmap"
)

func newByNamespaceConverter() outputConverter {
	return &byNamespaceConverter{}
}

type byNamespaceConverter struct {
}

func (*byNamespaceConverter) Element(policyInfo scheme.PolicyInfo, violation scheme.Violation) string {
	return policyInfo.Namespace
}
func (*byNamespaceConverter) NewScheme() *orderedmap.OrderedMap {
	return scheme.NewByTypeScheme()
}

func (c *byNamespaceConverter) Convert(output scheme.FlattenedScheme) (interface{}, error) {
	return ConvertToGroupBy(c, output)
}
