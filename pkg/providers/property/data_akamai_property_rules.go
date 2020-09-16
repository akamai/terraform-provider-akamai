package property

import (
	"context"
	"fmt"
	"log"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/jsonhooks-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"
)

func dataPropertyRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataPropertyRulesRead,
		Schema:      akamaiDataPropertyRulesSchema,
	}
}

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

var akamaiDataPropertyRulesSchema = map[string]*schema.Schema{
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
							"rule": {
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
										"rule": {
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
													"rule": {
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
}

func dataPropertyRulesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rules := papi.NewRules()

	// get rules from the TF config
	err := unmarshalRules(d, rules)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonBody, err := jsonhooks.Marshal(rules)
	if err != nil {
		return diag.FromErr(err)
	}

	sha := tools.GetSHAString(string(jsonBody))
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.FromErr(fmt.Errorf("%w:%q", tools.ErrValueSet, err.Error()))
	}

	d.SetId(sha)

	log.Println("[DEBUG] Done")
	return nil
}

// TODO: discuss how property rules should be handled
func unmarshalRules(d *schema.ResourceData, propertyRules *papi.Rules) error {
	// Default Rules
	rules, ok := d.GetOk("rules")
	if ok {
		for _, r := range rules.(*schema.Set).List() {
			ruleTree, ok := r.(map[string]interface{})
			if ok {
				behavior, ok := ruleTree["behavior"]
				if ok {
					for _, b := range behavior.(*schema.Set).List() {
						bb, ok := b.(map[string]interface{})
						if ok {
							beh := papi.NewBehavior()
							beh.Name = bb["name"].(string)
							boptions, ok := bb["option"]
							if ok {
								opts, err := extractOptions(boptions.(*schema.Set))
								if err != nil {
									return err
								}
								beh.Options = opts
							}

							// Fixup CPCode
							if beh.Name == "cpCode" {
								if _, ok := beh.Options["value"]; !ok {
									if _, ok := beh.Options["id"]; ok {
										cpCode := papi.NewCpCodes(nil, nil).NewCpCode()
										cpCode.CpcodeID = beh.Options["id"].(string)
										beh.Options = papi.OptionValue{"value": papi.OptionValue{"id": cpCode.ID()}}
									}
								}
							}

							// Fixup SiteShield
							if beh.Name == "siteShield" {
								if _, ok := beh.Options["ssmap"].(string); ok {
									beh.Options = papi.OptionValue{"ssmap": papi.OptionValue{"value": beh.Options["ssmap"].(string)}}
								}
							}

							propertyRules.Rule.MergeBehavior(beh)
						}
					}
				}

				criteria, ok := ruleTree["criteria"]
				if ok {
					for _, c := range criteria.(*schema.Set).List() {
						cc, ok := c.(map[string]interface{})
						if ok {
							newCriteria := papi.NewCriteria()
							newCriteria.Name = cc["name"].(string)
							coptions, ok := cc["option"]
							if ok {
								opts, err := extractOptions(coptions.(*schema.Set))
								if err != nil {
									return err
								}
								newCriteria.Options = opts
							}
							propertyRules.Rule.MergeCriteria(newCriteria)
						}
					}
				}

				isSecure, ok := ruleTree["is_secure"].(bool)
				if ok && isSecure {
					propertyRules.Rule.Options.IsSecure = isSecure
				}
			}

			childRules, ok := ruleTree["rule"]
			if ok {
				rules, err := extractRules(childRules.(*schema.Set))
				if err != nil {
					return err
				}
				for _, rule := range rules {
					propertyRules.Rule.MergeChildRule(rule)
				}
			}
		}

		// ADD vars from variables resource
		jsonvars, ok := d.GetOk("variables")
		if ok {
			log.Println("VARS from JSON ", jsonvars)
			variables := gjson.Parse(jsonvars.(string))
			result := gjson.Get(variables.String(), "variables") //gjson.GetMany(variables.String(),"variables.#.name","variables.#.description","variables.#.value","variables.#.hidden","variables.#.sensitive" )

			result.ForEach(func(key, value gjson.Result) bool {
				variableMap, ok := value.Value().(map[string]interface{})
				log.Println("VARS from JSON LOOP NAME ", variableMap["name"].(string))
				log.Println("VARS from JSON LOOP DESC ", variableMap["description"].(string))
				if ok {
					newVariable := papi.NewVariable()
					newVariable.Name = variableMap["name"].(string)
					newVariable.Description = variableMap["description"].(string)
					newVariable.Value = variableMap["value"].(string)
					newVariable.Hidden = variableMap["hidden"].(bool)
					newVariable.Sensitive = variableMap["sensitive"].(bool)
					propertyRules.Rule.AddVariable(newVariable)
				}

				return true
			}) //variables
		}
	}
	return nil
}
