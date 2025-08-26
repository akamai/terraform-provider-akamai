package date

import (
	"time"

	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/date"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TimeRFC3339Value returns a types.String representation of the given time.Time in RFC3339 format.
func TimeRFC3339Value(t time.Time) types.String {
	return types.StringValue(date.FormatRFC3339(t))
}

// TimeRFC3339PointerValue returns a types.String representation of the given *time.Time in RFC3339 format.
func TimeRFC3339PointerValue(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return TimeRFC3339Value(*t)
}

// TimeRFC3339NanoValue returns a types.String representation of the given time.Time in RFC3339Nano format.
func TimeRFC3339NanoValue(t time.Time) types.String {
	return types.StringValue(date.FormatRFC3339Nano(t))
}

// TimeRFC3339NanoPointerValue returns a types.String representation of the given *time.Time in RFC3339Nano format.
func TimeRFC3339NanoPointerValue(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return TimeRFC3339NanoValue(*t)
}
