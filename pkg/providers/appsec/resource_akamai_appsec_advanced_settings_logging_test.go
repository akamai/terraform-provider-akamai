package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsLogging_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsLogging ID", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAdvancedSettingsLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"), &getResponse)
		require.NoError(t, err)

		updateResponse := appsec.UpdateAdvancedSettingsLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"), &updateResponse)
		require.NoError(t, err)

		removeResponse := appsec.RemoveAdvancedSettingsLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"), &removeResponse)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		client.On("GetAdvancedSettingsLogging",
			mock.Anything,
			appsec.GetAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7},
		).Return(&getResponse, nil)

		updateAdvancedSettingsLoggingJSON := loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/UpdateAdvancedSettingsLogging.json")
		client.On("UpdateAdvancedSettingsLogging",
			mock.Anything,
			appsec.UpdateAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: "", JsonPayloadRaw: updateAdvancedSettingsLoggingJSON},
		).Return(&updateResponse, nil)

		client.On("RemoveAdvancedSettingsLogging",
			mock.Anything,
			appsec.RemoveAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: ""},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAdvancedSettingsLogging/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_logging.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
