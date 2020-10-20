package property

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Implements Error() for diag.Diagnostics
type diagsError diag.Diagnostics

func (diags diagsError) Error() string {
	if !diag.Diagnostics(diags).HasError() {
		// NOTE: Check your diags for lack of errors before using this
		panic("BUG: DiagsError can only be used when diag.Diagnostics contains at least one error (programmer error)")
	}

	errs := []string{}

	for _, d := range diags {
		if d.Severity != diag.Error {
			continue
		}

		s := d.Summary

		if d.Detail != "" {
			s = fmt.Sprintf("%s: %s", s, d.Detail)
		}

		errs = append(errs, s)
	}

	return strings.Join(errs, "\n")
}
