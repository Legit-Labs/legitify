package formatter

import (
	"encoding/json"

	"github.com/Legit-Labs/legitify/internal/outputer/scheme"
)

type JsonFormatter struct {
}

func NewJsonFormatter() OutputFormatter {
	return &JsonFormatter{}
}

func (f *JsonFormatter) Format(s scheme.Scheme, failedOnly bool) ([]byte, error) {
	schemeType, err := scheme.DetectSchemeType(s)
	if err != nil {
		return nil, err
	}
	typed := scheme.NewTypedMarshalable(schemeType, s)

	bytes, err := json.MarshalIndent(&typed, "", DefaultOutputIndent)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (f *JsonFormatter) IsSchemeSupported(schemeType string) bool {
	return true
}
