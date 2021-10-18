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

func dataSourceCloudletsEdgeRedirectorMatchRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: akamaiCloudletsEdgeRedirectorMatchRuleRead,
		Schema: map[string]*schema.Schema{
			"match_rules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A set of rules for policy",
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
						"matches": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "A set of match objects",
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
										Description: "Match value, depends on the matchType",
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
								},
							},
						},
						"use_relative_url": {
							Type:     schema.TypeString,
							Required: true,
							Description: "If set to relative_url, takes the path entered for the redirectUrl and sets it in the response’s Location header. " +
								"If set to copy_scheme_hostname, creates an absolute path by taking the protocol and hostname from the incoming request and combining them with path information entered for the redirectUrl. " +
								"If this property is not included, or is set to none, then the redirect_url should be fully-qualified URL",
						},
						"status_code": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The HTTP response status code (allowed values: 301, 302, 303, 307, 308)",
						},
						"redirect_url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URL Edge Redirector redirects the request to. If using use_relative_url, you can enter a path for the value",
						},
						"match_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "If using a URL match, this property is the URL that the Cloudlet uses to match the incoming request",
						},
						"use_incoming_query_string": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to true, the Cloudlet includes the query string from the request in the rewritten or forwarded URL",
						},
						"use_incoming_scheme_and_host": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to true, the Cloudlet copies both the protocol/scheme and the hostname from the incoming request to use in the redirect URL",
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
	matchRulesSet, err := tools.GetSetValue("match_rules", d)
	if err != nil {
		return diag.FromErr(err)
	}
	matchRules, err := GetMatchRules(matchRulesSet)
	if err != nil {
		return diag.Errorf("'match_rules' - %s", err)
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

// GetMatchCriteria returns contact information from Set object
func GetMatchCriteria(set *schema.Set) ([]cloudlets.MatchCriteriaER, error) {
	matches := set.List()
	result := make([]cloudlets.MatchCriteriaER, 0, len(matches))
	for _, criterion := range matches {
		criterionMap, ok := criterion.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("matches is of invalid type")
		}

		matchCriterion := cloudlets.MatchCriteriaER{
			MatchType:     getStringValue(criterionMap, "match_type"),
			MatchValue:    getStringValue(criterionMap, "match_value"),
			MatchOperator: cloudlets.MatchOperator(getStringValue(criterionMap, "match_operator")),
			CaseSensitive: getBoolValue(criterionMap, "case_sensitive"),
			Negate:        getBoolValue(criterionMap, "negate"),
			CheckIPs:      cloudlets.CheckIPs(getStringValue(criterionMap, "check_ips")),
		}
		result = append(result, matchCriterion)
	}
	return result, nil
}

// GetMatchRules returns contact information from Set object
func GetMatchRules(set *schema.Set) (*cloudlets.MatchRules, error) {
	matchRuleList := set.List()
	result := make(cloudlets.MatchRules, 0, len(matchRuleList))
	for _, mr := range matchRuleList {
		matchRuleMap, ok := mr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("match rule is of invalid type")
		}

		matches, err := GetMatchCriteria(matchRuleMap["matches"].(*schema.Set))
		if err != nil {
			return nil, err
		}

		matchRule := cloudlets.MatchRuleER{
			Name:                     getStringValue(matchRuleMap, "name"),
			Type:                     cloudlets.MatchRuleTypeER,
			Start:                    getIntValue(matchRuleMap, "start"),
			End:                      getIntValue(matchRuleMap, "end"),
			Matches:                  matches,
			UseRelativeURL:           getStringValue(matchRuleMap, "use_relative_url"),
			StatusCode:               getIntValue(matchRuleMap, "status_code"),
			RedirectURL:              getStringValue(matchRuleMap, "redirect_url"),
			MatchURL:                 getStringValue(matchRuleMap, "match_url"),
			UseIncomingQueryString:   getBoolValue(matchRuleMap, "use_incoming_query_string"),
			UseIncomingSchemeAndHost: getBoolValue(matchRuleMap, "use_incoming_scheme_and_host"),
		}
		result = append(result, matchRule)
	}
	return &result, nil
}
