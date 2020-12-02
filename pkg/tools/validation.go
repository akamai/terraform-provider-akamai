package tools

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
)

// AggregateValidations takes any number of schema.SchemaValidateDiagFunc and executes them one by one
// it returns a diagnostics object containing combined results of each validation function
func AggregateValidations(funcs ...schema.SchemaValidateDiagFunc) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		for _, f := range funcs {
			if err := f(i, path); err != nil {
				diags = append(diags, err...)
			}
		}
		return diags
	}
}

// IsNotBlank verifies whether given value is not blank and returns error if it is where "blank" means:
// - nil value
// - a collection with len == 0 in case the value is a map, array or slice
// - value equal to zero-value for given type (e.g. empty string)
func IsNotBlank(i interface{}, path cty.Path) diag.Diagnostics {
	val := reflect.ValueOf(i)
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		if val.Len() == 0 {
			return diag.Errorf("provided value cannot be blank")
		}
	default:
		if i == nil || reflect.DeepEqual(i, reflect.Zero(reflect.TypeOf(i)).Interface()) {
			return diag.Errorf("provided value cannot be blank")
		}
	}
	return nil
}

// ValidateJSON checks whether given value is a valid JSON object
func ValidateJSON(val interface{}, path cty.Path) diag.Diagnostics {
	if str, ok := val.(string); ok {
		var target map[string]interface{}
		if err := json.Unmarshal([]byte(str), &target); err != nil {
			return diag.FromErr(fmt.Errorf("invalid JSON: %s", err))
		}
		return nil
	}
	return diag.Errorf("value is not a string: %s", val)
}
