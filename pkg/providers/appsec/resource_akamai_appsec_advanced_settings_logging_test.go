package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAkamaiAdvancedSettingsLogging_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsLogging ID", func(t *testing.T) {
		client := &mockappsec{}

		configResponse := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)

		getResponse := appsec.GetAdvancedSettingsLoggingResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"), &getResponse)

		updateResponse := appsec.UpdateAdvancedSettingsLoggingResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"), &updateResponse)

		removeResponse := appsec.RemoveAdvancedSettingsLoggingResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"), &removeResponse)

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
				IsUnitTest: true,
				Providers:  testAccProviders,
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
