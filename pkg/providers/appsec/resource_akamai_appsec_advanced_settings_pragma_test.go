package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAkamaiAdvancedSettingsPragma_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsPragma ID", func(t *testing.T) {
		client := &mockappsec{}

		updateResponse := appsec.UpdateAdvancedSettingsPragmaResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsPragma/AdvancedSettingsPragma.json"), &updateResponse)

		deleteResponse := appsec.UpdateAdvancedSettingsPragmaResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsPragma/AdvancedSettingsPragma.json"), &deleteResponse)

		getResponse := appsec.GetAdvancedSettingsPragmaResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsPragma/AdvancedSettingsPragma.json"), &getResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

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
