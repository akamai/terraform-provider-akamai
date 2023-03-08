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

// TypeMappings returns a map of mappings for given rule format
// that is: mappings of values to its expected form in json
// the keys are "{behavior}.{option_name}.{value}".
func TypeMappings(ruleFormat string) map[string]any {
	return schemasRegistry.typeMappings(ruleFormat)
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
			Required:     true,
			MaxItems:     1,
			ExactlyOneOf: registeredVersions,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "",
					},
					"advanced_override": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "",
					},
					"comments": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "",
					},
					"is_secure": {
						Optional:    true,
						Type:        schema.TypeBool,
						Description: "",
					},
					"criteria_must_satisfy": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "",
					},
					"uuid": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "",
					},
					"template_uuid": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "",
					},
					"template_link": {
						Optional:    true,
						Type:        schema.TypeString,
						Description: "",
					},
					"criteria_locked": {
						Optional:    true,
						Type:        schema.TypeBool,
						Description: "",
					},
					"custom_override": {
						Optional: true,
						Type:     schema.TypeList,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Optional:    true,
									Type:        schema.TypeString,
									Description: "",
								},
								"override_id": {
									Optional:    true,
									Type:        schema.TypeString,
									Description: "",
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
						Description: "",
					},
					"criterion": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Resource{
							Schema: ruleFormat.criteriaSchemas,
						},
						Description: "",
					},
					"variable": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Required:    true,
									Type:        schema.TypeString,
									Description: "",
								},
								"value": {
									Required:    true,
									Type:        schema.TypeString,
									Description: "",
								},
								"description": {
									Required:    true,
									Type:        schema.TypeString,
									Description: "",
								},
								"hidden": {
									Required:    true,
									Type:        schema.TypeBool,
									Description: "",
								},
								"sensitive": {
									Required:    true,
									Type:        schema.TypeBool,
									Description: "",
								},
							},
						},
						Description: "",
					},
					"children": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Description: "",
					},
				},
			},
		}

	}

	return schemas
}
