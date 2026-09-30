package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

var testObjectAttrTypes = map[string]attr.Type{
	"foo": types.StringType,
	"bar": types.StringType,
}

func TestNonEmptyObject(t *testing.T) {
	t.Parallel()
	nonEmptyObject := NonEmptyObject()
	ctx := context.Background()

	tests := []struct {
		name          string
		input         types.Object
		expectError   bool
		errorContains string
	}{
		{
			name:        "null value",
			input:       types.ObjectNull(testObjectAttrTypes),
			expectError: false,
		},
		{
			name:        "unknown value",
			input:       types.ObjectUnknown(testObjectAttrTypes),
			expectError: false,
		},
		{
			name: "all attributes null",
			input: types.ObjectValueMust(testObjectAttrTypes, map[string]attr.Value{
				"foo": types.StringNull(),
				"bar": types.StringNull(),
			}),
			expectError:   true,
			errorContains: `test must contain at least one field, or be omitted entirely; an empty object is not allowed.`,
		},
		{
			name: "partially populated - one attribute set",
			input: types.ObjectValueMust(testObjectAttrTypes, map[string]attr.Value{
				"foo": types.StringValue("value"),
				"bar": types.StringNull(),
			}),
			expectError: false,
		},
		{
			name: "fully populated",
			input: types.ObjectValueMust(testObjectAttrTypes, map[string]attr.Value{
				"foo": types.StringValue("value"),
				"bar": types.StringValue("other"),
			}),
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request := validator.ObjectRequest{
				Path:        path.Root("test"),
				ConfigValue: tc.input,
			}
			var response validator.ObjectResponse

			nonEmptyObject.ValidateObject(ctx, request, &response)

			if tc.expectError {
				require.True(t, response.Diagnostics.HasError(), "expected error but got none")
				found := false
				for _, d := range response.Diagnostics.Errors() {
					if d.Detail() == tc.errorContains {
						found = true
						break
					}
				}
				require.True(t, found, "expected error %q, but got: %+v", tc.errorContains, response.Diagnostics.Errors())
			} else {
				require.False(t, response.Diagnostics.HasError(), "expected no error but got: %+v", response.Diagnostics.Errors())
			}
		})
	}
}
