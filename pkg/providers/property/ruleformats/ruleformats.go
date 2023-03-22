// Package ruleformats contains logic required for akamai_property_rules_builder data source.
package ruleformats

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	// RuleFormat contains information about a specific ruleformat schema.
	RuleFormat struct {
		version          string
		behaviorsSchemas map[string]*schema.Schema
		criteriaSchemas  map[string]*schema.Schema
		typeMappings     map[string]any
		nameMappings     map[string]string
		shouldFlatten    []string
	}

	// RuleVersion contains rule format version and allows for version format conversion.
	RuleVersion string
)

// SchemaKey returns schema key under which given rule version can be found.
func (v RuleVersion) SchemaKey() string {
	return string(v)
}

// Version returns schema version in normalized format.
func (v RuleVersion) Version() string {
	s := strings.TrimPrefix(string(v), "rules_")
	return strings.ReplaceAll(s, "_", "-")
}
