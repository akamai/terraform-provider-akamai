package cloudcertificates

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestKeySpecsValidator(t *testing.T) {
	t.Parallel()
	keySpecsValidator := KeySpecsValidator()
	ctx := context.Background()

	tests := []struct {
		name          string
		input         types.Map
		expectError   bool
		errorContains string
	}{
		{
			name: "valid - RSA and ECDSA",
			input: mustKeySpecsMapValue(t, map[string]string{
				"RSA":   "2048",
				"ECDSA": "P-256",
			}),
			expectError: false,
		},
		{
			name: "valid - single RSA with 4096",
			input: mustKeySpecsMapValue(t, map[string]string{
				"RSA": "4096",
			}),
			expectError: false,
		},
		{
			name: "valid - single ECDSA with P-384",
			input: mustKeySpecsMapValue(t, map[string]string{
				"ECDSA": "P-384",
			}),
			expectError: false,
		},
		{
			name: "invalid - unrecognized key_type",
			input: mustKeySpecsMapValue(t, map[string]string{
				"DSA": "2048",
			}),
			expectError:   true,
			errorContains: `key_type "DSA" is not valid; must be one of: RSA, ECDSA.`,
		},
		{
			name: "invalid - RSA key_type with ECDSA key_size",
			input: mustKeySpecsMapValue(t, map[string]string{
				"RSA": "P-256",
			}),
			expectError:   true,
			errorContains: `key_size "P-256" is not valid for key_type "RSA"; must be one of: 2048, 4096.`,
		},
		{
			name: "invalid - ECDSA key_type with RSA key_size",
			input: mustKeySpecsMapValue(t, map[string]string{
				"ECDSA": "2048",
			}),
			expectError:   true,
			errorContains: `key_size "2048" is not valid for key_type "ECDSA"; must be one of: P-256, P-384.`,
		},
		{
			name:        "null value",
			input:       types.MapNull(types.StringType),
			expectError: false,
		},
		{
			name:        "unknown value",
			input:       types.MapUnknown(types.StringType),
			expectError: false,
		},
		{
			// e.g. key_size interpolated from another resource's not-yet-known attribute at plan time.
			name: "valid - unknown key_size within an otherwise known map",
			input: types.MapValueMust(types.StringType, map[string]attr.Value{
				"RSA": types.StringUnknown(),
			}),
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request := validator.MapRequest{
				Path:        path.Root("key_specs"),
				ConfigValue: tc.input,
			}
			var response validator.MapResponse

			keySpecsValidator.ValidateMap(ctx, request, &response)

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

func mustKeySpecsMapValue(t *testing.T, specs map[string]string) types.Map {
	t.Helper()
	m, diags := types.MapValueFrom(context.Background(), types.StringType, specs)
	require.False(t, diags.HasError(), "failed to build key_specs map value: %+v", diags.Errors())
	return m
}
