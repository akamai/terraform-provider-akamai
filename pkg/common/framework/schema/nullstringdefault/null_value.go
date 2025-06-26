package nullstringdefault

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NullString returns a static null string default handler.
//
// Use NullString if a null value for a string should be set.
func NullString() defaults.String {
	return nullString{}
}

// nullString is static null value default handler.
type nullString struct{}

// Description returns a human-readable description of the default value handler.
func (d nullString) Description(_ context.Context) string {
	return "value defaults to null"
}

// MarkdownDescription returns a markdown description of the default value handler.
func (d nullString) MarkdownDescription(_ context.Context) string {
	return "value defaults to null"
}

// DefaultString implements the static default value logic.
func (d nullString) DefaultString(_ context.Context, _ defaults.StringRequest, resp *defaults.StringResponse) {
	resp.PlanValue = types.StringNull()
}
