package cloudlets

import (
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
		"match criteria ALB - ObjectMatchValueObjectSubtype": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_object.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/omv_object_rules.json",
		},
		"match criteria ALB - ObjectMatchValueRangeSubtype": {
			configPath:       "testdata/TestDataCloudletsLoadBalancerMatchRule/omv_range.tf",
			expectedJSONPath: "testdata/TestDataCloudletsLoadBalancerMatchRule/rules/omv_range_rules.json",
		},
		"match criteria ALB - empty ObjectMatchValue": {
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
						),
					},
				},
			})
		})
	}
}
