// Package convert provides helper functions to convert Go primitive pointer types
// into terraform-plugin-framework types, handling nil pointers as null values.
package convert

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Int64PtrToStringValue converts a *int64 to types.String.
// If the pointer is nil, it returns types.StringNull().
// Otherwise it formats the int64 value as a base-10 string.
func Int64PtrToStringValue(v *int64) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(strconv.FormatInt(*v, 10))
}
