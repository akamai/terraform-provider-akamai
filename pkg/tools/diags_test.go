package tools

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/stretchr/testify/assert"
)

func TestDiagsWithErrors(t *testing.T) {
	tests := map[string]struct {
		d        diag.Diagnostics
		errs     []error
		expected diag.Diagnostics
	}{
		"multiple errors to empty diag": {
			errs: []error{fmt.Errorf("error1"), fmt.Errorf("error2"), fmt.Errorf("error3")},
			expected: diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "error1",
				},
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "error2",
				},
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "error3",
				},
			},
		},
		"given diag not empty": {
			d: diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "error1",
				},
			},
			errs: []error{fmt.Errorf("error2"), fmt.Errorf("error3")},
			expected: diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "error1",
				},
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "error2",
				},
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "error3",
				},
			},
		},
		"no errors added": {
			d: diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "error1",
				},
			},
			expected: diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "error1",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := DiagsWithErrors(test.d, test.errs...)
			assert.Equal(t, test.expected, d)
		})
	}
}

func TestDiagWarningf(t *testing.T) {
	makeSliceOfInterface := func(inputs ...interface{}) []interface{} {
		out := make([]interface{}, len(inputs))
		copy(out, inputs)
		return out
	}

	tests := map[string]struct {
		format   string
		inputs   []interface{}
		expected diag.Diagnostics
	}{
		"summary with multiple arguments": {
			format: "%s-%d-%t-%v",
			inputs: makeSliceOfInterface("test", 1, true, []int{1, 2, 3, 4}),
			expected: diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary: "test-1-true-[1 2 3 4]",
				},
			},
		},
		"summary with no arguments": {
			format: "this is a whole summary",
			expected: diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary: "this is a whole summary",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := DiagWarningf(test.format, test.inputs...)
			assert.Equal(t, test.expected, d)
		})
	}
}