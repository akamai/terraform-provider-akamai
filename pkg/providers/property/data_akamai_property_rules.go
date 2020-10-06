package property

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
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

func dataPropertyRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// get rules from the TF config
	rules, err := unmarshalRules(d)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonBody, err := json.Marshal(papi.GetRuleTreeResponse{Rules: rules})
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

// TODO this function maps values from data source schema onto papi.Rules struct
// it will not be needed if we use plain json for property rules (after rewrite)
func unmarshalRules(d *schema.ResourceData) (papi.Rules, error) {
	propertyRules := papi.Rules{Name: "default"}
	// Default Rules
	rules, ok := d.GetOk("rules")
	if !ok {
		return papi.Rules{}, nil
	}
	for _, r := range rules.(*schema.Set).List() {
		ruleTree, ok := r.(map[string]interface{})
		if ok {
			behavior, ok := ruleTree["behavior"]
			if ok {
				for _, b := range behavior.(*schema.Set).List() {
					bb, ok := b.(map[string]interface{})
					if ok {
						beh := papi.RuleBehavior{}
						beh.Name = bb["name"].(string)
						boptions, ok := bb["option"]
						if ok {
							opts, err := extractOptions(boptions.(*schema.Set))
							if err != nil {
								return papi.Rules{}, err
							}
							beh.Options = opts
						}

						// Fixup CPCode
						if beh.Name == "cpCode" {
							if cpCodeOption, ok := beh.Options["id"]; ok {
								cpCodeID, err := tools.AddPrefix(tools.ConvertToString(cpCodeOption), "cpc_")
								if err != nil {
									return papi.Rules{}, err
								}
								beh.Options = papi.RuleOptionsMap{"value": map[string]interface{}{"id": cpCodeID}}
							}
						}

						// Fixup SiteShield
						if beh.Name == "siteShield" {
							if _, ok := beh.Options["ssmap"].(string); ok {
								beh.Options = papi.RuleOptionsMap{"ssmap": map[string]interface{}{"value": beh.Options["ssmap"].(string)}}
							}
						}

						propertyRules.Behaviors = mergeBehaviors(propertyRules.Behaviors, beh)
					}
				}
			}

			criteria, ok := ruleTree["criteria"]
			if ok {
				for _, c := range criteria.(*schema.Set).List() {
					cc, ok := c.(map[string]interface{})
					if ok {
						newCriteria := papi.RuleBehavior{}
						newCriteria.Name = cc["name"].(string)
						coptions, ok := cc["option"]
						if ok {
							opts, err := extractOptions(coptions.(*schema.Set))
							if err != nil {
								return papi.Rules{}, err
							}
							newCriteria.Options = opts
						}
						propertyRules.Criteria = append(propertyRules.Criteria, newCriteria)
					}
				}
			}

			if criteriamustsatisfy, ok := ruleTree["criteria_match"]; ok {
				s, _ := criteriamustsatisfy.(string)
				switch s {
				case "all":
					propertyRules.CriteriaMustSatisfy = papi.RuleCriteriaMustSatisfyAll
				case "any":
					propertyRules.CriteriaMustSatisfy = papi.RuleCriteriaMustSatisfyAny
				}
			}

			isSecure, ok := ruleTree["is_secure"].(bool)
			if ok && isSecure {
				propertyRules.Options = papi.RuleOptions{IsSecure: isSecure}
			}
		}

		childRules, ok := ruleTree["rule"]
		if ok {
			rules, err := extractRules(childRules.(*schema.Set))
			if err != nil {
				return papi.Rules{}, err
			}
			propertyRules.Children = append(propertyRules.Children, rules...)
		}
	}

	// ADD vars from variables resource
	jsonvars, ok := d.GetOk("variables")
	if ok {
		log.Println("VARS from JSON ", jsonvars)
		variables := gjson.Parse(jsonvars.(string))
		result := gjson.Get(variables.String(), "variables")

		result.ForEach(func(key, value gjson.Result) bool {
			variableMap, ok := value.Value().(map[string]interface{})
			log.Println("VARS from JSON LOOP NAME ", variableMap["name"].(string))
			log.Println("VARS from JSON LOOP DESC ", variableMap["description"].(string))
			if ok {
				newVariable := papi.RuleVariable{}
				newVariable.Name = variableMap["name"].(string)
				newVariable.Description = variableMap["description"].(string)
				newVariable.Value = variableMap["value"].(string)
				newVariable.Hidden = variableMap["hidden"].(bool)
				newVariable.Sensitive = variableMap["sensitive"].(bool)
				propertyRules.Variables = addVariable(propertyRules.Variables, newVariable)
			}

			return true
		}) //variables
	}
	return propertyRules, nil
}
