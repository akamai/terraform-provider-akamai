package apidefinitions

import (
	"context"
	"encoding/json"
	"strings"

	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions/v0"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = apiStateValidator{}

// apiStateValidator validates JSON-formatted information about the API configuration.
type apiStateValidator struct{}

// Description describes the validation in plain text formatting.
func (v apiStateValidator) Description(_ context.Context) string {
	return "Invalid JSON-formatted information about the API configuration provided"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v apiStateValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v apiStateValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	endpoint := v0.APIAttributes{}
	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&endpoint)
	if err != nil {
		response.Diagnostics.AddError("Invalid JSON provided", err.Error())
	}
}

// APIStateValidator returns a validator which ensures that JSON-formatted information
// about the API configuration is valid.
func APIStateValidator() validator.String {
	return apiStateValidator{}
}
