package cloudlets

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataCloudletsEdgeRedirectorMatchRule(t *testing.T) {
	t.Run("valid all vars map", func(t *testing.T) {
		client := mockcloudlets{}
		useClient(&client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDataCloudletsEdgeRedirectorMatchRule/vars_map.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_cloudlets_edge_redirector_match_rule.test", "json",
								loadFixtureString("testdata/TestDataCloudletsEdgeRedirectorMatchRule/rules/rules_out.json")),
						),
					},
				},
			})
		})
	})
	t.Run("valid minimal vars map", func(t *testing.T) {
		client := mockcloudlets{}
		useClient(&client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDataCloudletsEdgeRedirectorMatchRule/minimal_vars_map.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_cloudlets_edge_redirector_match_rule.test", "json",
								loadFixtureString("testdata/TestDataCloudletsEdgeRedirectorMatchRule/rules/minimal_rules_out.json")),
						),
					},
				},
			})
		})
	})
}
