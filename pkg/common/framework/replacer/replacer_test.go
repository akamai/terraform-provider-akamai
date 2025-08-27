// Package replacer_test ...
package replacer_test

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/modifiers"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/framework/replacer"
	"github.com/stretchr/testify/assert"
)

func TestReplacer(t *testing.T) {
	tests := map[string]struct {
		source       []string
		replacements []string
		expected     []string
		eqFunc       func(string, string) bool
	}{
		"all replaced with prefix": {
			source:       []string{"123", "456", "789"},
			replacements: []string{"prp_123", "prp_456", "prp_789"},
			expected:     []string{"prp_123", "prp_456", "prp_789"},
			eqFunc:       modifiers.EqualUpToPrefixFunc("prp_"),
		},
		"some replaced with prefix": {
			source:       []string{"123", "456", "789"},
			replacements: []string{"prp_123", "prp_456"},
			expected:     []string{"prp_123", "prp_456", "789"},
			eqFunc:       modifiers.EqualUpToPrefixFunc("prp_"),
		},
		"replacements without prefix": {
			source:       []string{"prp_123", "prp_456", "prp_789"},
			replacements: []string{"123", "456"},
			expected:     []string{"123", "456", "prp_789"},
			eqFunc:       modifiers.EqualUpToPrefixFunc("prp_"),
		},
		"no matches": {
			source:       []string{"123", "456", "789"},
			replacements: []string{"prp_1234", "prp_4567", "prp_7890"},
			expected:     []string{"123", "456", "789"},
			eqFunc:       modifiers.EqualUpToPrefixFunc("prp_"),
		},
		"no matches - wrong prefix": {
			source:       []string{"123", "456", "789"},
			replacements: []string{"prp_123", "prp_456"},
			expected:     []string{"123", "456", "789"},
			eqFunc:       modifiers.EqualUpToPrefixFunc("ctr_"),
		},
		"no eq func - should not modify source": {
			source:       []string{"123", "456", "789"},
			replacements: []string{"prp_123", "prp_456"},
			expected:     []string{"123", "456", "789"},
			eqFunc:       nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := replacer.Replacer{
				Source:       test.source,
				Replacements: test.replacements,
				EqFunc:       test.eqFunc,
			}.Replace()

			assert.Equal(t, test.expected, actual)
		})
	}
}
