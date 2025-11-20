package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// Int64KeepNonNullState returns a plan modifier that keeps the prior state value.
func Int64KeepNonNullState() planmodifier.Int64 {
	return int64KeepNonNullState{}
}

// More context about the modifier implementation and why it's needed:
// https://github.com/hashicorp/terraform-plugin-framework/issues/1117
// https://github.com/hashicorp/terraform-plugin-framework/issues/1197
// https://github.com/hashicorp/terraform-plugin-framework/issues/1211
type int64KeepNonNullState struct{}

// Description returns a human-readable description of the plan modifier.
func (m int64KeepNonNullState) Description(_ context.Context) string {
	return "Once set to a non-null value, the value of this attribute in state will not change."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m int64KeepNonNullState) MarkdownDescription(_ context.Context) string {
	return "Once set to a non-null value, the value of this attribute in state will not change."
}

// PlanModifyString implements the plan modification logic.
func (m int64KeepNonNullState) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// Do nothing if there the state value is null.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	resp.PlanValue = req.StateValue
}
