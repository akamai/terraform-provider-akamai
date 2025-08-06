package modifiers

import "github.com/hashicorp/terraform-plugin-framework/resource"

// IsCreate checks if the ModifyPlan method is being called for a create operation.
func IsCreate(request resource.ModifyPlanRequest) bool {
	return request.State.Raw.IsNull()
}

// IsUpdate checks if the ModifyPlan method is being called for an update operation.
func IsUpdate(request resource.ModifyPlanRequest) bool {
	return !request.State.Raw.IsNull() && !request.Plan.Raw.IsNull()
}

// IsDelete checks if the ModifyPlan method is being called for a delete operation.
func IsDelete(request resource.ModifyPlanRequest) bool {
	return !request.State.Raw.IsNull() && request.Plan.Raw.IsNull()
}
