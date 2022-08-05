package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAkamaiAdvancedSettingsEvasivePathMatch_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsLogging ID", func(t *testing.T) {
		client := &mockappsec{}

		configResponse := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)

		getResponse := appsec.GetAdvancedSettingsEvasivePathMatchResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsEvasivePathMatch/EvasivePathMatch.json"), &getResponse)

		updateResponse := appsec.UpdateAdvancedSettingsEvasivePathMatchResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsEvasivePathMatch/EvasivePathMatch.json"), &updateResponse)

		removeResponse := appsec.RemoveAdvancedSettingsEvasivePathMatchResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsEvasivePathMatch/EvasivePathMatch.json"), &removeResponse)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		client.On("UpdateAdvancedSettingsEvasivePathMatch",
			mock.Anything,
			appsec.UpdateAdvancedSettingsEvasivePathMatchRequest{ConfigID: 43253, Version: 7, PolicyID: "", EnablePathMatch: true},
		).Return(&updateResponse, nil)

		client.On("GetAdvancedSettingsEvasivePathMatch",
			mock.Anything,
			appsec.GetAdvancedSettingsEvasivePathMatchRequest{ConfigID: 43253, Version: 7},
		).Return(&getResponse, nil)

		client.On("RemoveAdvancedSettingsEvasivePathMatch",
			mock.Anything,
			appsec.RemoveAdvancedSettingsEvasivePathMatchRequest{ConfigID: 43253, Version: 7, PolicyID: ""},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAdvancedSettingsEvasivePathMatch/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_evasive_path_match.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
