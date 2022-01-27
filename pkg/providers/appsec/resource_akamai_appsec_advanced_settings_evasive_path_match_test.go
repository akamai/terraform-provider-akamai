package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAdvancedSettingsEvasivePathMatch_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsLogging ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateAdvancedSettingsEvasivePathMatchResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsEvasivePathMatch/EvasivePathMatch.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetAdvancedSettingsEvasivePathMatchResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsEvasivePathMatch/EvasivePathMatch.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveAdvancedSettingsEvasivePathMatchResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsEvasivePathMatch/EvasivePathMatch.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetAdvancedSettingsEvasivePathMatch",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAdvancedSettingsEvasivePathMatchRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateAdvancedSettingsEvasivePathMatch",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAdvancedSettingsEvasivePathMatchRequest{ConfigID: 43253, Version: 7, PolicyID: "", EnablePathMatch: true},
		).Return(&cu, nil)

		client.On("RemoveAdvancedSettingsEvasivePathMatch",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveAdvancedSettingsEvasivePathMatchRequest{ConfigID: 43253, Version: 7, PolicyID: ""},
		).Return(&crd, nil)

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
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
