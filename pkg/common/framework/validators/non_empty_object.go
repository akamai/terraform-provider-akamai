package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Object = nonEmptyObjectValidator{}

// nonEmptyObjectValidator rejects a configured object whose attributes are all null (e.g. `{}`). This matters
// for objects where every attribute is individually optional: nothing else would catch an all-null object, but
// it plans as a known, non-null value while a subsequent empty API response may collapse to a null object,
// tripping Terraform's post-apply consistency check.
type nonEmptyObjectValidator struct{}

// Description describes the validation in plain text formatting.
func (v nonEmptyObjectValidator) Description(_ context.Context) string {
	return "object must be omitted entirely or contain at least one field."
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v nonEmptyObjectValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateObject performs the validation.
func (v nonEmptyObjectValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	for _, attr := range req.ConfigValue.Attributes() {
		if !attr.IsNull() {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(req.Path, "Invalid Value",
		fmt.Sprintf("%s must contain at least one field, or be omitted entirely; an empty object is not allowed.", req.Path))
}

// NonEmptyObject returns a validator that rejects a configured object whose attributes are all null (e.g. `{}`),
// while still allowing the object to be omitted (null) or fully unknown.
func NonEmptyObject() validator.Object {
	return nonEmptyObjectValidator{}
}
