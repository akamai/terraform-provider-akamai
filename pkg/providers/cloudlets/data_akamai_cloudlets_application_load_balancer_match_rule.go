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
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Defines a set of rules for policy",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the rule",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of Cloudlet the rule is for",
						},
						"start": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The start time for this match (in seconds since the epoch)",
						},
						"end": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The end time for this match (in seconds since the epoch)",
						},
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Akamai internal use only",
						},
						"matches": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Defines a set of match objects",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The type of match used",
									},
									"match_value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Depends on the matchType",
									},
									"match_operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Valid entries for this property: contains, exists, and equals",
									},
									"case_sensitive": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the match is case sensitive",
									},
									"negate": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, negates the match",
									},
									"check_ips": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "For clientip, continent, countrycode, proxy, and regioncode match types, the part of the request that determines the IP address to use",
									},
									"object_match_value": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "An object used when a rule either includes more complex match criteria, like multiple value attributes, or a range match",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
													Description: "If using a match type that supports name attributes, enter the value in the incoming request to match on. " +
														"The following match types support this property: cookie, header, parameter, and query",
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
													Description: "The array type, which can be one of the following: object, range, or simple. " +
														"Use the simple option when adding only an array of string-based values",
												},
												"name_case_sensitive": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Set to true if the entry for the name property should be evaluated based on case sensitivity",
												},
												"name_has_wildcard": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Set to true if the entry for the name property includes wildcards",
												},
												"options": {
													Type:        schema.TypeSet,
													MaxItems:    1,
													Optional:    true,
													Description: "If using the object type, use this set to list the values to match on (use only with the object type)",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"value": {
																Type:        schema.TypeList,
																Elem:        &schema.Schema{Type: schema.TypeString},
																Optional:    true,
																Description: "The value attributes in the incoming request to match on",
															},
															"value_has_wildcard": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Set to true if the entries for the value property include wildcards",
															},
															"value_case_sensitive": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Set to true if the entries for the value property should be evaluated based on case sensitivity",
															},
															"value_escaped": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Set to true if provided value should be compared in escaped form",
															},
														},
													},
												},
												"value": {
													Type:        schema.TypeList,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Optional:    true,
													Description: "The value attributes in the incoming request to match on (use only with simple or range type)",
												},
											},
										},
									},
								},
							},
						},
						"match_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "If using a URL match, this property is the URL that the Cloudlet uses to match the incoming request",
						},
						"matches_always": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Is used in some cloudlets to support default rules (rule that is always matched)",
						},
						"forward_settings": {
							Type:     schema.TypeSet,
							Required: true,
							Description: "This property defines data used to construct a new request URL if all conditions are met. " +
								"If all of the conditions you set are true, then the Edge Server returns an HTTP response from the rewritten URL",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The ID of the Conditional Origin requests are forwarded to",
									},
								},
							},
						},
					},
				},
			},
			"json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A match_rules JSON structure generated from the schema",
			},
		},
	}
}

func dataSourceCloudletsLoadBalancerMatchRuleRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	matchRules, err := tools.GetListValue("match_rules", d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = setALBMatchRuleSchemaType(matchRules)
	if err != nil {
		return diag.FromErr(err)
	}

	rules := make(cloudlets.MatchRules, len(matchRules))
	for i, r := range matchRules {
		rawRule := r.(map[string]interface{})

		// types are guaranteed by the datasource schema -> no need for type assertions
		rule := cloudlets.MatchRuleALB{
			Name:          getStringValue(rawRule, "name"),
			Type:          cloudlets.MatchRuleTypeALB,
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

// setALBMatchRuleSchemaType takes ALB matchrules schema set and sets type field for every rule in set
func setALBMatchRuleSchemaType(matchRules []interface{}) error {
	for _, mr := range matchRules {
		matchRuleMap, ok := mr.(map[string]interface{})
		if !ok {
			return fmt.Errorf("match rule is of invalid type: %T", mr)
		}
		matchRuleMap["type"] = cloudlets.MatchRuleTypeALB
	}
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
	matchCriteriaALB := cloudlets.MatchCriteriaALB{
		MatchType:     getStringValue(m, "match_type"),
		MatchValue:    getStringValue(m, "match_value"),
		CaseSensitive: getBoolValue(m, "case_sensitive"),
		Negate:        getBoolValue(m, "negate"),
		MatchOperator: cloudlets.MatchOperator(getStringValue(m, "match_operator")),
	}

	if c, ok := m["check_ips"]; ok {
		if checkIPs, ok := c.(cloudlets.CheckIPs); ok {
			matchCriteriaALB.CheckIPs = checkIPs
		}
	}
	omv, err := parseALBObjectMatchValue(m)
	if err != nil {
		return nil, err
	}
	matchCriteriaALB.ObjectMatchValue = omv
	return &matchCriteriaALB, err
}

func parseALBObjectMatchValue(aMap map[string]interface{}) (interface{}, error) {
	v, ok := aMap["object_match_value"]
	if !ok {
		return nil, nil
	}
	rawObjects := v.(*schema.Set).List()
	if len(rawObjects) < 1 {
		return nil, nil
	}

	omv, ok := rawObjects[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: 'object_match_value' should be an object", tools.ErrInvalidType)
	}
	if omvType, ok := omv["type"]; ok {
		if cloudlets.ObjectMatchValueObjectType(omvType.(string)) == cloudlets.Object {
			return getOMVObjectType(omv)
		}
		if cloudlets.ObjectMatchValueSimpleType(omvType.(string)) == cloudlets.Simple {
			return getOMVSimpleType(omv), nil
		}
		if cloudlets.ObjectMatchValueRangeType(omvType.(string)) == cloudlets.Range {
			return getOMVRangeType(omv)
		}
		return nil, fmt.Errorf("'object_match_value' type '%T' is invalid. Must be one of: 'simple', 'range' or 'object'", omvType)
	}
	return nil, nil
}
