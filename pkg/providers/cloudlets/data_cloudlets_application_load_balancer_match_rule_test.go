package cloudlets

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAkamaiCloudletsLoadBalancerMatchRuleRead(t *testing.T) {

	tests := map[string]struct {
		configPath       string
		expectedJSONPath string
	}{
		"basic valid rule set": {
			configPath:       "testdata/TestAkamaiCloudletsLoadBalancerMatchRuleRead/basic.tf",
			expectedJSONPath: "testdata/TestAkamaiCloudletsLoadBalancerMatchRuleRead/rules/basic_rules.json",
		},
		"match criteria ALB - ObjectMatchValueObjectSubtype": {
			configPath:       "testdata/TestAkamaiCloudletsLoadBalancerMatchRuleRead/omv_object.tf",
			expectedJSONPath: "testdata/TestAkamaiCloudletsLoadBalancerMatchRuleRead/rules/omv_object_rules.json",
		},
		"match criteria ALB - ObjectMatchValueRangeSubtype": {
			configPath:       "testdata/TestAkamaiCloudletsLoadBalancerMatchRuleRead/omv_range.tf",
			expectedJSONPath: "testdata/TestAkamaiCloudletsLoadBalancerMatchRuleRead/rules/omv_range_rules.json",
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
