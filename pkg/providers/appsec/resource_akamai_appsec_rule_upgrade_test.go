package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiRuleUpgrade_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by RuleUpgrade ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updateRuleUpgradeResponse := appsec.UpdateRuleUpgradeResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRuleUpgrade/RuleUpgrade.json"), &updateRuleUpgradeResponse)
		require.NoError(t, err)

		getWAFModeResponse := appsec.GetWAFModeResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRuleUpgrade/WAFMode.json"), &getWAFModeResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetWAFMode",
			testutils.MockContext,
			appsec.GetWAFModeRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getWAFModeResponse, nil)

		client.APPSEC.On("UpdateRuleUpgrade",
			testutils.MockContext,
			appsec.UpdateRuleUpgradeRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Upgrade: true},
		).Return(&updateRuleUpgradeResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResRuleUpgrade/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_rule_upgrade.test", "id", "43253:AAAA_81230"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
