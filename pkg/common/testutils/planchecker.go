package testutils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// FieldsKnownAtPlan provides a way to verify if provided list of fields is known at plan level
type FieldsKnownAtPlan struct {
	FieldsKnown   []string
	FieldsUnknown []string
}

// CheckPlan checks a plan file and then returns an error if the plan file does not match what is expected
func (p FieldsKnownAtPlan) CheckPlan(_ context.Context, request plancheck.CheckPlanRequest, response *plancheck.CheckPlanResponse) {
	knownAtPlan := request.Plan.PlannedValues.RootModule.Resources[0].AttributeValues

	for _, field := range p.FieldsKnown {
		if _, present := knownAtPlan[field]; !present {
			response.Error = fmt.Errorf("%s not known at plan (expecting known)", field)
			return
		}
	}
	for _, field := range p.FieldsUnknown {
		if _, present := knownAtPlan[field]; present {
			response.Error = fmt.Errorf("%s known at plan (expecting not known)", field)
			return
		}
	}
}
