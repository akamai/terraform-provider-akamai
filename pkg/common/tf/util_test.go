package tf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestStringValueOrNullIfEmpty(t *testing.T) {

	tests := []struct {
		name     string
		s        string
		expected types.String
	}{
		{"empty string", "", types.StringNull()},
		{"non-empty string", "hello", types.StringValue("hello")},
		{"whitespace string", "   ", types.StringValue("   ")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := StringValueOrNullIfEmpty(tc.s)
			assert.False(t, res.IsUnknown())
			assert.True(t, tc.expected.Equal(res))
			assert.True(t, tc.expected == res)
		})
	}
}
