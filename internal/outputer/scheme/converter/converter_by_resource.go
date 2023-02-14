package converter

import (
	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

func newByResourceConverter() outputConverter {
	return &byResourceConverter{}
}

type byResourceConverter struct {
}

func (*byResourceConverter) Element(policyInfo scheme.PolicyInfo, violation scheme.Violation) string {
	return violation.CanonicalLink
}
func (*byResourceConverter) NewScheme() groupingScheme {
	return scheme.NewByResource()
}

func (c *byResourceConverter) Convert(output *scheme.Flattened) (scheme.Scheme, error) {
	converted, err := ConvertToGroupBy(c, output)
	if err != nil {
		return nil, err
	}
	return (*scheme.ByResource)(converted), nil
}
