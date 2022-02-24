package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRuleUpgrade_res_basic(t *testing.T) {
	t.Run("match by RuleUpgrade ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateRuleUpgradeResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResRuleUpgrade/RuleUpgrade.json")), &cu)

		wafModeRead := appsec.GetWAFModeResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResRuleUpgrade/WAFMode.json")), &wafModeRead)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetWAFMode",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetWAFModeRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&wafModeRead, nil)

		client.On("UpdateRuleUpgrade",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateRuleUpgradeRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Upgrade: true},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRuleUpgrade/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rule_upgrade.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
