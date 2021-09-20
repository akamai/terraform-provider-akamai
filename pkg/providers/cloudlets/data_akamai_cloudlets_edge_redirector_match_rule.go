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
							Default:  "erMatchRule",
						},
						"start": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"end": {
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
								},
							},
						},
						"use_relative_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"status_code": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"redirect_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"match_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"use_incoming_query_string": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"use_incoming_scheme_and_host": {
							Type:     schema.TypeBool,
							Optional: true,
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
			Type:                     cloudlets.MatchRuleType(getStringValue(matchRuleMap, "type")),
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
