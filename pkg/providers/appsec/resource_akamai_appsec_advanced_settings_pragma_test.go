package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsPragma_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsPragma ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateResponse := appsec.UpdateAdvancedSettingsPragmaResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsPragma/AdvancedSettingsPragma.json"), &updateResponse)
		require.NoError(t, err)

		deleteResponse := appsec.UpdateAdvancedSettingsPragmaResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsPragma/AdvancedSettingsPragma.json"), &deleteResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAdvancedSettingsPragmaResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsPragma/AdvancedSettingsPragma.json"), &getResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetAdvancedSettingsPragma",
			mock.Anything,
			appsec.GetAdvancedSettingsPragmaRequest{ConfigID: 43253, Version: 7},
		).Return(&getResponse, nil)

		client.On("UpdateAdvancedSettingsPragma",
			mock.Anything,
			appsec.UpdateAdvancedSettingsPragmaRequest{ConfigID: 43253, Version: 7, PolicyID: "", JsonPayloadRaw: json.RawMessage("{\"action\":\"REMOVE\"}\n")},
		).Return(&updateResponse, nil)

		client.On("UpdateAdvancedSettingsPragma",
			mock.Anything,
			appsec.UpdateAdvancedSettingsPragmaRequest{ConfigID: 43253, Version: 7, PolicyID: "", JsonPayloadRaw: json.RawMessage("{}")},
		).Return(&updateResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAdvancedSettingsPragma/match_by_id.tf"),
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
