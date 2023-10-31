package timeouts

import (
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// SDKDefaultTimeout is a timeout used by Hashicorp when no explicit timeout was specified
var SDKDefaultTimeout = time.Minute * 20

// ValidateDurationFormat validates if provided value is a time duration
func ValidateDurationFormat(v any, _ cty.Path) diag.Diagnostics {
	duration, ok := v.(string)
	if !ok {
		return diag.Errorf("incorrect format: expected duration")
	}
	_, err := time.ParseDuration(duration)
	if err != nil {
		return diag.Errorf("provided incorrect duration: %s", err)
	}
	return nil
}
