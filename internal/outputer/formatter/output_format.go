package formatter

import (
	"fmt"

	"github.com/Legit-Labs/legitify/internal/outputer/scheme/converter"
)

type FormatName = string

const (
	Human FormatName = "human"
	Json  FormatName = "json"
	Sarif FormatName = "sarif"
)

type OutputFormatter interface {
	Format(scheme interface{}, failedOnly bool) ([]byte, error)
	IsSchemeSupported(schemeType string) bool
}

const DefaultOutputIndent = "  "

type NewFormatFunc func(indent string) OutputFormatter

var outputFormatters = map[FormatName]NewFormatFunc{
	Human: NewHumanFormatter,
	Json:  NewJsonFormatter,
	Sarif: nil, // TODO pending implementation of Sarif output
}

func ValidateOutputFormat(outputFormat FormatName, schemeType converter.SchemeType) error {
	creator, ok := outputFormatters[outputFormat]
	if !ok {
		return fmt.Errorf("Unsupported output format: %s", outputFormat)
	}

	formatter := creator(DefaultOutputIndent)
	if !formatter.IsSchemeSupported(schemeType) {
		return fmt.Errorf("Scheme Type (%s) does not support output format: %s", schemeType, outputFormat)
	}

	return nil
}

func OutputFormats() []FormatName {
	formatNames := []FormatName{}
	for outputFormat, formatter := range outputFormatters {
		if formatter == nil {
			continue
		}
		formatNames = append(formatNames, outputFormat)
	}

	return formatNames
}

func Format(outputFormat FormatName, outputIndent string, scheme interface{}, failedOnly bool) ([]byte, error) {
	outputFormatterCreator := outputFormatters[outputFormat]
	if outputFormatterCreator == nil {
		return nil, fmt.Errorf("No output generator for %s", outputFormat)
	}

	outputFormatter := outputFormatterCreator(outputIndent)

	output, err := outputFormatter.Format(scheme, failedOnly)
	if err != nil {
		return nil, err
	}

	return output, nil
}

type UnsupportedScheme struct {
	scheme interface{}
}

func (e UnsupportedScheme) Error() string {
	return fmt.Sprintf("Unsupported scheme type: %T", e.scheme)
}
