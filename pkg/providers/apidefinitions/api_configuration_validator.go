package apidefinitions

import (
	"context"
	"encoding/json"
	"strings"

	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions/v0"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = apiConfigurationValidator{}

// apiConfigurationValidator validates JSON-formatted information about the API configuration.
type apiConfigurationValidator struct{}

// Description describes the validation in plain text formatting.
func (v apiConfigurationValidator) Description(_ context.Context) string {
	return "Invalid JSON-formatted information about the API configuration provided"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v apiConfigurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v apiConfigurationValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
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

// APIConfigurationValidator returns a validator which ensures that JSON-formatted information
// about the API configuration is valid.
func APIConfigurationValidator() validator.String {
	return apiConfigurationValidator{}
}
