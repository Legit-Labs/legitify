package converter

import (
	"fmt"

	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

type SchemeType = string

const (
	Flattened        SchemeType = "flattened"
	GroupByNamespace SchemeType = "group-by-namespace"
	GroupByResource  SchemeType = "group-by-resource"
	GroupBySeverity  SchemeType = "group-by-severity"

	DefaultScheme = Flattened
)

func Convert(schemeType SchemeType, output scheme.FlattenedScheme) (interface{}, error) {
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
	Convert(output scheme.FlattenedScheme) (interface{}, error)
}

type newConvertFunc func() outputConverter

var outputConverters = map[SchemeType]newConvertFunc{
	Flattened:        newFlattenedConverter,
	GroupByNamespace: newByNamespaceConverter,
	GroupByResource:  newByResourceConverter,
	GroupBySeverity:  newBySeverityConverter,
}

func ValidateOutputScheme(schemeType SchemeType) error {
	_, ok := outputConverters[schemeType]
	if !ok {
		return fmt.Errorf("unsupported output scheme type: %s", schemeType)
	}

	return nil
}

func SchemeTypes() []SchemeType {
	converterNames := []SchemeType{}
	for outputFormat, formatter := range outputConverters {
		if formatter == nil {
			continue
		}
		converterNames = append(converterNames, outputFormat)
	}

	return converterNames
}
