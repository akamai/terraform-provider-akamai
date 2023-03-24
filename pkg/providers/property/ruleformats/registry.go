package ruleformats

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	registry struct {
		rules []RuleFormat
	}
)

var schemasRegistry registry

// Schemas returns map of terraform schemas for all rule formats that are in the registry.
func Schemas() map[string]*schema.Schema {
	return schemasRegistry.schemas()
}

// ShouldFlattenFunc returns a function for given rule format that can be called
// with given "{behavior_name}.{option_name}" as input, and it returns if it should get flatten
// it is useful for flattening lists that are expected to be objects in json.
func ShouldFlattenFunc(ruleFormat string) func(string) bool {
	return schemasRegistry.shouldFlattenFunc(ruleFormat)
}

// TypeMappings returns a map of mappings for a given rule format
// that is: mappings of values to its expected form in json
// the keys are "{behavior}.{option_name}.{value}".
//
// For example if a certain enum expects some of the fields to be of a different type than the others,
// type mapping will take care of the conversion from the common terraform type of that attribute.
func TypeMappings(ruleFormat string) map[string]any {
	return schemasRegistry.typeMappings(ruleFormat)
}

// NameMappings returns a map of option name mappings for a given rule format
//
// For example, field "detect_smart_dns_proxy" converted to camelCase would end up as "detectSmartDnsProxy"
// and it would be incorrect. Correct value that is expected by API should be "detectSmartDNSProxy"
// and the returned map would contain following mapping: {"detectSmartDnsProxy":"detectSmartDNSProxy"}
func NameMappings(ruleFormat string) map[string]string {
	return schemasRegistry.nameMappings(ruleFormat)
}

// RulesFormats returns a list of all rule formats that are in registry.
func RulesFormats() []RuleVersion {
	return schemasRegistry.rulesFormats()
}

func (r *registry) register(rf RuleFormat) {
	r.rules = append(r.rules, rf)
}

func (r *registry) rulesFormats() []RuleVersion {
	var rulesFormats []RuleVersion

	for _, rf := range r.rules {
		rulesFormats = append(rulesFormats, RuleVersion(rf.version))
	}

	return rulesFormats
}

func (r *registry) typeMappings(ruleFormat string) map[string]any {
	for _, r := range r.rules {
		if r.version == ruleFormat {
			return r.typeMappings
		}
	}
	return nil
}

func (r *registry) nameMappings(ruleFormat string) map[string]string {
	for _, r := range r.rules {
		if r.version == ruleFormat {
			return r.nameMappings
		}
	}
	return nil
}

func (r *registry) shouldFlattenFunc(ruleFormat string) func(string) bool {
	for _, r := range r.rules {
		if r.version != ruleFormat {
			continue
		}
		return func(s string) bool {
			for _, v := range r.shouldFlatten {
				if v == s {
					return true
				}
			}
			return false
		}
	}

	// should never happen, cannot continue
	panic("no flaten func for given rule format: " + ruleFormat)
}

func (r *registry) versions() []string {
	versions := make([]string, 0, len(r.rules))
	for _, ruleFormat := range r.rules {
		versions = append(versions, ruleFormat.version)
	}
	return versions
}

func (r *registry) schemas() map[string]*schema.Schema {
	registeredVersions := r.versions()
	schemas := map[string]*schema.Schema{}
	for _, ruleFormat := range r.rules {
		schemas[ruleFormat.version] = &schema.Schema{
			Type:         schema.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: registeredVersions,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The name of a rule",
					},
					"advanced_override": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "XML metadata of the rule",
					},
					"comments": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "The comments for a rule",
					},
					"is_secure": {
						Optional:    true,
						Type:        schema.TypeBool,
						Description: "States whether a rule is secure",
					},
					"criteria_must_satisfy": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "States whether 'all' criteria need to match or 'any'",
					},
					"uuid": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "The UUID of the rule",
					},
					"template_uuid": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "The UUID of a rule template",
					},
					"template_link": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "The template link for the rule",
					},
					"criteria_locked": {
						Optional:    true,
						Type:        schema.TypeBool,
						Description: "States whether changes to 'criterion' objects are prohibited",
					},
					"custom_override": {
						Optional: true,
						Type:     schema.TypeList,
						MaxItems: 1,
                        Description: "XML metadata of the rule"
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Optional:    true,
									Type:        schema.TypeString,
									Description: "The name of a custom override",
								},
								"override_id": {
									Optional:    true,
									Type:        schema.TypeString,
									Description: "The ID of a custom override",
								},
							},
						},
					},
					"behavior": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Resource{
							Schema: ruleFormat.behaviorsSchemas,
						},
						Description: "The list of behaviors for a rule",
					},
					"criterion": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Resource{
							Schema: ruleFormat.criteriaSchemas,
						},
						Description: "The list of criteria for a rule",
					},
					"variable": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Required:    true,
									Type:        schema.TypeString,
									Description: "The name of a variable",
								},
								"value": {
									Required:    true,
									Type:        schema.TypeString,
									Description: "The value for a variable",
								},
								"description": {
									Required:    true,
									Type:        schema.TypeString,
									Description: "The description for a variable",
								},
								"hidden": {
									Required:    true,
									Type:        schema.TypeBool,
									Description: "Whether a variable should be hidden",
								},
								"sensitive": {
									Required:    true,
									Type:        schema.TypeBool,
									Description: "States whether a variable contains sensitive information",
								},
							},
						},
						Description: "A list of variables for a rule",
					},
					"children": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Description: "A list of child rules for a particular rule in JSON format",
					},
				},
			},
		}

	}

	return schemas
}
