package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsAttackPayloadLoggingResBasic(t *testing.T) {
	t.Run("match by AdvancedSettingsAttackPayloadLogging ID", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAdvancedSettingsAttackPayloadLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"), &getResponse)
		require.NoError(t, err)

		updateResponse := appsec.UpdateAdvancedSettingsAttackPayloadLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"), &updateResponse)
		require.NoError(t, err)

		removeResponse := appsec.RemoveAdvancedSettingsAttackPayloadLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"), &removeResponse)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		client.On("GetAdvancedSettingsAttackPayloadLogging",
			mock.Anything,
			appsec.GetAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: 43253, Version: 7},
		).Return(&getResponse, nil)

		updateAdvancedSettingsAttackPayloadLoggingJSON := loadFixtureBytes("testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLogging.json")
		client.On("UpdateAdvancedSettingsAttackPayloadLogging",
			mock.Anything,
			appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: "", JSONPayloadRaw: updateAdvancedSettingsAttackPayloadLoggingJSON},
		).Return(&updateResponse, nil)

		client.On("RemoveAdvancedSettingsAttackPayloadLogging",
			mock.Anything,
			appsec.RemoveAdvancedSettingsAttackPayloadLoggingRequest{
				ConfigID:     43253,
				Version:      7,
				PolicyID:     "",
				Enabled:      true,
				RequestBody:  appsec.AttackPayloadLoggingRequestBody{Type: appsec.AttackPayload},
				ResponseBody: appsec.AttackPayloadLoggingResponseBody{Type: appsec.AttackPayload},
			},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAdvancedSettingsAttackPayloadLogging/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_attack_payload_logging.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiAdvancedSettingsAttackPayloadLoggingResBasicPolicy(t *testing.T) {
	t.Run("match by AdvancedSettingsAttackPayloadLoggingPolicy ID", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAdvancedSettingsAttackPayloadLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"), &getResponse)
		require.NoError(t, err)

		updateResponse := appsec.UpdateAdvancedSettingsAttackPayloadLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"), &updateResponse)
		require.NoError(t, err)

		removeResponse := appsec.RemoveAdvancedSettingsAttackPayloadLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"), &removeResponse)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&configResponse, nil)

		client.On("GetAdvancedSettingsAttackPayloadLogging",
			mock.Anything,
			appsec.GetAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: "test_policy"},
		).Return(&getResponse, nil)

		updateAdvancedSettingsAttackPayloadLoggingJSON := loadFixtureBytes("testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLogging.json")
		client.On("UpdateAdvancedSettingsAttackPayloadLogging",
			mock.Anything,
			appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: "test_policy", JSONPayloadRaw: updateAdvancedSettingsAttackPayloadLoggingJSON},
		).Return(&updateResponse, nil)

		client.On("RemoveAdvancedSettingsAttackPayloadLogging",
			mock.Anything,
			appsec.RemoveAdvancedSettingsAttackPayloadLoggingRequest{
				ConfigID:     43253,
				Version:      7,
				PolicyID:     "test_policy",
				Enabled:      true,
				RequestBody:  appsec.AttackPayloadLoggingRequestBody{Type: appsec.AttackPayload},
				ResponseBody: appsec.AttackPayloadLoggingResponseBody{Type: appsec.AttackPayload},
			},
		).Return(&removeResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAdvancedSettingsAttackPayloadLogging/update_by_policy_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_attack_payload_logging.policy", "id", "43253:test_policy"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
