package appsec

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ValidateActions ensure actions are correct for API call
func ValidateActions(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	//alert, deny, deny_custom_{custom_deny_id}, none
	m := map[string]bool{"alert": true, "deny": true, "none": true}
	if m[value] { // will be false if "a" is not in the map
		//it was in the map
		return warnings, errors
	}

	if !(strings.Contains(value, "deny_custom_")) {
		errors = append(errors, fmt.Errorf("%q may only contain alert, deny, deny_custom_{custom_deny_id}, none", k))
	}

	return warnings, errors
}

// VerifyIDUnchanged compares the configuration's value for the configuration ID with the resource's value
// specified in the resources's ID, to ensure that the user has not inadvertently modified the configuration's value;
// any such modifications indicate an incorrect understanding of the Update operation.
func VerifyIDUnchanged(_ context.Context, d *schema.ResourceDiff, m interface{}) error {
	meta := akamai.Meta(m)
	logger := meta.Log("APPSEC", "configIdCustomDiff")

	if d.HasChange("config_id") {
		old, new := d.GetChange("config_id")
		oldvalue := old.(int)
		newvalue := new.(int)
		if oldvalue > 0 {
			logger.Errorf("%s value %d specified in configuration differs from resource ID's value %d", "config_id", newvalue, oldvalue)
			return fmt.Errorf("%s value %d specified in configuration differs from resource ID's value %d", "config_id", newvalue, oldvalue)
		}
	}

	if d.Id() != "" {
		_, exists := d.GetOkExists("security_policy_id")

		if exists && d.HasChange("security_policy_id") {
			old, new := d.GetChange("security_policy_id")
			oldvalue := old.(string)
			newvalue := new.(string)
			if len(oldvalue) > 0 {
				logger.Errorf("%s value %s specified in configuration differs from resource ID's value %s", "security_policy_id", newvalue, oldvalue)
				return fmt.Errorf("%s value %s specified in configuration differs from resource ID's value %s", "security_policy_id", newvalue, oldvalue)
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

func splitID(id string, expectedNum int, example string) ([]string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != expectedNum {
		return nil, fmt.Errorf("ID '%s' incorrectly formatted: should be of form '%s'", id, example)
	}
	return parts, nil
}
