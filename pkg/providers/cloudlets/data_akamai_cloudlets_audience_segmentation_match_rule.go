package cloudlets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/cloudlets"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceCloudletsAudienceSegmentationMatchRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudletsAudienceSegmentationMatchRuleRead,
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
										ValidateDiagFunc: tools.ValidateStringInSlice([]string{"header", "hostname", "path", "extension", "query", "range",
											"regex", "cookie", "deviceCharacteristics", "clientip", "continent", "countrycode", "regioncode", "protocol", "method", "proxy"}),
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
													Description: "The array type, which can be one of the following: object or simple or range. " +
														"Use the simple option when adding only an array of string-based values",
													ValidateDiagFunc: tools.ValidateStringInSlice([]string{"simple", "object", "range"}),
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
						"forward_settings": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Description: "This property defines data used to construct a new request URL if all conditions are met. " +
								"If all of the conditions you set are true, then the Edge Server returns an HTTP response from the rewritten URL",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin_id": {
										Type:             schema.TypeString,
										Optional:         true,
										Description:      "The ID of the Conditional Origin requests are forwarded to",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(0, 8192)),
									},
									"use_incoming_query_string": {
										Type:     schema.TypeBool,
										Optional: true,
										Description: "If set to true, the Cloudlet includes the query string from the request " +
											"in the rewritten or forwarded URL.",
									},
									"path_and_qs": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "If a value is provided and match conditions are met, this property defines " +
											"the path/resource/query string to rewrite URL for the incoming request.",
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 8192)),
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

func dataSourceCloudletsAudienceSegmentationMatchRuleRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	matchRulesList, err := tools.GetListValue("match_rules", d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = setMatchRuleSchemaType(matchRulesList, cloudlets.MatchRuleTypeAS); err != nil {
		return diag.FromErr(err)
	}

	matchRules, err := getMatchRulesAS(matchRulesList)
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

func getMatchRulesAS(matchRules []interface{}) (cloudlets.MatchRules, error) {
	result := make(cloudlets.MatchRules, 0, len(matchRules))
	for _, mr := range matchRules {
		matchRuleMap, ok := mr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("match rule is of invalid type: %T", mr)
		}

		matches, err := getMatchCriteriaAS(matchRuleMap["matches"].([]interface{}))
		if err != nil {
			return nil, err
		}

		matchRule := cloudlets.MatchRuleAS{
			Name:     getStringValue(matchRuleMap, "name"),
			Type:     cloudlets.MatchRuleTypeAS,
			Start:    getInt64Value(matchRuleMap, "start"),
			End:      getInt64Value(matchRuleMap, "end"),
			Matches:  matches,
			MatchURL: getStringValue(matchRuleMap, "match_url"),
			Disabled: getBoolValue(matchRuleMap, "disabled"),
		}
		// Schema guarantees that "forward_settings" will be present and of type *schema.Set
		settings, ok := matchRuleMap["forward_settings"].(*schema.Set)
		if !ok {
			return nil, fmt.Errorf("%v: 'forward_settings' should be an *schema.Set", tools.ErrInvalidType)
		}
		for _, element := range settings.List() {
			entries := element.(map[string]interface{})
			matchRule.ForwardSettings = cloudlets.ForwardSettingsAS{
				OriginID:               entries["origin_id"].(string),
				PathAndQS:              entries["path_and_qs"].(string),
				UseIncomingQueryString: entries["use_incoming_query_string"].(bool),
			}
		}

		result = append(result, matchRule)
	}
	return result, nil
}

func getMatchCriteriaAS(matches []interface{}) ([]cloudlets.MatchCriteriaAS, error) {
	result := make([]cloudlets.MatchCriteriaAS, 0, len(matches))
	for _, criterion := range matches {
		criterionMap, ok := criterion.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("matches is of invalid type")
		}

		omv, err := parseObjectMatchValue(criterionMap, getObjectMatchValueObjectOrSimpleOrRange)
		if err != nil {
			return nil, err
		}

		matchCriterion := cloudlets.MatchCriteriaAS{
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
