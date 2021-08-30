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

		allProtectionsFalse := appsec.GetPolicyProtectionsResponse{}
		tempJSON := compactJSON(loadFixtureBytes("testdata/TestResIPGeoProtection/PolicyProtections.json"))
		json.Unmarshal([]byte(tempJSON), &allProtectionsFalse)

		cu := appsec.UpdateSlowPostProtectionSettingResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtectionSetting.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetSlowPostProtectionSettingsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtectionSetting.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		cup := appsec.UpdateSlowPostProtectionResponse{}
		expectJSUP := compactJSON(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtection.json"))
		json.Unmarshal([]byte(expectJSUP), &cup)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetSlowPostProtectionSettings",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSlowPostProtectionSettingsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("GetPolicyProtections",
			mock.Anything,
			appsec.GetPolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allProtectionsFalse, nil).Once()

		client.On("UpdatePolicyProtections",
			mock.Anything,
			appsec.UpdatePolicyProtectionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allProtectionsFalse, nil).Once()

		client.On("UpdateSlowPostProtectionSetting",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateSlowPostProtectionSettingRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", SlowRateThreshold: struct {
				Rate   int "json:\"rate\""
				Period int "json:\"period\""
			}{Rate: 10, Period: 30}, DurationThreshold: struct {
				Timeout int "json:\"timeout\""
			}{Timeout: 20}},
		).Return(&cu, nil)

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
