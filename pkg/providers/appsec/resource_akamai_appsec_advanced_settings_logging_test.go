package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsLogging_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by AdvancedSettingsLogging ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAdvancedSettingsLoggingResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"), &getResponse)
		require.NoError(t, err)

		updateResponse := appsec.UpdateAdvancedSettingsLoggingResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"), &updateResponse)
		require.NoError(t, err)

		removeResponse := appsec.RemoveAdvancedSettingsLoggingResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"), &removeResponse)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		client.APPSEC.On("GetAdvancedSettingsLogging",
			testutils.MockContext,
			appsec.GetAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7},
		).Return(&getResponse, nil)

		updateAdvancedSettingsLoggingJSON := testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsLogging/UpdateAdvancedSettingsLogging.json")
		client.APPSEC.On("UpdateAdvancedSettingsLogging",
			testutils.MockContext,
			appsec.UpdateAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: "", JsonPayloadRaw: updateAdvancedSettingsLoggingJSON},
		).Return(&updateResponse, nil)

		client.APPSEC.On("RemoveAdvancedSettingsLogging",
			testutils.MockContext,
			appsec.RemoveAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: ""},
		).Return(&removeResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsLogging/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_logging.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
