package cloudlets

import (
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataCloudletsLoadBalancerMatchRule(t *testing.T) {

	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
		matchRulesSize   int
		emptyRules       bool
	}{
		"basic valid rule set": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/basic.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/basic_rules.json",
			matchRulesSize:   1,
		},
		"match criteria ALB - ObjectMatchValue of Object type": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_object.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/omv_object_rules.json",
			matchRulesSize:   2,
		},
		"match criteria ALB - ObjectMatchValue of Range type": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_range.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/omv_range_rules.json",
			matchRulesSize:   1,
		},
		"match criteria ALB - ObjectMatchValue of Simple type": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_simple.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/omv_simple_rules.json",
			matchRulesSize:   2,
		},
		"match criteria ALB - without ObjectMatchValue": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_empty.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/omv_empty_rules.json",
			matchRulesSize:   2,
		},
		"no match rules": {
			configPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/no_match_rules.tf",
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
						Check: checkMatchRulesAttr(t, "albMatchRule", "data.akamai_cloudlets_application_load_balancer_match_rule.test",
							test.expectedJSONPath, test.emptyRules, test.matchRulesSize),
					},
				},
			})
		})
	}
}

func TestIncorrectDataCloudletsLoadBalancerMatchRule(t *testing.T) {
	tests := map[string]struct {
		configPath string
		withError  string
	}{
		"match criteria ALB - ObjectMatchValueRangeSubtype with incorrect value": {
			configPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_incorrect_range.tf",
			withError:  "cannot parse abc value as an integer: strconv.ParseInt: parsing \"abc\": invalid syntax",
		},
		"match criteria ALB - missed type field in ObjectMatchValue": {
			configPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_missed_type.tf",
			withError:  "Missing required argument",
		},
		"match criteria ALB - invalid type value for ObjectMatchValue": {
			configPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_invalid_type.tf",
			withError:  `expected type to be one of \['simple', 'object', 'range'\], got invalid_type`,
		},
		"match criteria ALB - invalid match_operator value for ObjectMatchValue": {
			configPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/matches_invalid_operator.tf",
			withError:  `expected match_operator to be one of \['contains', 'exists', 'equals', ''\], got invalid`,
		},
		"match criteria ALB - invalid check_ips value for ObjectMatchValue": {
			configPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/matches_invalid_checkips.tf",
			withError:  `expected check_ips to be one of \['CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS', ''\], got invalid`,
		},
		"match criteria ALB - match_value and object_match_value together": {
			configPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/match_value_and_omv_together.tf",
			withError:  `(?s)must be blank when ObjectMatchValue is set.*must be blank when MatchValue is set`,
		},
		"match criteria ALB - no match_value and object_match_value": {
			configPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/no_match_value_and_omv.tf",
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
								"data.akamai_cloudlets_application_load_balancer_match_rule.test", "match_rules.0.type", "albMatchRule"),
						),
						ExpectError: regexp.MustCompile(test.withError),
					},
				},
			})
		})
	}
}
