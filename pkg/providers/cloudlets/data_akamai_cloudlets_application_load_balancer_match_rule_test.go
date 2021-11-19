package cloudlets

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataCloudletsLoadBalancerMatchRule(t *testing.T) {

	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
	}{
		"basic valid rule set": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/basic.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/basic_rules.json",
		},
		"match criteria ALB - ObjectMatchValue of Object type": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_object.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/omv_object_rules.json",
		},
		"match criteria ALB - ObjectMatchValue of Range type": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_range.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/omv_range_rules.json",
		},
		"match criteria ALB - ObjectMatchValue of Simple type": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_simple.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/omv_simple_rules.json",
		},
		"match criteria ALB - without ObjectMatchValue": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_empty.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/omv_empty_rules.json",
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
								"data.akamai_cloudlets_application_load_balancer_match_rule.test", "json",
								loadFixtureString(test.expectedJSONPath)),
							resource.TestCheckResourceAttr(
								"data.akamai_cloudlets_application_load_balancer_match_rule.test", "match_rules.0.type", "albMatchRule"),
						),
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
								"data.akamai_cloudlets_application_load_balancer_match_rule.test", "match_rules.0.type", "albMatchRule"),
						),
						ExpectError: regexp.MustCompile(test.withError),
					},
				},
			})
		})
	}

}
