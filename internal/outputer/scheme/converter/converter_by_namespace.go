package converter

import (
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

func newByNamespaceConverter() outputConverter {
	return &byNamespaceConverter{}
}

type byNamespaceConverter struct {
}

func (*byNamespaceConverter) Element(policyInfo scheme.PolicyInfo, violation scheme.Violation) string {
	return policyInfo.Namespace
}
func (*byNamespaceConverter) NewScheme() groupingScheme {
	return scheme.NewByNamespace()
}

func (c *byNamespaceConverter) Convert(output *scheme.Flattened) (scheme.Scheme, error) {
	converted, err := ConvertToGroupBy(c, output)
	if err != nil {
		return nil, err
	}
	return (*scheme.ByNamespace)(converted), nil
}
