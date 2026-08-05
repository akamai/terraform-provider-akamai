package cloudcertificates

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/cloudcertificates"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.Map = keySpecsValidator{}

// validKeySizesByType maps each key type to its key sizes allowed by the API.
var validKeySizesByType = map[string][]string{
	string(cloudcertificates.CryptographicAlgorithmRSA):   {string(cloudcertificates.KeySize2048), string(cloudcertificates.KeySize4096)},
	string(cloudcertificates.CryptographicAlgorithmECDSA): {string(cloudcertificates.KeySizeP256), string(cloudcertificates.KeySizeP384)},
}

// keySpecsValidator is the sole validator for key_specs: it checks that each key_type (the map key) is
// recognized, and that each key_size (the map value) is valid for its key_type. Duplicate key_types are not
// possible to express: key_specs is a map keyed by key_type.
type keySpecsValidator struct{}

// Description describes the validation in plain text formatting.
func (v keySpecsValidator) Description(_ context.Context) string {
	return "key_specs keys must be a recognized key_type (RSA or ECDSA), and each value must be a key_size " +
		"valid for its key_type (RSA: 2048 or 4096; ECDSA: P-256 or P-384)."
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v keySpecsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateMap performs the validation.
func (v keySpecsValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var specs map[string]types.String
	if resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &specs, false)...); resp.Diagnostics.HasError() {
		return
	}

	for keyType, keySizeValue := range specs {
		if keySizeValue.IsNull() || keySizeValue.IsUnknown() {
			continue
		}
		keySize := keySizeValue.ValueString()

		validSizes, ok := validKeySizesByType[keyType]
		if !ok {
			resp.Diagnostics.AddAttributeError(req.Path.AtMapKey(keyType),
				"Invalid Key Type", fmt.Sprintf("key_type %q is not valid; must be one of: RSA, ECDSA.", keyType))
			continue
		}

		if !slices.Contains(validSizes, keySize) {
			resp.Diagnostics.AddAttributeError(req.Path.AtMapKey(keyType),
				"Invalid Key Size For Key Type",
				fmt.Sprintf("key_size %q is not valid for key_type %q; must be one of: %s.",
					keySize, keyType, strings.Join(validSizes, ", ")))
		}
	}
}

// KeySpecsValidator returns a validator ensuring key_specs keys are recognized key types, each with a key_size
// valid for that key_type.
func KeySpecsValidator() validator.Map {
	return keySpecsValidator{}
}
