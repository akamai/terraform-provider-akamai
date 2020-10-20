package tools

import (
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"testing"
)

func TestAddPrefix(t *testing.T) {
	tests := map[string]struct {
		givenStr, givenPrefix, expected string
	}{
		"append prefix":                {"test", "pre_", "pre_test"},
		"prefix exists, return string": {"pre_test", "pre_", "pre_test"},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, AddPrefix(test.givenStr, test.givenPrefix))
		})
	}
}

func TestGetIntID(t *testing.T) {
	tests := map[string]struct {
		givenStr, givenPrefix string
		expected              int
		withError             bool
	}{
		"remove prefix and convert to int": {"pre_123", "pre_", 123, false},
		"no prefix, convert to int":        {"123", "pre_", 123, false},
		"invalid string, return error":     {"pre_abc", "pre_", 0, true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := GetIntID(test.givenStr, test.givenPrefix)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}
