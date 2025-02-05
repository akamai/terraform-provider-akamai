package appsec

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/tf"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/meta"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ValidateActions ensure actions are correct for API call
func ValidateActions(v interface{}, path cty.Path) diag.Diagnostics {
	value := v.(string)
	schemaFieldName, err := tf.GetSchemaFieldNameFromPath(path)
	if err != nil {
		return diag.FromErr(err)
	}

	//alert, deny, deny_custom_{custom_deny_id}, none
	m := map[string]bool{"alert": true, "deny": true, "none": true}
	if m[value] { // will be false if "a" is not in the map
		//it was in the map
		return nil
	}

	if !(strings.Contains(value, "deny_custom_")) {
		return diag.Errorf("%q may only contain alert, deny, deny_custom_{custom_deny_id}, none", schemaFieldName)
	}

	return nil
}

// ValidateWithBotManActions ensure actions are correct for API call
func validateWithBotManActions(v interface{}, path cty.Path) diag.Diagnostics {
	schemaFieldName, err := tf.GetSchemaFieldNameFromPath(path)
	if err != nil {
		return diag.FromErr(err)
	}
	value, ok := v.(string)
	if !ok {
		return diag.Errorf("%q is not a string", schemaFieldName)
	}

	m := map[string]struct{}{"alert": {}, "delay": {}, "deny": {}, "monitor": {}, "none": {}, "slow": {}, "tarpit": {}}
	_, ok = m[value]
	if !(ok || strings.Contains(value, "deny_custom_") || strings.Contains(value, "cond_action_") || strings.Contains(value, "serve_alt_")) {
		return diag.Errorf("%q may only contain alert, cond_action_{action_id}, delay, deny, deny_custom_{action_id}, monitor, none, serve_alt_{action_id}, slow, tarpit", schemaFieldName)
	}

	return nil
}

// VerifyIDUnchanged compares the configuration's value for the configuration ID with the resource's value
// specified in the resources's ID, to ensure that the user has not inadvertently modified the configuration's value;
// any such modifications indicate an incorrect understanding of the Update operation.
func VerifyIDUnchanged(_ context.Context, d *schema.ResourceDiff, m interface{}) error {
	meta := meta.Must(m)
	logger := meta.Log("APPSEC", "VerifyIDUnchanged")

	if d.HasChange("config_id") {
		oldConfig, newConfig := d.GetChange("config_id")
		oldVal := oldConfig.(int)
		newVal := newConfig.(int)
		if oldVal > 0 {
			logger.Errorf("%s value %d specified in configuration differs from resource ID's value %d", "config_id", newVal, oldVal)
			return fmt.Errorf("%s value %d specified in configuration differs from resource ID's value %d", "config_id", newVal, oldVal)
		}
	}

	if d.Id() != "" {
		_, exists := d.GetOkExists("security_policy_id")

		if exists && d.HasChange("security_policy_id") {
			oldPolicy, newPolicy := d.GetChange("security_policy_id")
			oldVal := oldPolicy.(string)
			newVal := newPolicy.(string)
			if len(oldVal) > 0 {
				logger.Errorf("%s value %s specified in configuration differs from resource ID's value %s", "security_policy_id", newVal, oldVal)
				return fmt.Errorf("%s value %s specified in configuration differs from resource ID's value %s", "security_policy_id", newVal, oldVal)
			}
		}

	}
	return nil
}

func validateActionAndConditionException(action, conditionexception string) error {
	if action == "none" && conditionexception != "" {
		return fmt.Errorf("action cannot be 'none' if non-empty condition/exception is supplied")
	}
	return nil
}

func validateEmptyElementsInList(v interface{}, path cty.Path) diag.Diagnostics {
	attrStep, ok := path[0].(cty.GetAttrStep)
	if !ok {
		return diag.Errorf("value must be of the specified type")
	}
	if v.(string) == "" {
		return diag.Errorf("empty or invalid string value for config parameter %s", attrStep.Name)
	}
	return nil
}
