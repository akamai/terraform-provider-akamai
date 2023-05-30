package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsRequestBodyDataBasic(t *testing.T) {
	t.Run("match by AdvancedSettingsRequestBody ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRequestBodyResponse := appsec.GetAdvancedSettingsRequestBodyResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json"), &getRequestBodyResponse)
		require.NoError(t, err)

		client.On("GetAdvancedSettingsRequestBody",
			mock.Anything,
			appsec.GetAdvancedSettingsRequestBodyRequest{ConfigID: 43253, Version: 7},
		).Return(&getRequestBodyResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAdvancedSettingsRequestBody/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_request_body.test", "id", "43253:"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiAdvancedSettingsRequestBodyDataBasicPolicyID(t *testing.T) {
	t.Run("match by AdvancedSettingsRequestBodyPolicy ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRequestBodyResponse := appsec.GetAdvancedSettingsRequestBodyResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json"), &getRequestBodyResponse)
		require.NoError(t, err)

		client.On("GetAdvancedSettingsRequestBody",
			mock.Anything,
			appsec.GetAdvancedSettingsRequestBodyRequest{ConfigID: 43253, Version: 7, PolicyID: "test_policy"},
		).Return(&getRequestBodyResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAdvancedSettingsRequestBody/match_by_policy_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_request_body.policy", "id", "43253:test_policy"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
