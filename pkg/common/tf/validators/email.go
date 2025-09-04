package validators

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// EmailValidator validates that a string is valid email address.
type EmailValidator struct{}

// Description describes the validation in plain text formatting.
func (v EmailValidator) Description(context.Context) string {
	return "The validator ensures that the provided input is a valid email address"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v EmailValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v EmailValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {

	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	email := req.ConfigValue.ValueString()

	_, err := mail.ParseAddress(email)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Email Address",
			fmt.Sprintf("The provided email address '%s' is not valid: %s", email, err.Error()),
		)
	}
}
