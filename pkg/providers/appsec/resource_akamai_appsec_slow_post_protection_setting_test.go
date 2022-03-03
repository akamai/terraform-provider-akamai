package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSlowPostProtectionSetting_res_basic(t *testing.T) {
	t.Run("match by SlowPostProtectionSetting ID", func(t *testing.T) {
		client := &mockappsec{}

		getConfigurationResponse := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &getConfigurationResponse)
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&getConfigurationResponse, nil)

		updateSlowPostProtectionSettingResponse := appsec.UpdateSlowPostProtectionSettingResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtectionSetting.json")), &updateSlowPostProtectionSettingResponse)
		client.On("UpdateSlowPostProtectionSetting",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateSlowPostProtectionSettingRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", SlowRateThreshold: struct {
				Rate   int "json:\"rate\""
				Period int "json:\"period\""
			}{Rate: 10, Period: 30}, DurationThreshold: struct {
				Timeout int "json:\"timeout\""
			}{Timeout: 20}},
		).Return(&updateSlowPostProtectionSettingResponse, nil)

		getSlowPostProtectionSettingsResponse := appsec.GetSlowPostProtectionSettingsResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtectionSetting.json")), &getSlowPostProtectionSettingsResponse)
		client.On("GetSlowPostProtectionSettings",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSlowPostProtectionSettingsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getSlowPostProtectionSettingsResponse, nil).Twice()

		updateSlowPostProtectionResponse := appsec.UpdateSlowPostProtectionResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtection.json")), &updateSlowPostProtectionResponse)
		client.On("UpdateSlowPostProtection",
			mock.Anything,
			appsec.UpdateSlowPostProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateSlowPostProtectionResponse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSlowPostProtectionSetting/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_slow_post.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
