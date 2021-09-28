package cloudlets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudletsApplicationLoadBalancerMatchRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudletsLoadBalancerMatchRuleRead,
		Schema: map[string]*schema.Schema{
			"match_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"end": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"matches": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"match_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"match_operator": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"case_sensitive": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"negate": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"check_ips": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"object_match_value": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"type": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"name_case_sensitive": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"name_has_wildcard": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"options": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:     schema.TypeString,
																Optional: true,
															},
															"value_has_wildcard": {
																Type:     schema.TypeBool,
																Optional: true,
															},
															"value_case_sensitive": {
																Type:     schema.TypeBool,
																Optional: true,
															},
															"value_escaped": {
																Type:     schema.TypeBool,
																Optional: true,
															},
														},
													},
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
						"match_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"matches_always": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"forward_settings": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCloudletsLoadBalancerMatchRuleRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {

	matchRules, err := tools.GetSetValue("match_rules", d)
	if err != nil {
		return diag.FromErr(err)
	}
	rules := make(cloudlets.MatchRules, matchRules.Len())

	for i, r := range matchRules.List() {
		rawRule := r.(map[string]interface{})

		// types are guaranteed by the datasource schema -> no need for type assertions
		rule := cloudlets.MatchRuleALB{
			Name:          getStringValue(rawRule, "name"),
			Type:          cloudlets.MatchRuleType(getStringValue(rawRule, "type")),
			MatchURL:      getStringValue(rawRule, "match_url"),
			Start:         getIntValue(rawRule, "start"),
			End:           getIntValue(rawRule, "end"),
			ID:            getInt64Value(rawRule, "id"),
			MatchesAlways: getBoolValue(rawRule, "matches_always"),
		}

		rule.Matches, err = parseRuleMatches(rawRule, "matches")
		if err != nil {
			return diag.FromErr(err)
		}

		// Schema guarantees that "forward_settings" will be present and of type *schema.Set
		settings, ok := rawRule["forward_settings"].(*schema.Set)
		if !ok {
			return diag.Errorf("%v: 'forward_settings' should be an *schema.Set", tools.ErrInvalidType)
		}
		for _, element := range settings.List() {
			entries := element.(map[string]interface{})
			// Schema guarantees that "origin_id" will be present
			rule.ForwardSettings = cloudlets.ForwardSettings{
				OriginID: entries["origin_id"].(string),
			}
		}

		rules[i] = rule
	}

	rulesJSON, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(rulesJSON)); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	hashID, err := getMatchRulesHashID(&rules)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hashID)

	return nil
}

func parseRuleMatches(rawRule map[string]interface{}, field string) ([]cloudlets.MatchCriteriaALB, error) {
	matches, ok := rawRule[field]
	if !ok {
		return nil, nil
	}

	rawMatches := matches.(*schema.Set).List()
	matchCriteriaALBData := make([]cloudlets.MatchCriteriaALB, len(rawMatches))
	for i, rawMatch := range rawMatches {
		alb, err := parseMatchCriteriaALB(rawMatch)
		if err != nil {
			return nil, err
		}

		matchCriteriaALBData[i] = *alb
	}
	return matchCriteriaALBData, nil
}

func parseMatchCriteriaALB(match interface{}) (*cloudlets.MatchCriteriaALB, error) {
	m := match.(map[string]interface{})
	matchCriteriaALB := cloudlets.MatchCriteriaALB{}
	matchCriteriaALB.MatchType = getStringValue(m, "match_type")
	matchCriteriaALB.MatchValue = getStringValue(m, "match_value")
	matchCriteriaALB.CaseSensitive = getBoolValue(m, "case_sensitive")
	matchCriteriaALB.Negate = getBoolValue(m, "negate")
	matchCriteriaALB.MatchOperator = cloudlets.MatchOperator(getStringValue(m, "match_operator"))
	if c, ok := m["check_ips"]; ok {
		if checkIPs, ok := c.(cloudlets.CheckIPs); ok {
			matchCriteriaALB.CheckIPs = checkIPs
		}
	}
	omv, err := parseObjectMatchValue(m)
	if err != nil {
		return nil, err
	}
	matchCriteriaALB.ObjectMatchValue = omv
	return &matchCriteriaALB, err
}

func parseObjectMatchValue(aMap map[string]interface{}) (interface{}, error) {
	v, ok := aMap["object_match_value"]
	if !ok {
		return nil, nil
	}
	rawObjects := v.(*schema.Set).List()
	for _, rawObject := range rawObjects {
		omv, ok := rawObject.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: 'object_match_value' should be an object", tools.ErrInvalidType)
		}
		if t, ok := omv["type"]; ok {
			if cloudlets.ObjectMatchValueObjectTypeSubtype(t.(string)) == cloudlets.ObjectMatchValueObjectTypeSubtypeObject {
				object := cloudlets.ObjectMatchValueObjectSubtype{Type: cloudlets.ObjectMatchValueObjectTypeSubtypeObject}
				object.Name = getStringValue(omv, "name")
				object.NameCaseSensitive = getBoolValue(omv, "name_case_sensitive")
				object.NameHasWildcard = getBoolValue(omv, "name_has_wildcard")
				opts, err := parseOMVOptions(omv)
				if err != nil {
					return nil, err
				}
				object.Options = opts
				return object, nil
			}

			rangeOrSimpleTypeSubtype := cloudlets.ObjectMatchValueRangeOrSimpleTypeSubtype(t.(string))
			if rangeOrSimpleTypeSubtype == cloudlets.ObjectMatchValueRangeOrSimpleTypeSubtypeRange ||
				rangeOrSimpleTypeSubtype == cloudlets.ObjectMatchValueRangeOrSimpleTypeSubtypeSimple {
				rangeOrSimpleType := cloudlets.ObjectMatchValueRangeOrSimpleSubtype{Type: rangeOrSimpleTypeSubtype}
				if value, ok := omv["value"]; ok {
					var val []interface{}
					if err := json.Unmarshal([]byte(value.(string)), &val); err != nil {
						return nil, err
					}
					rangeOrSimpleType.Value = val
				}
				return rangeOrSimpleType, nil
			}
		}
	}
	return nil, nil
}

func parseOMVOptions(aMap map[string]interface{}) (*cloudlets.Options, error) {
	o, ok := aMap["options"]
	if !ok {
		return nil, nil
	}
	for _, f := range o.(*schema.Set).List() {
		optionField, ok := f.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: 'options' should be an object", tools.ErrInvalidType)
		}
		options := cloudlets.Options{}
		options.Value = getStringValue(optionField, "value")
		options.ValueHasWildcard = getBoolValue(optionField, "value_has_wildcard")
		options.ValueCaseSensitive = getBoolValue(optionField, "value_case_sensitive")
		options.ValueEscaped = getBoolValue(optionField, "value_escaped")
		return &options, nil
	}
	return nil, nil
}
