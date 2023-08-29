package modifiers

import (
	"context"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/framework/replacer"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SetUseStateIf returns a plan modifier that copies a known prior state
// value into the planned value if given function returns true.
// Use this when you want to suppress value changes.
// Value will remain the same after a resource update.
//
// To prevent Terraform errors, the framework automatically sets unconfigured
// and Computed attributes to an unknown value "(known after apply)" on update.
// Using this plan modifier will instead display the prior state value in the
// plan, unless a prior plan modifier adjusts the value.
func SetUseStateIf(eqFunc func(string, string) bool) planmodifier.Set {
	return setUseStateIfModifier{
		equal: eqFunc,
	}
}

// setUseStateIfModifier implements the plan modifier.
type setUseStateIfModifier struct {
	equal func(string, string) bool
}

// Description returns a human-readable description of the plan modifier.
func (m setUseStateIfModifier) Description(_ context.Context) string {
	return "Use if you want to suppress value changes. " +
		"Once provided function returns true, the value of this attribute in state will not change. "
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m setUseStateIfModifier) MarkdownDescription(_ context.Context) string {
	return "Use if you want to suppress value changes. " +
		"Once provided function returns true, the value of this attribute in state will not change. "
}

// PlanModifyString implements the plan modification logic.
func (m setUseStateIfModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// Do nothing if there is no state value.
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	planValues := []string{}
	stateValues := []string{}
	resp.Diagnostics.Append(req.PlanValue.ElementsAs(ctx, &planValues, false)...)
	resp.Diagnostics.Append(req.StateValue.ElementsAs(ctx, &stateValues, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	replaced := replacer.Replacer{
		Source:       planValues,
		Replacements: stateValues,
		EqFunc:       m.equal,
	}.Replace()

	replacedSet, diags := types.SetValueFrom(ctx, types.StringType, replaced)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.PlanValue = replacedSet
}
