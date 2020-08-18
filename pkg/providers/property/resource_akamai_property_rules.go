package property

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePropertyRules() *schema.Resource {
	return &schema.Resource{
		Create: resourcePropertyRulesCreate,
		Read:   resourcePropertyRulesRead,
		Update: resourcePropertyRulesUpdate,
		Delete: resourcePropertyRulesDelete,
		Exists: resourcePropertyRulesExists,
		Schema: akamaiPropertyRulesSchema,
	}
}

var akamaiPropertyRulesSchema = map[string]*schema.Schema{
	// rules tree can go max 5 levels deep
	"variables": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"rules": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"criteria_match": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "all",
				},
				"behavior": akpsBehavior,
				"is_secure": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"rule": &schema.Schema{
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"comment": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"criteria_match": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  "all",
							},
							"criteria": akpsCriteria,
							"behavior": akpsBehavior,
							"rule": &schema.Schema{
								Type:     schema.TypeSet,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"comment": {
											Type:     schema.TypeString,
											Optional: true,
										},
										"criteria_match": {
											Type:     schema.TypeString,
											Optional: true,
											Default:  "all",
										},
										"criteria": akpsCriteria,
										"behavior": akpsBehavior,
										"rule": &schema.Schema{
											Type:     schema.TypeSet,
											Optional: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"name": {
														Type:     schema.TypeString,
														Required: true,
													},
													"comment": {
														Type:     schema.TypeString,
														Optional: true,
													},
													"criteria_match": {
														Type:     schema.TypeString,
														Optional: true,
														Default:  "all",
													},
													"criteria": akpsCriteria,
													"behavior": akpsBehavior,
													"rule": &schema.Schema{
														Type:     schema.TypeSet,
														Optional: true,
														Elem: &schema.Resource{
															Schema: map[string]*schema.Schema{
																"name": {
																	Type:     schema.TypeString,
																	Required: true,
																},
																"comment": {
																	Type:     schema.TypeString,
																	Optional: true,
																},
																"criteria_match": {
																	Type:     schema.TypeString,
																	Optional: true,
																	Default:  "all",
																},
																"criteria": akpsCriteria,
																"behavior": akpsBehavior,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"variable": &schema.Schema{
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"description": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"hidden": {
								Type:     schema.TypeBool,
								Required: true,
							},
							"sensitive": {
								Type:     schema.TypeBool,
								Required: true,
							},
							"value": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
			},
		},
	},
	"json": {
		Type: schema.TypeString,
		//Type: schema.TypeSet,
		Computed:    true,
		Description: "JSON Rule representation",
	},
}

func resourcePropertyRulesCreate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("The akamai_property_rules resource has moved to a data source, please change 'resource \"akamai_property_rules\"' to 'data \"akamai_property_rules\"")
}

func resourcePropertyRulesDelete(d *schema.ResourceData, meta interface{}) error {
	return errors.New("The akamai_property_rules resource has moved to a data source, please change 'resource \"akamai_property_rules\"' to 'data \"akamai_property_rules\"")
}

func resourcePropertyRulesExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return false, errors.New("The akamai_property_rules resource has moved to a data source, please change 'resource \"akamai_property_rules\"' to 'data \"akamai_property_rules\"")
}

func resourcePropertyRulesRead(d *schema.ResourceData, meta interface{}) error {
	return errors.New("The akamai_property_rules resource has moved to a data source, please change 'resource \"akamai_property_rules\"' to 'data \"akamai_property_rules\"")
}

func resourcePropertyRulesUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("The akamai_property_rules resource has moved to a data source, please change 'resource \"akamai_property_rules\"' to 'data \"akamai_property_rules\"")
}
