// Package jsonutil provides JSON comparison utilities shared across SDK v2 and Plugin Framework code.
package jsonutil

import (
	"encoding/json"
	"reflect"
	"strings"
)

// EqualIgnoringArrayOrder reports whether two JSON strings are semantically
// equal, treating arrays as unordered sets at every nesting level.
// Numbers are preserved precisely via json.Number rather than float64.
func EqualIgnoringArrayOrder(a, b string) bool {
	aVal, err := decodeWithNumber(a)
	if err != nil {
		return false
	}
	bVal, err := decodeWithNumber(b)
	if err != nil {
		return false
	}
	return deepEqual(aVal, bVal)
}

func decodeWithNumber(s string) (interface{}, error) {
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func deepEqual(a, b interface{}) bool {
	switch aVal := a.(type) {
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for k, av := range aVal {
			bv, exists := bVal[k]
			if !exists {
				return false
			}
			if !deepEqual(av, bv) {
				return false
			}
		}
		return true
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		return slicesEqual(aVal, bVal)
	default:
		return reflect.DeepEqual(a, b)
	}
}

func slicesEqual(a, b []interface{}) bool {
	used := make([]bool, len(b))
	for _, av := range a {
		found := false
		for j, bv := range b {
			if !used[j] && deepEqual(av, bv) {
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
