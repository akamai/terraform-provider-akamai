package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreventBoolUpdate returns a plan modifier that ensures given field cannot be updated after creation.
func PreventBoolUpdate() planmodifier.Bool {
	return preventBoolUpdateModifier{}
}

// preventBoolUpdateModifier implements the plan modifier.
type preventBoolUpdateModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m preventBoolUpdateModifier) Description(_ context.Context) string {
	return "Use if you want to ensure that no update is available for given field."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m preventBoolUpdateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyBool implements the plan modification logic.
func (m preventBoolUpdateModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if req.PlanValue != req.StateValue {
		resp.Diagnostics.AddError("Update not Supported", fmt.Sprintf("updating field `%s` is not possible", req.Path.String()))
	}
}
