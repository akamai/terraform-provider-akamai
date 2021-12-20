package cloudlets

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataCloudletsEdgeRedirectorMatchRule(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
	}{
		"valid all vars map": {
			configPath:       "testdata/TestDataCloudletsEdgeRedirectorMatchRule/vars_map.tf",
			expectedJSONPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/rules/rules_out.json",
		},
		"valid minimal vars map": {
			configPath:       "testdata/TestDataCloudletsEdgeRedirectorMatchRule/minimal_vars_map.tf",
			expectedJSONPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/rules/minimal_rules_out.json",
		},
		"valid vars map wth empty use_relative_url": {
			configPath:       "testdata/TestDataCloudletsEdgeRedirectorMatchRule/empty_relative_url_vars_map.tf",
			expectedJSONPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/rules/empty_relative_url_rules_out.json",
		},
		"match criteria ER - without ObjectMatchValue": {
			configPath:       "testdata/TestDataCloudletsEdgeRedirectorMatchRule/omv_empty.tf",
			expectedJSONPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/rules/omv_empty_rules.json",
		},
		"match criteria ER -ObjectMatchValue of Simple type": {
			configPath:       "testdata/TestDataCloudletsEdgeRedirectorMatchRule/omv_simple.tf",
			expectedJSONPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/rules/omv_simple_rules.json",
		},
		"match criteria ER -ObjectMatchValue of Object type": {
			configPath:       "testdata/TestDataCloudletsEdgeRedirectorMatchRule/omv_object.tf",
			expectedJSONPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/rules/omv_object_rules.json",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(test.configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_edge_redirector_match_rule.test", "json",
								loadFixtureString(test.expectedJSONPath)),
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_edge_redirector_match_rule.test", "match_rules.0.type", "erMatchRule"),
						),
					},
				},
			})
		})
	}
}

func TestIncorrectDataCloudletsEdgeRedirectorMatchRule(t *testing.T) {
	tests := map[string]struct {
		configPath string
		withError  string
	}{
		"match criteria ER - missed type field in ObjectMatchValue": {
			configPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/omv_missed_type.tf",
			withError:  "Missing required argument",
		},
		"match criteria ER - invalid type value for ObjectMatchValue": {
			configPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/omv_invalid_type.tf",
			withError:  `expected type to be one of \['simple', 'object'\], got invalid_type`,
		},
		"match criteria ER - invalid match_operator value": {
			configPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/matches_invalid_operator.tf",
			withError:  `expected match_operator to be one of \['contains', 'exists', 'equals', ''\], got invalid`,
		},
		"match criteria ER - invalid check_ips value": {
			configPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/matches_invalid_checkips.tf",
			withError:  `expected check_ips to be one of \['CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS', ''\], got invalid`,
		},
		"match criteria ER - invalid status_code": {
			configPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/invalid_status_code.tf",
			withError:  `expected status_code to be one of \[301 302 303 307 308\], got 111`,
		},
		"match criteria ER - match_value and object_match_value together": {
			configPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/match_value_and_omv_together.tf",
			withError:  `(?s)must be blank when ObjectMatchValue is set.*must be blank when MatchValue is set`,
		},
		"match criteria ER - no match_value and object_match_value": {
			configPath: "testdata/TestDataCloudletsEdgeRedirectorMatchRule/no_match_value_and_omv.tf",
			withError:  `(?s)cannot be blank when ObjectMatchValue is blank.*cannot be blank when MatchValue is blank`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(test.configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_edge_redirector_match_rule.test", "match_rules.0.type", "erMatchRule"),
						),
						ExpectError: regexp.MustCompile(test.withError),
					},
				},
			})
		})
	}
}
