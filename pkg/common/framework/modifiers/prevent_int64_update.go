package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreventInt64Update returns a plan modifier that ensures given field cannot be updated after creation.
func PreventInt64Update() planmodifier.Int64 {
	return preventInt64UpdateModifier{}
}

// preventInt64UpdateModifier implements the plan modifier.
type preventInt64UpdateModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m preventInt64UpdateModifier) Description(_ context.Context) string {
	return "Use if you want to ensure that no update is available for given field."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m preventInt64UpdateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyInt64 implements the plan modification logic.
func (m preventInt64UpdateModifier) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.StateValue.IsNull() {
		return
	}
	if req.PlanValue != req.StateValue {
		resp.Diagnostics.AddError("Update not Supported", fmt.Sprintf("updating field `%s` is not possible", req.Path.String()))
	}
}
