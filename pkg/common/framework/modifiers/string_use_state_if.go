package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// StringUseStateIf returns a plan modifier that copies a known prior state
// value into the planned value if given function returns true.
// Use this when you want to suppress value changes.
// Value will remain the same after a resource update.
//
// To prevent Terraform errors, the framework automatically sets unconfigured
// and Computed attributes to an unknown value "(known after apply)" on update.
// Using this plan modifier will instead display the prior state value in the
// plan, unless a prior plan modifier adjusts the value.
func StringUseStateIf(pred func(string, string) bool) planmodifier.String {
	return stringUseStateIfModifier{
		pred: pred,
	}
}

// stringUseStateIfModifier implements the plan modifier.
type stringUseStateIfModifier struct {
	pred func(string, string) bool
}

// Description returns a human-readable description of the plan modifier.
func (m stringUseStateIfModifier) Description(_ context.Context) string {
	return "Use if you want to suppress value changes. " +
		"Once provided function returns true, the value of this attribute in state will not change. "
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m stringUseStateIfModifier) MarkdownDescription(_ context.Context) string {
	return "Use if you want to suppress value changes. " +
		"Once provided function returns true, the value of this attribute in state will not change. "
}

// PlanModifyString implements the plan modification logic.
func (m stringUseStateIfModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	if m.pred(req.PlanValue.ValueString(), req.StateValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}
