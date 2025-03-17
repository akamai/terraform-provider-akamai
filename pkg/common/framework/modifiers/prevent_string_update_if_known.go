package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreventStringUpdateIfKnown returns a plan modifier that ensures given field cannot be updated after creation
// if its planned value is known and not equal to empty string.
func PreventStringUpdateIfKnown(a string) planmodifier.String {
	return preventStringUpdateIfKnownModifier{path: a}
}

// preventStringUpdateIfKnownModifier implements the plan modifier.
type preventStringUpdateIfKnownModifier struct {
	path string
}

// Description returns a human-readable description of the plan modifier.
func (p preventStringUpdateIfKnownModifier) Description(_ context.Context) string {
	return "Use if you want to ensure that no update is available for given field if " +
		"the planned value is known and it's not empty."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (p preventStringUpdateIfKnownModifier) MarkdownDescription(ctx context.Context) string {
	return p.Description(ctx)
}

// PlanModifyString implements the plan modification logic.
func (p preventStringUpdateIfKnownModifier) PlanModifyString(_ context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	if !request.StateValue.IsNull() && !request.PlanValue.IsNull() && request.PlanValue.ValueString() != "" {
		if !request.StateValue.Equal(request.PlanValue) {
			response.Diagnostics.AddError("validation error", fmt.Sprintf("updating '%s' is not allowed", p.path))
			return
		}
	}
}
