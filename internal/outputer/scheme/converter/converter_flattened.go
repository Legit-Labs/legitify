package converter

import "github.com/Legit-Labs/legitify/internal/outputer/scheme"

func newFlattenedConverter() outputConverter {
	return &flattenedConverter{}
}

type flattenedConverter struct {
}

func (c *flattenedConverter) Convert(output *scheme.Flattened) (scheme.Scheme, error) {
	return output, nil
}
