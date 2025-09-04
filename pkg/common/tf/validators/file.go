package validators

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// fileReadableValidator validates that a file is readable at provided file path
type fileReadableValidator struct{}

// Description describes the validation in plain text formatting.
func (v fileReadableValidator) Description(_ context.Context) string {
	return "The validator ensures that the provided file path refers to a readable file"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v fileReadableValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v fileReadableValidator) ValidateString(_ context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() || request.ConfigValue.ValueString() == "" {
		return
	}

	path := request.ConfigValue.ValueString()

	f, err := os.Open(path)

	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		response.Diagnostics.AddAttributeError(request.Path, "The provided path does not refer to a readable file", err.Error())
		return
	}
}

// FileReadable returns an instance of fileReadableValidator
func FileReadable() validator.String {
	return fileReadableValidator{}
}
