package modifiers

import (
	"context"
	"strings"

	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/jsonutil"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// PreventJSONReorder returns a plan modifier that suppresses diffs when two
// JSON values are semantically equal, treating arrays as unordered sets at
// every nesting level.
func PreventJSONReorder() planmodifier.String {
	return stringSuppressJSONListReorderModifier{}
}

type stringSuppressJSONListReorderModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m stringSuppressJSONListReorderModifier) Description(_ context.Context) string {
	return "Ignore diffs when the JSON values are semantically equal, treating arrays as unordered sets at every nesting level."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m stringSuppressJSONListReorderModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyString sets the plan equal to state when both JSON values are
// semantically equal ignoring array order, preventing a spurious diff.
func (m stringSuppressJSONListReorderModifier) PlanModifyString(
	_ context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	stateStr := strings.TrimSpace(req.StateValue.ValueString())
	planStr := strings.TrimSpace(req.PlanValue.ValueString())

	if stateStr == planStr {
		return
	}

	if stateStr == "" || planStr == "" {
		return
	}

	if jsonutil.EqualIgnoringArrayOrder(stateStr, planStr) {
		resp.PlanValue = req.StateValue
	}
}
