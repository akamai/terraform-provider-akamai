package property

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataPropertyRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropRulesOperation,
		Schema:      propertyRulesSchemaAttrs(),
	}
}

func dataPropRulesOperation(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// Always fail for new resources and changed values
	if d.Id() == "" || d.IsNewResource() || d.HasChanges("variables", "rules") {
		return diag.Errorf(`data "akamai_property_rules" is no longer supported - See Akamai Terraform Upgrade Guide`)
	}

	// No changes and resource already exists, must be from previous version of plugin
	return nil
}

// Returns the property rules schema attributes map
func propertyRulesSchemaAttrs() map[string]*schema.Schema {
	nameOptionSchema := func() *schema.Schema {
		return &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {Type: schema.TypeString, Required: true},
					"option": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key":    {Type: schema.TypeString, Required: true},
								"values": {Type: schema.TypeSet, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
								"value":  {Type: schema.TypeString, Optional: true},
							},
						},
					},
				},
			},
		}
	}

	// Build a nested rule schema to the given depth
	var ruleSchema func(int) *schema.Schema
	ruleSchema = func(depth int) *schema.Schema {
		attrs := map[string]*schema.Schema{
			"name":           {Type: schema.TypeString, Required: true},
			"comment":        {Type: schema.TypeString, Optional: true},
			"criteria_match": {Type: schema.TypeString, Optional: true, Default: "all"},
			"criteria":       nameOptionSchema(),
			"behavior":       nameOptionSchema(),
		}

		if depth > 0 {
			attrs["rule"] = ruleSchema(depth - 1)
		}

		return &schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Resource{Schema: attrs},
		}
	}

	return map[string]*schema.Schema{
		"variables": {Type: schema.TypeString, Optional: true},
		"json":      {Type: schema.TypeString, Computed: true, Description: "JSON Rule representation"},
		"rules": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"criteria_match": {Type: schema.TypeString, Optional: true, Default: "all"},
					"is_secure":      {Type: schema.TypeBool, Optional: true, Default: false},
					"behavior":       nameOptionSchema(),
					"rule":           ruleSchema(4), // rules tree can go max 5 levels deep (this one + 4 more)
					"variable": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name":        {Type: schema.TypeString, Required: true},
								"description": {Type: schema.TypeString, Optional: true},
								"hidden":      {Type: schema.TypeBool, Required: true},
								"sensitive":   {Type: schema.TypeBool, Required: true},
								"value":       {Type: schema.TypeString, Optional: true},
							},
						},
					},
				},
			},
		},
	}
}
