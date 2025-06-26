package nullint64default

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NullInt64 returns a static int64 null value default handler.
//
// Use NullInt64 if a static default null value for an int64 should be set.
func NullInt64() defaults.Int64 {
	return nullInt64{}
}

// nullInt64 is static value default handler that
// sets a value on an int64 attribute.
type nullInt64 struct{}

// Description returns a human-readable description of the default value handler.
func (d nullInt64) Description(_ context.Context) string {
	return "value defaults to null"
}

// MarkdownDescription returns a markdown description of the default value handler.
func (d nullInt64) MarkdownDescription(_ context.Context) string {
	return "value defaults to null"
}

// DefaultInt64 implements the static default value logic.
func (d nullInt64) DefaultInt64(_ context.Context, _ defaults.Int64Request, resp *defaults.Int64Response) {
	resp.PlanValue = types.Int64Null()
}
