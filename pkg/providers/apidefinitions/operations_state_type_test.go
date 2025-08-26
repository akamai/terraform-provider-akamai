package apidefinitions

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
)

func TestOperationsStateValidator(t *testing.T) {
	apiValidator := OperationsStateValidator()
	ctx := context.Background()

	tests := []struct {
		name        string
		input       types.String
		expectError bool
		errorRegex  *regexp.Regexp
	}{
		{
			name:        "valid JSON",
			input:       types.StringValue(readJSONFile("resource-operations-01.json")),
			expectError: false,
		},
		{
			name:        "invalid JSON format",
			input:       types.StringValue(readJSONFile("resource-operations-invalid.json")),
			expectError: true,
			errorRegex:  regexp.MustCompile(`unexpected EOF`),
		},
		{
			name:        "unknown field",
			input:       types.StringValue(readJSONFile("resource-operations-unknown-field.json")),
			expectError: true,
			errorRegex:  regexp.MustCompile(`json: unknown field "unknown"`),
		},
		{
			name:        "null value",
			expectError: false,
			input:       basetypes.NewStringNull(),
		},
		{
			name:        "unknown value",
			expectError: false,
			input:       basetypes.NewStringUnknown(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := validator.StringRequest{
				ConfigValue: tc.input,
			}
			var response validator.StringResponse

			apiValidator.ValidateString(ctx, request, &response)

			if tc.expectError {
				require.True(t, response.Diagnostics.HasError(), "expected error but got none")
				if tc.errorRegex != nil {
					found := false
					for _, diag := range response.Diagnostics.Errors() {
						if tc.errorRegex.MatchString(diag.Detail()) {
							found = true
							break
						}
					}
					require.True(t, found, "expected error matching %q, but got: %+v", tc.errorRegex.String(), response.Diagnostics.Errors())
				}
			} else {
				require.False(t, response.Diagnostics.HasError(), "expected no error but got: %+v", response.Diagnostics.Errors())
			}
		})
	}
}
