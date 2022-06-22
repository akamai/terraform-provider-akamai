package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAdvancedSettingsEvasivePathMatch_data_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsEvasivePatchMatch ID", func(t *testing.T) {
		client := &mockappsec{}

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getEvasivePathMatchResponse := appsec.GetAdvancedSettingsEvasivePathMatchResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSAdvancedSettingsEvasivePathMatch/EvasivePathMatch.json"), &getEvasivePathMatchResponse)

		client.On("GetAdvancedSettingsEvasivePathMatch",
			mock.Anything,
			appsec.GetAdvancedSettingsEvasivePathMatchRequest{ConfigID: 43253, Version: 7},
		).Return(&getEvasivePathMatchResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAdvancedSettingsEvasivePathMatch/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_evasive_path_match.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
