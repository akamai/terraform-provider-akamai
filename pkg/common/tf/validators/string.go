package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = notEmptyStringValidator{}

// notEmptyValidator validates that a string Attribute's length is at least a certain value.
type notEmptyStringValidator struct{}

// Description describes the validation in plain text formatting.
func (v notEmptyStringValidator) Description(_ context.Context) string {
	return "Attribute cannot be empty"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v notEmptyStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v notEmptyStringValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	if value != "" {
		return
	}

	attr, _ := request.PathExpression.Steps().LastStep()
	response.Diagnostics.AddAttributeError(request.Path, v.Description(ctx), fmt.Sprintf("Attribute %s cannot be empty", attr))
}

// NotEmptyString returns an validator which ensures that any configured
// attribute value is not an empty string.
func NotEmptyString() validator.String {
	return notEmptyStringValidator{}
}
