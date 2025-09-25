package modifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// StringRequiredForCreate returns a plan modifier that ensures given field is only required during creation.
func StringRequiredForCreate() planmodifier.String {
	return stringRequiredForCreate{}
}

// stringRequiredForCreate implements the plan modifier.
type stringRequiredForCreate struct{}

// Description returns a human-readable description of the plan modifier.
func (m stringRequiredForCreate) Description(_ context.Context) string {
	return "Use if you want to ensure that the field is only required during creation."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m stringRequiredForCreate) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyString implements the plan modification logic.
func (m stringRequiredForCreate) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() && req.PlanValue.IsNull() {
		resp.Diagnostics.AddError("Required Field Missing", fmt.Sprintf("field `%s` is required during creation", req.Path.String()))
		return
	}
}
