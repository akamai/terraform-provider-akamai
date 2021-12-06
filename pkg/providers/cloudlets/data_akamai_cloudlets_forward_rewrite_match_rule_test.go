package cloudlets

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataCloudletsForwardRewriteMatchRule(t *testing.T) {

	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
	}{
		"basic valid rule set": {
			configPath:       "testdata/TestDataCloudletsForwardRewriteMatchRule/basic.tf",
			expectedJSONPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/rules/basic_rules.json",
		},
		"match criteria FR - ObjectMatchValue of Object type": {
			configPath:       "testdata/TestDataCloudletsForwardRewriteMatchRule/omv_object.tf",
			expectedJSONPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/rules/omv_object_rules.json",
		},
		"match criteria FR - ObjectMatchValue of Simple type": {
			configPath:       "testdata/TestDataCloudletsForwardRewriteMatchRule/omv_simple.tf",
			expectedJSONPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/rules/omv_simple_rules.json",
		},
		"match criteria FR - without ObjectMatchValue": {
			configPath:       "testdata/TestDataCloudletsForwardRewriteMatchRule/omv_empty.tf",
			expectedJSONPath: "testdata/TestDataCloudletsForwardRewriteMatchRule/rules/omv_empty_rules.json",
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
								"data.akamai_cloudlets_forward_rewrite_match_rule.test", "json",
								loadFixtureString(test.expectedJSONPath)),
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_forward_rewrite_match_rule.test", "match_rules.0.type", "frMatchRule"),
						),
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
			withError:  "is invalid. Must be one of: 'simple' or 'object'",
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
								"data.akamai_cloudlets_forward_rewrite_match_rule.test", "match_rules.0.type", "frMatchRule"),
						),
						ExpectError: regexp.MustCompile(test.withError),
					},
				},
			})
		})
	}
}
