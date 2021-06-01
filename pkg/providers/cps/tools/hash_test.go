package tools

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tj/assert"
)

func TestHashFromChallengesMap(t *testing.T) {
	tests := map[string]struct {
		val      interface{}
		expected int
	}{
		"valid type passed, return hashed value": {
			val:      map[string]interface{}{"domain": "test"},
			expected: schema.HashString("test"),
		},
		"map does not have 'domain' key": {
			val:      map[string]interface{}{},
			expected: 0,
		},
		"passed value is of invalid type": {
			val:      "test",
			expected: 0,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := HashFromChallengesMap(test.val)
			assert.Equal(t, test.expected, res)
		})
	}
}
