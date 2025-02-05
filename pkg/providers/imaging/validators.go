package imaging

import (
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stringAsIntBetween(lowerLimitRaw, upperLimitRaw interface{}) schema.SchemaValidateDiagFunc {
	return func(i interface{}, _ cty.Path) diag.Diagnostics {
		value, err := strconv.Atoi(i.(string))
		if err != nil {
			return diag.Errorf("incorrect attribute value type: expect int")
		}

		if lowerLimitRaw != nil {
			lowerLimit := lowerLimitRaw.(int)
			if value < lowerLimit {
				return diag.Errorf("provided value %d must be at least %d", value, lowerLimit)
			}
		}

		if upperLimitRaw != nil {
			upperLimit := upperLimitRaw.(int)
			if value > upperLimit {
				return diag.Errorf("provided value %d must be at most %d", value, upperLimit)
			}
		}
		return nil
	}
}

func validateIsTypeFloat() schema.SchemaValidateDiagFunc {
	return func(i interface{}, _ cty.Path) diag.Diagnostics {
		if _, err := strconv.ParseFloat(i.(string), 64); err != nil {
			return diag.Errorf("incorrect attribute value type: expect float64")
		}
		return nil
	}
}

func validateIsTypeInt() schema.SchemaValidateDiagFunc {
	return func(i interface{}, _ cty.Path) diag.Diagnostics {
		_, err := strconv.Atoi(i.(string))
		if err != nil {
			return diag.Errorf("incorrect attribute value type: expect int")
		}
		return nil
	}
}

func validateIsTypeBool() schema.SchemaValidateDiagFunc {
	return func(i interface{}, _ cty.Path) diag.Diagnostics {
		value := i.(string)
		if value != "false" && value != "true" {
			return diag.Errorf("incorrect attribute value type: expect boolean")
		}
		return nil
	}
}
