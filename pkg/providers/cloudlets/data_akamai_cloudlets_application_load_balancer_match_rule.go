package cloudlets

import (
	"context"
	"encoding/json"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "The name of the rule",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 8192)),
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of Cloudlet the rule is for",
						},
						"start": {
							Type:             schema.TypeInt,
							Optional:         true,
							Description:      "The start time for this match (in seconds since the epoch)",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"end": {
							Type:             schema.TypeInt,
							Optional:         true,
							Description:      "The end time for this match (in seconds since the epoch)",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
						},
						"id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Akamai internal use only",
						},
						"matches": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Defines a set of match objects",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The type of match used",
										ValidateDiagFunc: tf.ValidateStringInSlice([]string{"header", "hostname", "path", "extension", "query",
											"cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy", "range"}),
									},
									"match_value": {
										Type:             schema.TypeString,
										Optional:         true,
										Description:      "Depends on the matchType",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 8192)),
									},
									"match_operator": {
										Type:             schema.TypeString,
										Optional:         true,
										Description:      "Valid entries for this property: contains, exists, and equals",
										ValidateDiagFunc: tf.ValidateStringInSlice([]string{"contains", "exists", "equals", ""}),
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
										Type:             schema.TypeString,
										Optional:         true,
										Description:      "For clientip, continent, countrycode, proxy, and regioncode match types, the part of the request that determines the IP address to use",
										ValidateDiagFunc: tf.ValidateStringInSlice([]string{"CONNECTING_IP", "XFF_HEADERS", "CONNECTING_IP XFF_HEADERS", ""}),
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
													ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 8192)),
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
													Description: "The array type, which can be one of the following: object, range, or simple. " +
														"Use the simple option when adding only an array of string-based values",
													ValidateDiagFunc: tf.ValidateStringInSlice([]string{"simple", "object", "range"}),
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
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "If using a URL match, this property is the URL that the Cloudlet uses to match the incoming request",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 8192)),
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
										Type:             schema.TypeString,
										Required:         true,
										Description:      "The ID of the Conditional Origin requests are forwarded to",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 8192)),
									},
								},
							},
						},
						"disabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to true, disables a rule so it is not evaluated against incoming requests.",
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
	matchRules, err := tf.GetListValue("match_rules", d)
	if err != nil {
		return handleEmptyMatchRules(err, d, "data_akamai_cloudlets_application_load_balancer_match_rule")
	}

	err = setMatchRuleSchemaType(matchRules, cloudlets.MatchRuleTypeALB)
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
			Start:         getInt64Value(rawRule, "start"),
			End:           getInt64Value(rawRule, "end"),
			ID:            getInt64Value(rawRule, "id"),
			MatchesAlways: getBoolValue(rawRule, "matches_always"),
			Disabled:      getBoolValue(rawRule, "disabled"),
		}

		rule.Matches, err = parseRuleMatches(rawRule, "matches")
		if err != nil {
			return diag.Errorf("'match_rules' - %s", err)
		}

		// Schema guarantees that "forward_settings" will be present and of type *schema.Set
		settings, ok := rawRule["forward_settings"].(*schema.Set)
		if !ok {
			return diag.Errorf("%v: 'forward_settings' should be an *schema.Set", tf.ErrInvalidType)
		}
		for _, element := range settings.List() {
			entries := element.(map[string]interface{})
			// Schema guarantees that "origin_id" will be present
			rule.ForwardSettings = cloudlets.ForwardSettingsALB{
				OriginID: entries["origin_id"].(string),
			}
		}

		rules[i] = rule
	}

	if err := rules.Validate(); err != nil {
		return diag.FromErr(err)
	}

	rulesJSON, err := json.MarshalIndent(rules, "", "  ")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("json", string(rulesJSON)); err != nil {
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	hashID, err := getMatchRulesHashID(rules)
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

	rawMatches := matches.([]interface{})
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
	omv, err := parseObjectMatchValue(m, getObjectMatchValueObjectOrSimpleOrRange)
	if err != nil {
		return nil, err
	}
	matchCriteriaALB.ObjectMatchValue = omv
	return &matchCriteriaALB, err
}
