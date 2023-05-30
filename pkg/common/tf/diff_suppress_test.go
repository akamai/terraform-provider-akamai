package tf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldPrefixSuppress(t *testing.T) {
	tests := map[string]struct {
		old, new, prefix string
		expected         bool
	}{
		"equal strings with prefixes":       {"grp_123", "grp_123", "grp_", true},
		"equal strings, old has prefix":     {"grp_123", "123", "grp_", true},
		"equal strings, neither has prefix": {"123", "123", "grp_", true},
		"empty strings":                     {"", "", "grp", true},
		"empty prefix":                      {"grp_123", "grp_123", "", true},
		"different strings with prefix":     {"grp_123", "grp_234", "grp_", false},
		"different strings without prefix":  {"123", "234", "grp_", false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := FieldPrefixSuppress(test.prefix)("", test.old, test.new, nil)
			assert.Equal(t, test.expected, res)
		})
	}
}
