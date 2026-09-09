package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsEvasivePathMatch_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by AdvancedSettingsLogging ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

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

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		client.APPSEC.On("UpdateAdvancedSettingsEvasivePathMatch",
			testutils.MockContext,
			appsec.UpdateAdvancedSettingsEvasivePathMatchRequest{ConfigID: 43253, Version: 7, PolicyID: "", EnablePathMatch: true},
		).Return(&updateResponse, nil)

		client.APPSEC.On("GetAdvancedSettingsEvasivePathMatch",
			testutils.MockContext,
			appsec.GetAdvancedSettingsEvasivePathMatchRequest{ConfigID: 43253, Version: 7},
		).Return(&getResponse, nil)

		client.APPSEC.On("RemoveAdvancedSettingsEvasivePathMatch",
			testutils.MockContext,
			appsec.RemoveAdvancedSettingsEvasivePathMatchRequest{ConfigID: 43253, Version: 7, PolicyID: ""},
		).Return(&removeResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsEvasivePathMatch/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_evasive_path_match.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
