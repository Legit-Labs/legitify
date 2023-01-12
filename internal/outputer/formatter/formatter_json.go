package formatter

import (
	"encoding/json"
)

type JsonFormatter struct {
}

func NewJsonFormatter() OutputFormatter {
	return &JsonFormatter{}
}

func (f *JsonFormatter) Format(scheme interface{}, failedOnly bool) ([]byte, error) {
	bytes, err := json.MarshalIndent(scheme, "", DefaultOutputIndent)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (f *JsonFormatter) IsSchemeSupported(schemeType string) bool {
	return true
}
