package cloudlets

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataCloudletsAPIPrioritizationMatchRule(t *testing.T) {
	workdir := "testdata/TestDataCloudletsAPIPrioritizationMatchRule"

	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
		matchRulesSize   int
		emptyRules       bool
	}{
		"valid all vars map": {
			configPath:       fmt.Sprintf("%s/vars_map.tf", workdir),
			expectedJSONPath: fmt.Sprintf("%s/rules/rules_out.json", workdir),
			matchRulesSize:   3,
		},
		"valid minimal vars map": {
			configPath:       fmt.Sprintf("%s/minimal_vars_map.tf", workdir),
			expectedJSONPath: fmt.Sprintf("%s/rules/minimal_rules_out.json", workdir),
			matchRulesSize:   1,
		},
		"match criteria AP - ObjectMatchValue of Simple type": {
			configPath:       fmt.Sprintf("%s/omv_simple.tf", workdir),
			expectedJSONPath: fmt.Sprintf("%s/rules/omv_simple_rules.json", workdir),
			matchRulesSize:   1,
		},
		"match criteria AP - ObjectMatchValue of Object type": {
			configPath:       fmt.Sprintf("%s/omv_object.tf", workdir),
			expectedJSONPath: fmt.Sprintf("%s/rules/omv_object_rules.json", workdir),
			matchRulesSize:   1,
		},
		"no match rules": {
			configPath: fmt.Sprintf("%s/no_match_rules.tf", workdir),
			emptyRules: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, test.configPath),
						Check: checkMatchRulesAttr(t, "apMatchRule", "data.akamai_cloudlets_api_prioritization_match_rule.test",
							test.expectedJSONPath, test.emptyRules, test.matchRulesSize),
					},
				},
			})
		})
	}
}

func TestIncorrectDataCloudletsAPIPrioritizationMatchRule(t *testing.T) {
	workdir := "testdata/TestDataCloudletsAPIPrioritizationMatchRule"

	tests := map[string]struct {
		configPath     string
		withError      string
		matchRulesSize int
	}{
		"missing passThroughPercent": {
			configPath:     fmt.Sprintf("%s/missing_argument.tf", workdir),
			withError:      "Missing required argument",
			matchRulesSize: 1,
		},
		"match criteria AP - invalid type value for ObjectMatchValue": {
			configPath:     fmt.Sprintf("%s/omv_invalid_type.tf", workdir),
			withError:      `expected type to be one of \['simple', 'object'\], got range`,
			matchRulesSize: 1,
		},
		"match criteria AP - invalid match_operator value": {
			configPath: fmt.Sprintf("%s/matches_invalid_operator.tf", workdir),
			withError:  `expected match_operator to be one of \['contains', 'exists', 'equals', ''\], got invalid`,
		},
		"match criteria AP - invalid check_ips value": {
			configPath: fmt.Sprintf("%s/matches_invalid_checkips.tf", workdir),
			withError:  `expected check_ips to be one of \['CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS', ''\], got invalid`,
		},
		"match criteria AP - invalid pass_through_percent value": {
			configPath:     fmt.Sprintf("%s/invalid_pass_through_percent.tf", workdir),
			withError:      `expected pass_through_percent to be in the range \(-1.000000 - 100.000000\), got -2.000000`,
			matchRulesSize: 1,
		},
		"match criteria AP - match_value and object_match_value together": {
			configPath:     fmt.Sprintf("%s/match_value_and_omv_together.tf", workdir),
			withError:      `(?s)must be blank when ObjectMatchValue is set.*must be blank when MatchValue is set`,
			matchRulesSize: 1,
		},
		"match criteria AP - no match_value and object_match_value": {
			configPath:     fmt.Sprintf("%s/no_match_value_and_omv.tf", workdir),
			withError:      `(?s)cannot be blank when ObjectMatchValue is blank.*cannot be blank when MatchValue is blank`,
			matchRulesSize: 1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, test.configPath),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_api_prioritization_match_rule.test", "match_rules.0.type", "apMatchRule"),
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_api_prioritization_match_rule.test", "match_rules.#", strconv.Itoa(test.matchRulesSize)),
						),
						ExpectError: regexp.MustCompile(test.withError),
					},
				},
			})
		})
	}
}
