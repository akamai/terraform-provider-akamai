package cloudlets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/tf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceCloudletsEdgeRedirectorMatchRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: akamaiCloudletsEdgeRedirectorMatchRuleRead,
		Schema: map[string]*schema.Schema{
			"match_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A set of rules for policy",
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
										ValidateDiagFunc: tf.ValidateStringInSlice([]string{"header", "hostname", "path", "extension", "query", "regex",
											"cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy"}),
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
													ValidateDiagFunc: tf.ValidateStringInSlice([]string{"simple", "object"}),
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
						"matches_always": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Is used in some cloudlets to support default rules (rule that is always matched)",
						},
						"use_relative_url": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "If set to relative_url, takes the path entered for the redirectUrl and sets it in the response’s Location header. " +
								"If set to copy_scheme_hostname, creates an absolute path by taking the protocol and hostname from the incoming request and combining them with path information entered for the redirectUrl. " +
								"If this property is not included, or is set to none, then the redirect_url should be fully-qualified URL",
							ValidateDiagFunc: tf.ValidateStringInSlice([]string{"relative_url", "copy_scheme_hostname", "none", ""}),
						},
						"status_code": {
							Type:             schema.TypeInt,
							Required:         true,
							Description:      "The HTTP response status code (allowed values: 301, 302, 303, 307, 308)",
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{301, 302, 303, 307, 308})),
						},
						"redirect_url": {
							Type:             schema.TypeString,
							Required:         true,
							Description:      "The URL Edge Redirector redirects the request to. If using use_relative_url, you can enter a path for the value",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 8192)),
						},
						"match_url": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "If using a URL match, this property is the URL that the Cloudlet uses to match the incoming request",
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 8192)),
						},
						"use_incoming_query_string": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to true, the Cloudlet includes the query string from the request in the rewritten or forwarded URL",
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

func akamaiCloudletsEdgeRedirectorMatchRuleRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	matchRulesList, err := tf.GetListValue("match_rules", d)
	if err != nil {
		return handleEmptyMatchRules(err, d, "data_akamai_cloudlets_edge_redirector_match_rule")
	}

	err = setMatchRuleSchemaType(matchRulesList, cloudlets.MatchRuleTypeER)
	if err != nil {
		return diag.FromErr(err)
	}

	matchRules, err := getMatchRulesER(matchRulesList)
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
		return diag.Errorf("%v: %s", tf.ErrValueSet, err.Error())
	}

	hashID, err := getMatchRulesHashID(matchRules)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashID)
	return nil
}

func getMatchCriteriaER(matches []interface{}) ([]cloudlets.MatchCriteriaER, error) {
	result := make([]cloudlets.MatchCriteriaER, 0, len(matches))
	for _, criterion := range matches {
		criterionMap, ok := criterion.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("matches is of invalid type")
		}

		omv, err := parseObjectMatchValue(criterionMap, getObjectMatchValueObjectOrSimple)
		if err != nil {
			return nil, err
		}

		matchCriterion := cloudlets.MatchCriteriaER{
			MatchType:        getStringValue(criterionMap, "match_type"),
			MatchValue:       getStringValue(criterionMap, "match_value"),
			MatchOperator:    cloudlets.MatchOperator(getStringValue(criterionMap, "match_operator")),
			CaseSensitive:    getBoolValue(criterionMap, "case_sensitive"),
			Negate:           getBoolValue(criterionMap, "negate"),
			CheckIPs:         cloudlets.CheckIPs(getStringValue(criterionMap, "check_ips")),
			ObjectMatchValue: omv,
		}

		result = append(result, matchCriterion)
	}
	return result, nil
}

func getMatchRulesER(matchRules []interface{}) (cloudlets.MatchRules, error) {
	result := make(cloudlets.MatchRules, 0, len(matchRules))
	for _, mr := range matchRules {
		matchRuleMap, ok := mr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("match rule is of invalid type: %T", mr)
		}

		matches, err := getMatchCriteriaER(matchRuleMap["matches"].([]interface{}))
		if err != nil {
			return nil, err
		}

		matchRule := cloudlets.MatchRuleER{
			Name:                     getStringValue(matchRuleMap, "name"),
			Type:                     cloudlets.MatchRuleTypeER,
			Start:                    getInt64Value(matchRuleMap, "start"),
			End:                      getInt64Value(matchRuleMap, "end"),
			MatchesAlways:            getBoolValue(matchRuleMap, "matches_always"),
			Matches:                  matches,
			UseRelativeURL:           getStringValue(matchRuleMap, "use_relative_url"),
			StatusCode:               getIntValue(matchRuleMap, "status_code"),
			RedirectURL:              getStringValue(matchRuleMap, "redirect_url"),
			MatchURL:                 getStringValue(matchRuleMap, "match_url"),
			UseIncomingQueryString:   getBoolValue(matchRuleMap, "use_incoming_query_string"),
			UseIncomingSchemeAndHost: getStringValue(matchRuleMap, "use_relative_url") == "copy_scheme_hostname",
			Disabled:                 getBoolValue(matchRuleMap, "disabled"),
		}
		result = append(result, matchRule)
	}
	return result, nil
}
