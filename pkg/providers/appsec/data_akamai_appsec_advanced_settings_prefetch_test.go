package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAkamaiAdvancedSettingsPrefetch_data_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsPrefetch ID", func(t *testing.T) {
		client := &mockappsec{}

		getPrefetchResponse := appsec.GetAdvancedSettingsPrefetchResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSAdvancedSettingsPrefetch/AdvancedSettingsPrefetch.json"), &getPrefetchResponse)

		client.On("GetAdvancedSettingsPrefetch",
			mock.Anything,
			appsec.GetAdvancedSettingsPrefetchRequest{ConfigID: 43253, Version: 7},
		).Return(&getPrefetchResponse, nil)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAdvancedSettingsPrefetch/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_prefetch.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
