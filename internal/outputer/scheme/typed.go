package scheme

import (
	"encoding/json"
	"fmt"

	"github.com/iancoleman/orderedmap"
)

type TypedScheme[T any] struct {
	Type    SchemeType `json:"type"`
	Content T          `json:"content"`
}

func NewTyped[T any](t SchemeType, content T) *TypedScheme[T] {
	return &TypedScheme[T]{
		Type:    t,
		Content: content,
	}
}

func NewTypedMarshalable(t SchemeType, content Scheme) *TypedScheme[*orderedmap.OrderedMap] {
	return NewTyped(t, content.AsOrderedMap())
}

// Unmarshal unmarshalls a typed json and returns the underlying Flattened scheme
func Unmarshal(data []byte) (*Flattened, error) {
	var typedScheme TypedScheme[json.RawMessage]
	if err := json.Unmarshal(data, &typedScheme); err != nil {
		return nil, fmt.Errorf("failed to parse input: %v", err)
	}

	if typedScheme.Type != TypeFlattened {
		return nil, fmt.Errorf("unmarshaling is only supported for the flattened scheme")
	}

	var flattened Flattened
	if err := json.Unmarshal(typedScheme.Content, &flattened); err != nil {
		return nil, fmt.Errorf("failed to parse flattened scheme: %v", err)
	}

	return &flattened, nil
}
