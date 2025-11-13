package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreventObjectUpdate returns a plan modifier that ensures given field cannot be updated after creation.
func PreventObjectUpdate() planmodifier.Object {
	return preventObjectUpdateModifier{}
}

// preventObjectUpdateModifier implements the plan modifier.
type preventObjectUpdateModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m preventObjectUpdateModifier) Description(_ context.Context) string {
	return "Use if you want to ensure that no update is available for given field."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m preventObjectUpdateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyObject implements the plan modification logic.
func (m preventObjectUpdateModifier) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.State.Raw.IsNull() {
		return
	}
	if !req.PlanValue.Equal(req.StateValue) {
		resp.Diagnostics.AddError("Update not Supported", fmt.Sprintf("updating field `%s` is not possible", req.Path.String()))
	}
}
