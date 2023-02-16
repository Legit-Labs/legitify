package converter

import (
	"fmt"

	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

func Convert(schemeType scheme.SchemeType, output *scheme.Flattened) (scheme.Scheme, error) {
	outputConverterCreator := outputConverters[schemeType]
	if outputConverterCreator == nil {
		return nil, fmt.Errorf("no output converter for %s", schemeType)
	}

	outputConverter := outputConverterCreator()

	converted, err := outputConverter.Convert(output)
	if err != nil {
		return nil, err
	}

	return converted, nil
}

type outputConverter interface {
	Convert(output *scheme.Flattened) (scheme.Scheme, error)
}

type newConvertFunc func() outputConverter

var outputConverters = map[scheme.SchemeType]newConvertFunc{
	scheme.TypeFlattened:        newFlattenedConverter,
	scheme.TypeGroupByNamespace: newByNamespaceConverter,
	scheme.TypeGroupByResource:  newByResourceConverter,
	scheme.TypeGroupBySeverity:  newBySeverityConverter,
}

func ValidateOutputScheme(schemeType scheme.SchemeType) error {
	_, ok := outputConverters[schemeType]
	if !ok {
		return fmt.Errorf("unsupported output scheme type: %s", schemeType)
	}

	return nil
}
