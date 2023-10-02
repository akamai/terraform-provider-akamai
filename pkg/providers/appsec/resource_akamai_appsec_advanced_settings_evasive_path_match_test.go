package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsEvasivePathMatch_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsLogging ID", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAdvancedSettingsEvasivePathMatchResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsEvasivePathMatch/EvasivePathMatch.json"), &getResponse)
		require.NoError(t, err)

		updateResponse := appsec.UpdateAdvancedSettingsEvasivePathMatchResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsEvasivePathMatch/EvasivePathMatch.json"), &updateResponse)
		require.NoError(t, err)

		removeResponse := appsec.RemoveAdvancedSettingsEvasivePathMatchResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsEvasivePathMatch/EvasivePathMatch.json"), &removeResponse)
		require.NoError(t, err)

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
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsEvasivePathMatch/match_by_id.tf"),
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
