package appsec

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// conditionValidator validates that condition attributes are valid for the given condition type
type conditionValidator struct{}

// Description returns a plain text description of the validator's behavior.
func (v conditionValidator) Description(_ context.Context) string {
	return "Condition attributes must be appropriate for the condition type"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior.
func (v conditionValidator) MarkdownDescription(_ context.Context) string {
	return "Condition attributes must be appropriate for the condition type"
}

// ValidateList performs the validation.
func (v conditionValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	var conditions []urlEvasionDefenseConditionModel
	diags := req.ConfigValue.ElementsAs(ctx, &conditions, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for i, condModel := range conditions {
		cond, convDiags := buildConditionFromModel(ctx, condModel)
		resp.Diagnostics.Append(convDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if err := cond.Validate(); err != nil {
			resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
				req.Path.AtListIndex(i),
				"Invalid Condition",
				fmt.Sprintf("Invalid condition at index %d: %s", i, err.Error()),
			))
		}
	}
}

// ValidateCondition returns a validator that validates a list of conditions
func ValidateCondition() validator.List {
	return conditionValidator{}
}
