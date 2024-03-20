package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSiemSettings_res_basic(t *testing.T) {
	t.Run("match by SiemSettings ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateSiemSettingsResponse := appsec.UpdateSiemSettingsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettings.json"), &updateSiemSettingsResponse)
		require.NoError(t, err)

		getSiemSettingsResponse := appsec.GetSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettings.json"), &getSiemSettingsResponse)
		require.NoError(t, err)

		removeSiemSettingsResponse := appsec.RemoveSiemSettingsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResSiemSettings/SiemSettings.json"), &removeSiemSettingsResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSiemSettings",
			mock.Anything,
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&getSiemSettingsResponse, nil)

		client.On("UpdateSiemSettings",
			mock.Anything,
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, EnabledBotmanSiemEvents: true, SiemDefinitionID: 1, FirewallPolicyIds: []string{"12345"}},
		).Return(&updateSiemSettingsResponse, nil)

		client.On("RemoveSiemSettings",
			mock.Anything,
			appsec.RemoveSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: false, EnabledBotmanSiemEvents: false, FirewallPolicyIds: []string(nil)},
		).Return(&removeSiemSettingsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResSiemSettings/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_siem_settings.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
