package jsonutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualIgnoringArrayOrder(t *testing.T) {
	tests := map[string]struct {
		a        string
		b        string
		expected bool
	}{
		"identical flat objects": {
			a:        `{"a":1,"b":2}`,
			b:        `{"a":1,"b":2}`,
			expected: true,
		},
		"object key order differs": {
			a:        `{"a":1,"b":2}`,
			b:        `{"b":2,"a":1}`,
			expected: true,
		},
		"flat array same order": {
			a:        `[1,2,3]`,
			b:        `[1,2,3]`,
			expected: true,
		},
		"flat array different order": {
			a:        `[1,2,3]`,
			b:        `[3,1,2]`,
			expected: true,
		},
		"object with nested array reordered": {
			a:        `{"tags":["b","a","c"]}`,
			b:        `{"tags":["a","b","c"]}`,
			expected: true,
		},
		"deeply nested arrays reordered": {
			a:        `{"x":{"items":[{"id":2},{"id":1}]}}`,
			b:        `{"x":{"items":[{"id":1},{"id":2}]}}`,
			expected: true,
		},
		"top-level array of objects reordered": {
			a:        `[{"id":1,"name":"one"},{"id":2,"name":"two"}]`,
			b:        `[{"name":"two","id":2},{"name":"one","id":1}]`,
			expected: true,
		},
		"duplicate elements matched by multiplicity": {
			a:        `[1,1,2]`,
			b:        `[2,1,1]`,
			expected: true,
		},
		"different multiplicity": {
			a:        `[1,1]`,
			b:        `[1,2]`,
			expected: false,
		},
		"objects with different values": {
			a:        `{"a":1}`,
			b:        `{"a":2}`,
			expected: false,
		},
		"objects with extra key": {
			a:        `{"a":1}`,
			b:        `{"a":1,"b":2}`,
			expected: false,
		},
		"arrays with different values": {
			a:        `[1,2,3]`,
			b:        `[1,2,4]`,
			expected: false,
		},
		"arrays with different lengths": {
			a:        `[1,2]`,
			b:        `[1,2,3]`,
			expected: false,
		},
		"large integer preserved precisely": {
			a:        `{"id":9007199254740993}`,
			b:        `{"id":9007199254740993}`,
			expected: true,
		},
		"large integers differ": {
			a:        `{"id":9007199254740993}`,
			b:        `{"id":9007199254740994}`,
			expected: false,
		},
		"invalid old JSON": {
			a:        `not-json`,
			b:        `{"a":1}`,
			expected: false,
		},
		"invalid new JSON": {
			a:        `{"a":1}`,
			b:        `not-json`,
			expected: false,
		},
		"null equals null": {
			a:        `null`,
			b:        `null`,
			expected: true,
		},
		"null vs object": {
			a:        `null`,
			b:        `{}`,
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, EqualIgnoringArrayOrder(tc.a, tc.b))
		})
	}
}
