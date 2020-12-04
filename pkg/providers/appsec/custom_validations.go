package appsec

import (
	"fmt"
	"strings"
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
