package test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// PlanCheckFunc is a function bridge to the plancheck.PlanCheck interface.
type PlanCheckFunc func(context.Context, plancheck.CheckPlanRequest, *plancheck.CheckPlanResponse)

// CheckPlan implements the plancheck.PlanCheck interface.
func (f PlanCheckFunc) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest,
	resp *plancheck.CheckPlanResponse) {

	f(ctx, req, resp)
}

// PrintPlan prints the plan details to stdout. Used for debugging purposes.
// Example usage:
//
//	ConfigPlanChecks: resource.ConfigPlanChecks{
//	    PreApply: []plancheck.PlanCheck{
//		       test.PlanCheckFunc(test.PrintPlan),
//	    },
//	},
func PrintPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {

	if req.Plan == nil {
		resp.Error = fmt.Errorf("plan is nil")
	}

	for _, resourceChange := range req.Plan.ResourceChanges {
		fmt.Printf("Resource Address: %s\n", resourceChange.Address)

		beforeJSON, err := json.MarshalIndent(resourceChange.Change.Before, "", "  ")
		if err != nil {
			resp.Error = fmt.Errorf("error marshaling before state: %s", err)
			return
		}
		fmt.Println("Before State:")
		fmt.Println(string(beforeJSON))

		afterJSON, err := json.MarshalIndent(resourceChange.Change.After, "", "  ")
		if err != nil {
			resp.Error = fmt.Errorf("error marshaling after state: %s", err)
			return
		}
		fmt.Println("After State:")
		fmt.Println(string(afterJSON))

		afterUnknownJSON, err := json.MarshalIndent(resourceChange.Change.AfterUnknown, "", "  ")
		if err != nil {
			resp.Error = fmt.Errorf("error marshaling after unknown state: %s", err)
			return
		}
		fmt.Println("After Unknown:")
		fmt.Println(string(afterUnknownJSON))
		fmt.Println("--------------------------------------------------")
	}
}
