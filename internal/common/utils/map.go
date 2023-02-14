package utils

import (
	"encoding/json"
	"log"
	"sort"

	"github.com/iancoleman/orderedmap"
)

// UnsafeGet returns the value associated with the key without checking if the key exists.
// This is useful for loop iterations:
// The loop runs over the keys using Keys(), so it is safe to use Get(),
// Yet the original interface (with comma-ok) does not allow for inline-conversion,
// which forces the creation of an intermediate temporary variable.
func UnsafeGet[T any](m *orderedmap.OrderedMap, key string) T {
	v := UnsafeGetUntyped(m, key)
	converted, ok := v.(T)
	if v != nil && !ok {
		log.Panicf("found wrong type %T instead of %T, for key %v with val %v", v, converted, key, v)
	}
	return converted
}

func UnsafeGetUntyped(m *orderedmap.OrderedMap, key string) interface{} {
	v, ok := m.Get(key)
	if !ok {
		log.Panicf("Used unsafe get with %s, but key does not exist [keys: %s]", key, m.Keys())
	}
	return v
}

func ToOrderedMap[T any](src map[string]T, sortFunc func([]string)) *orderedmap.OrderedMap {
	if src == nil {
		return nil
	}

	m := orderedmap.New()
	for k, v := range src {
		m.Set(k, v)
	}
	if sortFunc != nil {
		m.SortKeys(sortFunc)
	}

	return m
}

func ToKeySortedMap[T any](src map[string]T) *orderedmap.OrderedMap {
	return ToOrderedMap(src, sort.Strings)
}

func ShallowCloneOrderedMap(m *orderedmap.OrderedMap) *orderedmap.OrderedMap {
	clone := orderedmap.New()
	for _, k := range m.Keys() {
		v := UnsafeGetUntyped(m, k)
		clone.Set(k, v)
	}
	return clone
}

func UnorderMap(m *orderedmap.OrderedMap) map[string]interface{} {
	return UnorderMapTypedValues[interface{}](m)
}
func UnorderMapTypedValues[T any](m *orderedmap.OrderedMap) map[string]T {
	newM := make(map[string]T, len(m.Keys()))
	for _, k := range m.Keys() {
		newM[k] = UnsafeGetUntyped(m, k).(T)
	}
	return newM
}

func ShallowUnmarshalMap(m map[string]interface{}, v any) error {
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func ShallowUnmarshalOrderedMap(m *orderedmap.OrderedMap, v any) error {
	return ShallowUnmarshalMap(UnorderMap(m), v)
}
