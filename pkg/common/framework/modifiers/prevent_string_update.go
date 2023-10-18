package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreventStringUpdate returns a plan modifier that ensures given field cannot be updated after creation.
func PreventStringUpdate() planmodifier.String {
	return preventStringUpdateModifier{}
}

// preventStringUpdateModifier implements the plan modifier.
type preventStringUpdateModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m preventStringUpdateModifier) Description(_ context.Context) string {
	return "Use if you want to ensure that no update is available for given field."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m preventStringUpdateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyString implements the plan modification logic.
func (m preventStringUpdateModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if req.PlanValue != req.StateValue {
		resp.Diagnostics.AddError("Update not Supported", fmt.Sprintf("updating field `%s` is not possible", req.Path.String()))
	}
}
