package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreventSetUpdate returns a plan modifier that ensures given field cannot be updated after creation.
func PreventSetUpdate() planmodifier.Set {
	return preventSetUpdateModifier{}
}

// preventSetUpdateModifier implements the plan modifier.
type preventSetUpdateModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m preventSetUpdateModifier) Description(_ context.Context) string {
	return "Use if you want to ensure that no update is available for given field."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m preventSetUpdateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifySet implements the plan modification logic.
func (m preventSetUpdateModifier) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.State.Raw.IsNull() {
		return
	}
	if !req.PlanValue.Equal(req.StateValue) {
		resp.Diagnostics.AddError("Update not Supported", fmt.Sprintf("updating field `%s` is not possible", req.Path.String()))
	}
}
