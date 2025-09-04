package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsPragma_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsPragma ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateResponse := appsec.UpdateAdvancedSettingsPragmaResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsPragma/AdvancedSettingsPragma.json"), &updateResponse)
		require.NoError(t, err)

		deleteResponse := appsec.UpdateAdvancedSettingsPragmaResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsPragma/AdvancedSettingsPragma.json"), &deleteResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAdvancedSettingsPragmaResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsPragma/AdvancedSettingsPragma.json"), &getResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetAdvancedSettingsPragma",
			testutils.MockContext,
			appsec.GetAdvancedSettingsPragmaRequest{ConfigID: 43253, Version: 7},
		).Return(&getResponse, nil)

		client.On("UpdateAdvancedSettingsPragma",
			testutils.MockContext,
			appsec.UpdateAdvancedSettingsPragmaRequest{ConfigID: 43253, Version: 7, PolicyID: "", JsonPayloadRaw: json.RawMessage("{\"action\":\"REMOVE\"}\n")},
		).Return(&updateResponse, nil)

		client.On("UpdateAdvancedSettingsPragma",
			testutils.MockContext,
			appsec.UpdateAdvancedSettingsPragmaRequest{ConfigID: 43253, Version: 7, PolicyID: "", JsonPayloadRaw: json.RawMessage("{}")},
		).Return(&updateResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsPragma/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_pragma_header.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
