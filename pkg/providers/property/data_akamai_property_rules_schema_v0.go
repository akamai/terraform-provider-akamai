package property

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var akpsOption = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	},
}

var akpsCriteria = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"option": akpsOption,
		},
	},
}

var akpsBehavior = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"option": akpsOption,
		},
	},
}

// SchemaVersion 0 of the property resource -- this is referenced in migrations to SchemaVersion 1
func dataAkamaiPropertyRuleSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// rules tree can go max 5 levels deep
			"variables": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
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
						"is_secure": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"rule": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"comments": {
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
									"rule": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"comments": {
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
												"rule": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:     schema.TypeString,
																Required: true,
															},
															"comments": {
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
															"rule": {
																Type:     schema.TypeSet,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"comments": {
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
						"variable": {
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
		},
	}
}

// Upgrade state from schema version 0 to 1
func upgradeAkamaiPropertyRuleStateV0(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	// Delete computed and removed attributes
	removed := []string{
		"variables",
		"rules",
	}
	for _, attr := range removed {
		delete(rawState, attr)
	}

	// json field is now changed to rules, so updating json field as rules.
	if v, ok := rawState["json"]; ok {
		rawState["rules"] = v
		delete(rawState, "json")
	}

	return rawState, nil
}
