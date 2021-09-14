package cloudlets

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAkamaiCloudletsEdgeRedirectorMatchRuleRead(t *testing.T) {
	t.Run("valid all vars map", func(t *testing.T) {
		client := mockcloudlets{}
		useClient(&client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestAkamaiCloudletsEdgeRedirectorMatchRuleRead/vars_map.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_cloudlets_edge_redirector_match_rule.test", "json",
								loadFixtureString("testdata/TestAkamaiCloudletsEdgeRedirectorMatchRuleRead/rules/rules_out.json")),
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
						Config: loadFixtureString("testdata/TestAkamaiCloudletsEdgeRedirectorMatchRuleRead/minimal_vars_map.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_cloudlets_edge_redirector_match_rule.test", "json",
								loadFixtureString("testdata/TestAkamaiCloudletsEdgeRedirectorMatchRuleRead/rules/minimal_rules_out.json")),
						),
					},
				},
			})
		})
	})
}
