package ruleformats

import (
	"fmt"

	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/dlclark/regexp2"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateRegex(pattern string) schema.SchemaValidateDiagFunc {
	re := regexp2.MustCompile(pattern, regexp2.RE2)
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value := fmt.Sprintf("%v", i)
		matchString, err := re.MatchString(value)
		if value != "" && (err != nil || !matchString) {
			errorSummary := fmt.Sprintf("value %q does not match the pattern %q", value, pattern)
			schemaField, err := tf.GetSchemaFieldNameFromPath(path)
			if err == nil {
				errorSummary = fmt.Sprintf("value %s: %q does not match the pattern %q", schemaField, value, pattern)
			}
			return diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       errorSummary,
				AttributePath: path.Copy(),
			}}
		}
		return nil
	}
}

func validateRegexOrVariable(pattern string) schema.SchemaValidateDiagFunc {
	return validateRegex(fmt.Sprintf("%s|{{.+}}", pattern))
}

func validateAny(funcs ...schema.SchemaValidateDiagFunc) schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		diags := diag.Diagnostics{}
		for _, f := range funcs {
			d := f(i, p)
			if !d.HasError() {
				return nil
			}
			diags = append(diags, d...)
		}
		return diags
	}
}
