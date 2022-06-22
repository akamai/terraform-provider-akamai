package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAdvancedSettingsLogging_data_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsLogging ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getLoggingResponse := appsec.GetAdvancedSettingsLoggingResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSAdvancedSettingsLogging/AdvancedSettingsLogging.json"), &getLoggingResponse)

		client.On("GetAdvancedSettingsLogging",
			mock.Anything,
			appsec.GetAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7},
		).Return(&getLoggingResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAdvancedSettingsLogging/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_logging.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
