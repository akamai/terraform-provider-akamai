package tools

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// DiagsWithErrors appends several errors to a diag.Diagnostics
func DiagsWithErrors(d diag.Diagnostics, errs ...error) diag.Diagnostics {
	for _, e := range errs {
		d = append(d, diag.FromErr(e)...)
	}
	return d
}

// DiagWarningf creates a diag.Diagnostics with a single Warning level diag.Diagnostic entry
func DiagWarningf(format string, a ...interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf(format, a...),
		},
	}
}
