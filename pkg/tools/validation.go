package tools

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"reflect"
)

func IsBlank(i interface{}, path cty.Path) diag.Diagnostics {
	if i == nil || reflect.DeepEqual(i, reflect.Zero(reflect.TypeOf(i)).Interface()) {
		return diag.Errorf("required value cannot be blank")
	}
	return nil
}
