package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiRuleUpgrade_res_basic(t *testing.T) {
	t.Run("match by RuleUpgrade ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateRuleUpgradeResponse := appsec.UpdateRuleUpgradeResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResRuleUpgrade/RuleUpgrade.json"), &updateRuleUpgradeResponse)
		require.NoError(t, err)

		getWAFModeResponse := appsec.GetWAFModeResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResRuleUpgrade/WAFMode.json"), &getWAFModeResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetWAFMode",
			mock.Anything,
			appsec.GetWAFModeRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getWAFModeResponse, nil)

		client.On("UpdateRuleUpgrade",
			mock.Anything,
			appsec.UpdateRuleUpgradeRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Upgrade: true},
		).Return(&updateRuleUpgradeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
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
