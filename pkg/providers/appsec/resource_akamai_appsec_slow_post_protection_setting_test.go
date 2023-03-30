package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiSlowPostProtectionSetting_res_basic(t *testing.T) {
	t.Run("match by SlowPostProtectionSetting ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getConfigurationResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &getConfigurationResponse)
		require.NoError(t, err)
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&getConfigurationResponse, nil)

		updateSlowPostProtectionSettingResponse := appsec.UpdateSlowPostProtectionSettingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtectionSetting.json"), &updateSlowPostProtectionSettingResponse)
		require.NoError(t, err)
		client.On("UpdateSlowPostProtectionSetting",
			mock.Anything,
			appsec.UpdateSlowPostProtectionSettingRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", SlowRateThreshold: struct {
				Rate   int `json:"rate"`
				Period int `json:"period"`
			}{Rate: 10, Period: 30}, DurationThreshold: struct {
				Timeout int `json:"timeout"`
			}{Timeout: 20}},
		).Return(&updateSlowPostProtectionSettingResponse, nil)

		getSlowPostProtectionSettingsResponse := appsec.GetSlowPostProtectionSettingsResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtectionSetting.json"), &getSlowPostProtectionSettingsResponse)
		require.NoError(t, err)
		client.On("GetSlowPostProtectionSettings",
			mock.Anything,
			appsec.GetSlowPostProtectionSettingsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getSlowPostProtectionSettingsResponse, nil).Twice()

		updateSlowPostProtectionResponse := appsec.UpdateSlowPostProtectionResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtection.json"), &updateSlowPostProtectionResponse)
		require.NoError(t, err)
		client.On("UpdateSlowPostProtection",
			mock.Anything,
			appsec.UpdateSlowPostProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateSlowPostProtectionResponse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
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
