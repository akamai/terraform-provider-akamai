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

		cu := appsec.UpdateAdvancedSettingsLoggingResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetAdvancedSettingsLoggingResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveAdvancedSettingsLoggingResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetAdvancedSettingsLogging",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		updateAdvancedSettingsLoggingJSON := loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/UpdateAdvancedSettingsLogging.json")
		client.On("UpdateAdvancedSettingsLogging",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: "", JsonPayloadRaw: updateAdvancedSettingsLoggingJSON},
		).Return(&cu, nil)

		client.On("RemoveAdvancedSettingsLogging",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: ""},
		).Return(&crd, nil)

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
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
