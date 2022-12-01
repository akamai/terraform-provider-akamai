package cloudlets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceCloudletsRequestControlMatchRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudletsRequestControlMatchRuleRead,
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
										ValidateDiagFunc: tools.ValidateStringInSlice([]string{"header", "hostname", "path", "extension", "query", "cookie",
											"deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy",
										}),
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
										ValidateDiagFunc: tools.ValidateStringInSlice([]string{"contains", "exists", "equals", ""}),
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
										ValidateDiagFunc: tools.ValidateStringInSlice([]string{"CONNECTING_IP", "XFF_HEADERS", "CONNECTING_IP XFF_HEADERS", ""}),
									},
									"object_match_value": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "An object used when a rule either includes more complex match criteria, like multiple value attributes",
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
													Description: "The array type, which can be one of the following: object or simple. " +
														"Use the simple option when adding only an array of string-based values",
													ValidateDiagFunc: tools.ValidateStringInSlice([]string{"simple", "object"}),
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
													Description: "The value attributes in the incoming request to match on (use only with simple type)",
												},
											},
										},
									},
								},
							},
						},
						"allow_deny": {
							Type:     schema.TypeString,
							Required: true,
							Description: "If set to allow, the request is sent to origin when all conditions are true. " +
								"If deny, the request is denied when all conditions are true. " +
								"If denybranded, the request is denied and rerouted according to the configuration of the Request Control behavior",
							ValidateDiagFunc: tools.ValidateStringInSlice([]string{"allow", "deny", "denybranded"}),
						},
						"matches_always": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Is used in some cloudlets to support default rules (rule that is always matched)",
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

func dataSourceCloudletsRequestControlMatchRuleRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	matchRulesList, err := tools.GetListValue("match_rules", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setMatchRuleSchemaType(matchRulesList, cloudlets.MatchRuleTypeRC); err != nil {
		return diag.FromErr(err)
	}

	matchRules, err := getMatchRulesRC(matchRulesList)
	if err != nil {
		return diag.Errorf("'match_rules' - %s", err)
	}

	if err := matchRules.Validate(); err != nil {
		return diag.FromErr(err)
	}

	jsonBody, err := json.MarshalIndent(matchRules, "", "  ")
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("json", string(jsonBody)); err != nil {
		return diag.Errorf("%v: %s", tools.ErrValueSet, err.Error())
	}

	hashID, err := getMatchRulesHashID(matchRules)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashID)
	return nil
}

func getMatchRulesRC(matchRules []interface{}) (cloudlets.MatchRules, error) {
	result := make(cloudlets.MatchRules, 0, len(matchRules))
	for _, mr := range matchRules {
		matchRuleMap, ok := mr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("match rule is of invalid type: %T", mr)
		}

		matches, err := getMatchCriteriaRC(matchRuleMap["matches"].([]interface{}))
		if err != nil {
			return nil, err
		}

		matchRule := cloudlets.MatchRuleRC{
			Name:          getStringValue(matchRuleMap, "name"),
			Type:          cloudlets.MatchRuleTypeRC,
			Start:         getInt64Value(matchRuleMap, "start"),
			End:           getInt64Value(matchRuleMap, "end"),
			Matches:       matches,
			MatchesAlways: getBoolValue(matchRuleMap, "matches_always"),
			AllowDeny:     cloudlets.AllowDeny(getStringValue(matchRuleMap, "allow_deny")),
			Disabled:      getBoolValue(matchRuleMap, "disabled"),
		}
		result = append(result, matchRule)
	}
	return result, nil
}

func getMatchCriteriaRC(matches []interface{}) ([]cloudlets.MatchCriteriaRC, error) {
	result := make([]cloudlets.MatchCriteriaRC, 0, len(matches))
	for _, criteria := range matches {
		criteriaMap, ok := criteria.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("matches is of invalid type")
		}

		omv, err := parseObjectMatchValue(criteriaMap, getObjectMatchValueObjectOrSimple)
		if err != nil {
			return nil, err
		}

		matchCriterion := cloudlets.MatchCriteriaRC{
			MatchType:        getStringValue(criteriaMap, "match_type"),
			MatchValue:       getStringValue(criteriaMap, "match_value"),
			MatchOperator:    cloudlets.MatchOperator(getStringValue(criteriaMap, "match_operator")),
			CaseSensitive:    getBoolValue(criteriaMap, "case_sensitive"),
			Negate:           getBoolValue(criteriaMap, "negate"),
			CheckIPs:         cloudlets.CheckIPs(getStringValue(criteriaMap, "check_ips")),
			ObjectMatchValue: omv,
		}

		result = append(result, matchCriterion)
	}
	return result, nil
}
