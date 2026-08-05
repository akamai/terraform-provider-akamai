// Package test contains reusable structure and functions used during testing of sub-providers
package test

import (
	"fmt"
	"maps"
	"slices"
	"strings"

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
	resourceName    string
	attributes      map[string]checkData
	setElements     []setElementCheck
	typeSetElements []resource.TestCheckFunc
}

// setElementCheck records an expectation that a Set-typed attribute contains an element with the given value.
type setElementCheck struct {
	attr string
	val  string
}

// AttributeBatch is a type that allows to gather a group of unprefixed attributes with their values
type AttributeBatch map[string]string

// NewStateChecker creates a new instance of a StateChecker that checks attributes for a resource with provided name.
func NewStateChecker(resourceName string) StateChecker {
	return StateChecker{
		attributes:   map[string]checkData{},
		resourceName: resourceName,
	}
}

func (c StateChecker) clone() StateChecker {
	copied := NewStateChecker(c.resourceName)
	maps.Copy(copied.attributes, c.attributes)
	copied.typeSetElements = append(copied.typeSetElements, c.typeSetElements...)
	return copied
}

// Build processes all attributes and creates checks for them based on assigned values.
func (c StateChecker) Build() resource.TestCheckFunc {

	if len(c.attributes) == 0 && len(c.setElements) == 0 && len(c.typeSetElements) == 0 {
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
	for _, e := range c.setElements {
		checks = append(checks, resource.TestCheckTypeSetElemAttr(c.resourceName, e.attr+".*", e.val))
	}
	checks = append(checks, c.typeSetElements...)

	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// CheckEqual adds a check for provided attribute name and corresponding value.
func (c StateChecker) CheckEqual(attr, val string) StateChecker {
	copied := c.clone()
	// TODO
	// If we check equal timeouts.delete = 5m
	// and we check missing for timeouts in the parent, we have a collision
	// Add a check, probably in Build
	copied.setElements = append(copied.setElements, c.setElements...)
	copied.attributes[attr] = checkData{
		value: val,
	}
	return copied
}

// CheckEqualBatch adds checks for a batch of attributes with their values.
// Prefix parameter defines that all attributes in the batch should have common prefix.
// Usually common prefix is needed for all attributes that are part of a nested structure or list.
func (c StateChecker) CheckEqualBatch(prefix string, batch AttributeBatch) StateChecker {
	checker := c
	for attr, val := range batch {
		checker = checker.CheckEqual(prefix+attr, val)
	}
	return checker
}

// CheckMissing adds a check for a provided attribute name to not be present in the state.
func (c StateChecker) CheckMissing(attr string) StateChecker {
	copied := c.clone()
	copied.setElements = append(copied.setElements, c.setElements...)
	copied.attributes[attr] = checkData{
		isMissing: true,
	}
	return copied
}

// CheckTypeSetElemAttr adds a check that a TypeSet contains the provided value.
func (c StateChecker) CheckTypeSetElemAttr(attr, value string) StateChecker {
	copied := c.clone()
	copied.typeSetElements = append(copied.typeSetElements, resource.TestCheckTypeSetElemAttr(c.resourceName, attr, value))
	return copied
}

// CheckSetContains adds a check that the Set-typed attribute at attr contains an element equal to val, without
// relying on its position: Sets aren't ordered, so a plain index-based CheckEqual isn't reliable for them.
func (c StateChecker) CheckSetContains(attr, val string) StateChecker {
	copied := NewStateChecker(c.resourceName)
	maps.Copy(copied.attributes, c.attributes)
	copied.setElements = append(slices.Clone(c.setElements), setElementCheck{attr: attr, val: val})
	return copied
}

// ImportChecker allows to check the attributes in the state after terraform import.
type ImportChecker struct {
	attributes  map[string]checkData
	setElements []setElementCheck
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
		if len(c.attributes) == 0 && len(c.setElements) == 0 {
			panic("there must be at least one check in order to build the checker")
		}

		state := s[0]

		for key, data := range c.attributes {
			if err := assertAttributeFor(state, key, data); err != nil {
				return err
			}
		}
		for _, e := range c.setElements {
			if err := assertSetContainsFor(state, e); err != nil {
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
	copied.setElements = append(copied.setElements, c.setElements...)
	copied.attributes[attr] = checkData{
		value: val,
	}
	return copied
}

// CheckEqualBatch adds checks for a batch of attributes with their values.
// Prefix parameter defines that all attributes in the batch should have common prefix.
// Usually common prefix is needed for all attributes that are part of a nested structure or list.
func (c ImportChecker) CheckEqualBatch(prefix string, batch AttributeBatch) ImportChecker {
	checker := c
	for attr, val := range batch {
		checker = checker.CheckEqual(prefix+attr, val)
	}
	return checker
}

// CheckMissing adds a check for a provided attribute name to not be present in the state.
func (c ImportChecker) CheckMissing(attr string) ImportChecker {
	copied := NewImportChecker()
	maps.Copy(copied.attributes, c.attributes)
	copied.setElements = append(copied.setElements, c.setElements...)
	copied.attributes[attr] = checkData{
		isMissing: true,
	}
	return copied
}

// CheckSetContains adds a check that the Set-typed attribute at attr contains an element equal to val, without
// relying on its position: Sets aren't ordered, so a plain index-based CheckEqual isn't reliable for them.
func (c ImportChecker) CheckSetContains(attr, val string) ImportChecker {
	copied := NewImportChecker()
	maps.Copy(copied.attributes, c.attributes)
	copied.setElements = append(slices.Clone(c.setElements), setElementCheck{attr: attr, val: val})
	return copied
}

// assertAttributeFor checks whether given attribute is present in the state and has a correct value.
func assertAttributeFor(state *terraform.InstanceState, key string, data checkData) error {
	valueInState, exists := state.Attributes[key]

	if data.isMissing && exists {
		return fmt.Errorf("attribute %q was present and has a value: %q, but shouldn't be", key, data.value)
	}
	if !data.isMissing && !exists {
		return fmt.Errorf("attribute %q was not present, but should have a value: %q", key, data.value)
	}
	if !data.isMissing && (data.value != valueInState) {
		return fmt.Errorf("attribute %q has incorrect value %q, but should have %q", key, valueInState, data.value)
	}

	return nil
}

// assertSetContainsFor checks whether the Set-typed attribute e.attr contains an element with value e.val,
// regardless of its position in the flatmap representation.
func assertSetContainsFor(state *terraform.InstanceState, e setElementCheck) error {
	prefix := e.attr + "."
	for key, val := range state.Attributes {
		if !strings.HasPrefix(key, prefix) || key == e.attr+".#" || key == e.attr+".%" {
			continue
		}
		if val == e.val {
			return nil
		}
	}
	return fmt.Errorf("attribute %q does not contain expected element %q", e.attr, e.val)
}
