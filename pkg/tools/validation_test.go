package tools

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tj/assert"
)

func TestIsBlank(t *testing.T) {
	tests := map[string]struct {
		value     interface{}
		withError bool
	}{
		"empty string": {"", true},
		"nil map": {
			value: func() interface{} {
				var x map[string]string
				return x
			}(),
			withError: true,
		},
		"empty map":        {make(map[string]string), true},
		"nil":              {nil, true},
		"non empty string": {"abc", false},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := IsNotBlank(test.value, nil)
			if test.withError {
				assert.NotNil(t, res)
				return
			}
			assert.Empty(t, res)
		})
	}
}

func TestAggregateValidations(t *testing.T) {
	tests := map[string]struct {
		funcs            []schema.SchemaValidateDiagFunc
		expectedDiagsLen int
	}{
		"some functions return errors": {
			funcs: []schema.SchemaValidateDiagFunc{
				func(i interface{}, path cty.Path) diag.Diagnostics {
					return diag.Errorf("1")
				},
				func(i interface{}, path cty.Path) diag.Diagnostics {
					return nil
				},
				func(i interface{}, path cty.Path) diag.Diagnostics {
					return diag.Diagnostics{diag.Diagnostic{Summary: "1"}, diag.Diagnostic{Summary: "2"}}
				},
			},
			expectedDiagsLen: 3,
		},
		"no errors": {
			funcs: []schema.SchemaValidateDiagFunc{
				func(i interface{}, path cty.Path) diag.Diagnostics {
					return nil
				},
			},
			expectedDiagsLen: 0,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := AggregateValidations(test.funcs...)(nil, nil)
			assert.Len(t, res, test.expectedDiagsLen)
		})
	}
}

func TestValidateJSON(t *testing.T) {
	tests := map[string]struct {
		givenVal      interface{}
		expectedError string
	}{
		"valid json passed": {
			givenVal: `{"abc":"cba","number":1}`,
		},
		"passed value is not a string": {
			givenVal:      1,
			expectedError: "value is not a string",
		},
		"invalid json provided": {
			givenVal:      "abc",
			expectedError: "invalid JSON",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := ValidateJSON(test.givenVal, nil)
			if test.expectedError != "" {
				assert.NotEmpty(t, res)
				assert.Contains(t, res[0].Summary, test.expectedError)
				return
			}
			assert.Empty(t, res)
		})
	}
}

func TestEmailValidation(t *testing.T) {
	tests := map[string]struct {
		givenVal      interface{}
		expectedError bool
	}{
		"empty email": {
			givenVal:      "",
			expectedError: true,
		},
		"invalid email": {
			givenVal:      "test",
			expectedError: true,
		},
		"no domain": {
			givenVal:      "test@akamai",
			expectedError: true,
		},
		"valid email": {
			givenVal:      "test@akamai.com",
			expectedError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			diags := ValidateEmail(test.givenVal, nil)
			assert.Equal(t, test.expectedError, diags.HasError())
		})
	}
}
