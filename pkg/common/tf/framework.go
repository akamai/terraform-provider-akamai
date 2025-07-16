package tf

import "github.com/hashicorp/terraform-plugin-framework/attr"

// IsKnown checks if the given attr.Value is neither null nor unknown.
func IsKnown(value attr.Value) bool {
	return !value.IsNull() && !value.IsUnknown()
}
