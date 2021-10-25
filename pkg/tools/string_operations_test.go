package tools

import (
	"encoding/json"
	"testing"

	"github.com/tj/assert"
)

func TestConvertToString(t *testing.T) {
	type testStr string

	tests := map[string]struct {
		val      interface{}
		expected string
	}{
		"float32": {
			val:      float32(1.23),
			expected: "1.23",
		},
		"float64": {
			val:      float64(1.23),
			expected: "1.23",
		},
		"int": {
			val:      -123,
			expected: "-123",
		},
		"int8": {
			val:      int8(123),
			expected: "123",
		},
		"int16": {
			val:      int16(123),
			expected: "123",
		},
		"int32": {
			val:      int32(123),
			expected: "123",
		},
		"int64": {
			val:      int64(123),
			expected: "123",
		},
		"uint": {
			val:      uint(123),
			expected: "123",
		},
		"json.Number": {
			val:      json.Number("123"),
			expected: "123",
		},
		"string": {
			val:      "123",
			expected: "123",
		},
		"[]byte": {
			val:      []byte("123"),
			expected: "123",
		},
		"boolean": {
			val:      true,
			expected: "true",
		},
		"different type": {
			val:      testStr("123"),
			expected: "123",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var res string
			assert.NotPanics(t, func() {
				res = ConvertToString(test.val)
			})
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestGetFirstNotEmpty(t *testing.T) {
	assert.Equal(t, GetFirstNotEmpty("", "def"), "def")
	assert.Equal(t, GetFirstNotEmpty("val", "def"), "val")
	assert.Equal(t, GetFirstNotEmpty("val", ""), "val")
	assert.Equal(t, GetFirstNotEmpty("", ""), "")
}

func TestContainsString(t *testing.T) {
	assert.False(t, ContainsString([]string{}, "a"))
	assert.True(t, ContainsString([]string{"a"}, "a"))
	assert.True(t, ContainsString([]string{"b", "a"}, "a"))
	assert.False(t, ContainsString([]string{"b", "c"}, "a"))
	assert.False(t, ContainsString([]string{"", "b"}, "a"))
}
