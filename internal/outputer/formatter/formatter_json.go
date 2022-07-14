package formatter

import (
	"encoding/json"
)

type JsonFormatter struct {
	indent string
}

func NewJsonFormatter(indent string) OutputFormatter {
	return &JsonFormatter{indent: indent}
}

func (f *JsonFormatter) Format(scheme interface{}, failedOnly bool) ([]byte, error) {
	bytes, err := json.MarshalIndent(scheme, "", f.indent)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (f *JsonFormatter) IsSchemeSupported(schemeType string) bool {
	return true
}
