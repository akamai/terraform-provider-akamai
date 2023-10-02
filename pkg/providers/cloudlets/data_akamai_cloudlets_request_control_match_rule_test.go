package cloudlets

import (
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataCloudletsRequestControlMatchRule(t *testing.T) {
	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
		matchRulesSize   int
		emptyRules       bool
	}{
		"basic valid rule set": {
			configPath:       "testdata/TestDataCloudletsRequestControlMatchRule/basic.tf",
			expectedJSONPath: "testdata/TestDataCloudletsRequestControlMatchRule/rules/basic_rules.json",
			matchRulesSize:   1,
		},
		"valid rule set with simple value": {
			configPath:       "testdata/TestDataCloudletsRequestControlMatchRule/omv_simple.tf",
			expectedJSONPath: "testdata/TestDataCloudletsRequestControlMatchRule/rules/omv_simple_rules.json",
			matchRulesSize:   1,
		},
		"valid rule set with object value": {
			configPath:       "testdata/TestDataCloudletsRequestControlMatchRule/omv_object.tf",
			expectedJSONPath: "testdata/TestDataCloudletsRequestControlMatchRule/rules/omv_object_rules.json",
			matchRulesSize:   1,
		},
		"valid rule set without value": {
			configPath:       "testdata/TestDataCloudletsRequestControlMatchRule/omv_empty.tf",
			expectedJSONPath: "testdata/TestDataCloudletsRequestControlMatchRule/rules/omv_empty_rules.json",
			matchRulesSize:   1,
		},
		"valid complex rule set": {
			configPath:       "testdata/TestDataCloudletsRequestControlMatchRule/omv_complex.tf",
			expectedJSONPath: "testdata/TestDataCloudletsRequestControlMatchRule/rules/omv_complex_rules.json",
			matchRulesSize:   2,
		},
		"no match rules": {
			configPath: "testdata/TestDataCloudletsRequestControlMatchRule/no_match_rules.tf",
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
						Check: checkMatchRulesAttr(t, "igMatchRule", "data.akamai_cloudlets_request_control_match_rule.test",
							test.expectedJSONPath, test.emptyRules, test.matchRulesSize),
					},
				},
			})
		})
	}
}

func TestIncorrectDataCloudletsRequestControlMatchRule(t *testing.T) {
	tests := map[string]struct {
		configPath string
		withError  string
	}{
		"match criteria RC - missed allow_deny field": {
			configPath: "testdata/TestDataCloudletsRequestControlMatchRule/missing_allow_deny.tf",
			withError:  `The argument "allow_deny" is required, but no definition was found.`,
		},
		"match criteria RC - missed type field in ObjectMatchValue": {
			configPath: "testdata/TestDataCloudletsRequestControlMatchRule/missing_type.tf",
			withError:  "Missing required argument",
		},
		"match criteria RC - no match_value and object_match_value": {
			configPath: "testdata/TestDataCloudletsRequestControlMatchRule/missing_value.tf",
			withError:  `(?s)cannot be blank when ObjectMatchValue is blank.*cannot be blank when MatchValue is blank`,
		},
		"match criteria RC - match_value and object_match_value together": {
			configPath: "testdata/TestDataCloudletsRequestControlMatchRule/duplicate_values.tf",
			withError:  `(?s)must be blank when ObjectMatchValue is set.*must be blank when MatchValue is set`,
		},
		"match criteria RC - incorrect value of type field": {
			configPath: "testdata/TestDataCloudletsRequestControlMatchRule/invalid_enum_type.tf",
			withError:  `expected type to be one of \['simple', 'object'\], got incorrect_type`,
		},
		"match criteria RC - incorrect value of match_operator field": {
			configPath: "testdata/TestDataCloudletsRequestControlMatchRule/invalid_enum_match_operator.tf",
			withError:  `expected match_operator to be one of \['contains', 'exists', 'equals', ''\], got invalid_operator`,
		},
		"match criteria RC - incorrect value of check_ips field": {
			configPath: "testdata/TestDataCloudletsRequestControlMatchRule/invalid_enum_check_ips.tf",
			withError:  `expected check_ips to be one of \['CONNECTING_IP', 'XFF_HEADERS', 'CONNECTING_IP XFF_HEADERS', ''\], got incorrect_check_ips`,
		},
		"match criteria RC - incorrect value of match_type field": {
			configPath: "testdata/TestDataCloudletsRequestControlMatchRule/invalid_enum_match_type.tf",
			withError:  `expected match_type to be one of \['header', 'hostname', 'path', 'extension', 'query', 'cookie', 'deviceCharacteristics', 'clientip', 'continent', 'countrycode', 'regioncode', 'protocol', 'method', 'proxy'\], got invalid_match_type`,
		},
		"match criteria RC - incorrect value of allow_deny field": {
			configPath: "testdata/TestDataCloudletsRequestControlMatchRule/invalid_enum_allow_deny.tf",
			withError:  `expected allow_deny to be one of \['allow', 'deny', 'denybranded'\], got invalid_value`,
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
								"data.akamai_cloudlets_request_control_match_rule.test", "match_rules.0.type", "igMatchRule"),
						),
						ExpectError: regexp.MustCompile(test.withError),
					},
				},
			})
		})
	}
}
