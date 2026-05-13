package str

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddPrefix(t *testing.T) {
	tests := map[string]struct {
		givenStr, givenPrefix, expected string
	}{
		"blank string":                 {"", "pre_", ""},
		"blank prefix":                 {"test", "", "test"},
		"append prefix":                {"test", "pre_", "pre_test"},
		"prefix exists, return string": {"pre_test", "pre_", "pre_test"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, AddPrefix(test.givenStr, test.givenPrefix))
		})
	}
}

var getIDTestCases = []struct {
	name, givenStr, givenPrefix string
	expected                    int64
	withError                   bool
}{
	{"remove prefix and convert", "pre_123", "pre_", 123, false},
	{"no prefix, convert", "123", "pre_", 123, false},
	{"invalid string, return error", "pre_abc", "pre_", 0, true},
}

func TestGetIntID(t *testing.T) {
	for _, test := range getIDTestCases {
		t.Run(test.name, func(t *testing.T) {
			res, err := GetIntID(test.givenStr, test.givenPrefix)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, int(test.expected), res)
		})
	}
}

func TestGetInt64ID(t *testing.T) {
	for _, test := range getIDTestCases {
		t.Run(test.name, func(t *testing.T) {
			res, err := GetInt64ID(test.givenStr, test.givenPrefix)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}
