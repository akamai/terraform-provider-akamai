// Package test contains reusable structure and functions used during testing of sub-providers
package test

import (
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// checkData contains information about expected value of an attribute and whether such attribute should be present.
type checkData struct {
	isMissing bool
	value     string
}

// StateChecker allows to check the attributes in the terraform state.
type StateChecker struct {
	resourceName string
	attributes   map[string]checkData
}

// NewStateChecker creates a new instance of a StateChecker that checks attributes for a resource with provided name.
func NewStateChecker(resourceName string) StateChecker {
	return StateChecker{
		attributes:   map[string]checkData{},
		resourceName: resourceName,
	}
}

// Build processes all attributes and creates checks for them based on assigned values.
func (c StateChecker) Build() resource.TestCheckFunc {

	if len(c.attributes) == 0 {
		panic("there must be at least one check in order to build the checker")
	}

	var checks []resource.TestCheckFunc
	for key, data := range c.attributes {
		if data.isMissing {
			checks = append(checks, resource.TestCheckNoResourceAttr(c.resourceName, key))
		} else {
			checks = append(checks, resource.TestCheckResourceAttr(c.resourceName, key, data.value))
		}
	}

	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// CheckEqual adds a check for provided attribute name and corresponding value.
func (c StateChecker) CheckEqual(attr, val string) StateChecker {
	copied := NewStateChecker(c.resourceName)
	maps.Copy(copied.attributes, c.attributes)
	copied.attributes[attr] = checkData{
		value: val,
	}
	return copied
}

// CheckMissing adds a check for a provided attribute name to not be present in the state.
func (c StateChecker) CheckMissing(attr string) StateChecker {
	copied := NewStateChecker(c.resourceName)
	maps.Copy(copied.attributes, c.attributes)
	copied.attributes[attr] = checkData{
		isMissing: true,
	}
	return copied
}

// ImportChecker allows to check the attributes in the state after terraform import.
type ImportChecker struct {
	attributes map[string]checkData
}

// NewImportChecker creates a new instance of a ImportChecker that checks attributes for provided resource name after the resource is imported.
func NewImportChecker() ImportChecker {
	return ImportChecker{
		attributes: map[string]checkData{},
	}
}

// Build processes all attributes and creates checks for them based on assigned values.
func (c ImportChecker) Build() resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(c.attributes) == 0 {
			panic("there must be at least one check in order to build the checker")
		}

		state := s[0]

		for key, data := range c.attributes {
			if err := assertAttributeFor(state, key, data); err != nil {
				return err
			}
		}
		return nil
	}
}

// CheckEqual adds a check for provided attribute name and corresponding value.
func (c ImportChecker) CheckEqual(attr, val string) ImportChecker {
	copied := NewImportChecker()
	maps.Copy(copied.attributes, c.attributes)
	copied.attributes[attr] = checkData{
		value: val,
	}
	return copied
}

// CheckMissing adds a check for a provided attribute name to not be present in the state.
func (c ImportChecker) CheckMissing(attr string) ImportChecker {
	copied := NewImportChecker()
	maps.Copy(copied.attributes, c.attributes)
	copied.attributes[attr] = checkData{
		isMissing: true,
	}
	return copied
}

// assertAttributeFor checks whether given attribute is present in the state and has a correct value.
func assertAttributeFor(state *terraform.InstanceState, key string, data checkData) error {
	valueInState, exists := state.Attributes[key]

	if data.isMissing && exists {
		return fmt.Errorf("attribute %q was present and has a value: %q, but shouldn't be", key, data.value)
	}
	if !data.isMissing && !exists {
		return fmt.Errorf(fmt.Sprintf("attribute %q was not present, but should have a value: %q", key, data.value))
	}
	if !data.isMissing && (data.value != valueInState) {
		return fmt.Errorf("attribute %q has incorrect value %q, but should have %q", key, valueInState, data.value)
	}

	return nil
}
