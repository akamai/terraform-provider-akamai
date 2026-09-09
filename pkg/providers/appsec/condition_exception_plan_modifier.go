package appsec

import (
	"context"
	"encoding/json"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// ruleConditionExceptionPlanModifier compares condition_exception JSON semantically
// by unmarshaling to SDK struct and using deep.Equal. If equal, uses state value
// to preserve the exact JSON string and avoid spurious diffs.
type ruleConditionExceptionPlanModifier struct{}

// NormalizeRuleConditionException returns a plan modifier that compares rule condition
// exception JSON semantically using SDK struct comparison.
func NormalizeRuleConditionException() planmodifier.String {
	return ruleConditionExceptionPlanModifier{}
}

func (m ruleConditionExceptionPlanModifier) Description(_ context.Context) string {
	return "Compares condition_exception JSON semantically using SDK struct deep equality"
}

func (m ruleConditionExceptionPlanModifier) MarkdownDescription(_ context.Context) string {
	return "Compares condition_exception JSON semantically using SDK struct deep equality"
}

func (m ruleConditionExceptionPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If state is null/unknown or plan is null/unknown, let default behavior handle it
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	stateJSON := req.StateValue.ValueString()
	planJSON := req.PlanValue.ValueString()

	// If strings are identical, no work needed
	if stateJSON == planJSON {
		return
	}

	// Unmarshal both to SDK struct
	var stateStruct, planStruct appsec.RuleConditionException
	if err := json.Unmarshal([]byte(stateJSON), &stateStruct); err != nil {
		return // Can't parse state, let plan value through
	}
	if err := json.Unmarshal([]byte(planJSON), &planStruct); err != nil {
		return // Can't parse plan, let plan value through
	}

	// Compare using deep.Equal
	if diff := deep.Equal(stateStruct, planStruct); diff == nil {
		// Semantically equal - use state value to preserve exact JSON
		resp.PlanValue = req.StateValue
	}
}

// attackGroupConditionExceptionPlanModifier compares condition_exception JSON semantically
// by unmarshaling to SDK struct and using deep.Equal. If equal, uses state value
// to preserve the exact JSON string and avoid spurious diffs.
type attackGroupConditionExceptionPlanModifier struct{}

// NormalizeAttackGroupConditionException returns a plan modifier that compares attack group
// condition exception JSON semantically using SDK struct comparison.
func NormalizeAttackGroupConditionException() planmodifier.String {
	return attackGroupConditionExceptionPlanModifier{}
}

func (m attackGroupConditionExceptionPlanModifier) Description(_ context.Context) string {
	return "Compares condition_exception JSON semantically using SDK struct deep equality"
}

func (m attackGroupConditionExceptionPlanModifier) MarkdownDescription(_ context.Context) string {
	return "Compares condition_exception JSON semantically using SDK struct deep equality"
}

func (m attackGroupConditionExceptionPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If state is null/unknown or plan is null/unknown, let default behavior handle it
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() || req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	stateJSON := req.StateValue.ValueString()
	planJSON := req.PlanValue.ValueString()

	// If strings are identical, no work needed
	if stateJSON == planJSON {
		return
	}

	// Unmarshal both to SDK struct
	var stateStruct, planStruct appsec.AttackGroupConditionException
	if err := json.Unmarshal([]byte(stateJSON), &stateStruct); err != nil {
		return // Can't parse state, let plan value through
	}
	if err := json.Unmarshal([]byte(planJSON), &planStruct); err != nil {
		return // Can't parse plan, let plan value through
	}

	// Compare using deep.Equal
	if diff := deep.Equal(stateStruct, planStruct); diff == nil {
		// Semantically equal - use state value to preserve exact JSON
		resp.PlanValue = req.StateValue
	}
}
