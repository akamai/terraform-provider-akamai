package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAdvancedSettingsLogging_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsLogging ID", func(t *testing.T) {
		client := &mockappsec{}

		configResponse := appsec.GetConfigurationResponse{}
		configResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(configResponseJSON), &configResponse)

		getResponse := appsec.GetAdvancedSettingsLoggingResponse{}
		getResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
		json.Unmarshal([]byte(getResponseJSON), &getResponse)

		updateResponse := appsec.UpdateAdvancedSettingsLoggingResponse{}
		expectResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
		json.Unmarshal([]byte(expectResponseJSON), &updateResponse)

		removeResponse := appsec.RemoveAdvancedSettingsLoggingResponse{}
		removeResponseJSON := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
		json.Unmarshal([]byte(removeResponseJSON), &removeResponse)

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
