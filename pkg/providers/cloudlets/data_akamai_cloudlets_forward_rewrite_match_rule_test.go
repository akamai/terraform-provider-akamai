package cloudlets

import (
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataCloudletsForwardRewriteMatchRule(t *testing.T) {

	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
		matchRulesSize   int
		emptyRules       bool
	}{
		"basic valid rule set": {
			configPath:       "testdata/TestDataCloudletsForwardRewriteMatchRule/basic.tf",
			expectedJSONPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/rules/basic_rules.json",
			matchRulesSize:   1,
		},
		"match criteria FR - ObjectMatchValue of Object type": {
			configPath:       "testdata/TestDataCloudletsForwardRewriteMatchRule/omv_object.tf",
			expectedJSONPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/rules/omv_object_rules.json",
			matchRulesSize:   2,
		},
		"match criteria FR - ObjectMatchValue of Simple type": {
			configPath:       "testdata/TestDataCloudletsForwardRewriteMatchRule/omv_simple.tf",
			expectedJSONPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/rules/omv_simple_rules.json",
			matchRulesSize:   2,
		},
		"match criteria FR - without ObjectMatchValue": {
			configPath:       "testdata/TestDataCloudletsForwardRewriteMatchRule/omv_empty.tf",
			expectedJSONPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/rules/omv_empty_rules.json",
			matchRulesSize:   2,
		},
		"no match rules": {
			configPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/no_match_rules.tf",
			emptyRules: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, test.configPath),
						Check: checkMatchRulesAttr(t, "frMatchRule", "data.akamai_cloudlets_forward_rewrite_match_rule.test",
							test.expectedJSONPath, test.emptyRules, test.matchRulesSize),
					},
				},
			})
		})
	}
}

func TestIncorrectDataCloudletsForwardRewriteMatchRule(t *testing.T) {
	tests := map[string]struct {
		configPath string
		withError  string
	}{
		"match criteria FR - missed type field in ObjectMatchValue": {
			configPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/omv_missed_type.tf",
			withError:  "Missing required argument",
		},
		"match criteria FR - invalid type value for ObjectMatchValue": {
			configPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/omv_invalid_type.tf",
			withError:  `expected type to be one of \['simple', 'object'\], got invalid_type`,
		},
		"match criteria FR - invalid match_operator value": {
			configPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/matches_invalid_operator.tf",
			withError:  `expected match_operator to be one of \['contains', 'exists', 'equals', ''\], got invalid`,
		},
		"match criteria FR - invalid check_ips value": {
			configPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/matches_invalid_checkips.tf",
			withError:  `expected check_ips to be one of \['CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS', ''\], got invalid`,
		},
		"match criteria FR - match_value and object_match_value together": {
			configPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/match_value_and_omv_together.tf",
			withError:  `(?s)must be blank when ObjectMatchValue is set.*must be blank when MatchValue is set`,
		},
		"match criteria FR - no match_value and object_match_value": {
			configPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/no_match_value_and_omv.tf",
			withError:  `(?s)cannot be blank when ObjectMatchValue is blank.*cannot be blank when MatchValue is blank`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, test.configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_forward_rewrite_match_rule.test", "match_rules.0.type", "frMatchRule"),
						),
						ExpectError: regexp.MustCompile(test.withError),
					},
				},
			})
		})
	}
}
