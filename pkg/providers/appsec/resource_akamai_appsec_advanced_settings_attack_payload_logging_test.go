package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsAttackPayloadLoggingConfig(t *testing.T) {
	var (
		configVersion = func(t *testing.T, configId int, client *appsec.Mock) appsec.GetConfigurationResponse {
			configResponse := appsec.GetConfigurationResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
			require.NoError(t, err)

			client.On("GetConfiguration",
				mock.Anything,
				appsec.GetConfigurationRequest{ConfigID: configId},
			).Return(&configResponse, nil)

			return configResponse
		}

		attackPayloadLoggingRead = func(t *testing.T, configId int, version int, policyId string, client *appsec.Mock, payloadPath string, numberOfTimes int) {
			attackPayloadLoggingResponse := appsec.GetAdvancedSettingsAttackPayloadLoggingResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, payloadPath), &attackPayloadLoggingResponse)
			require.NoError(t, err)

			client.On("GetAdvancedSettingsAttackPayloadLogging",
				mock.Anything,
				appsec.GetAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: configId, Version: version, PolicyID: policyId},
			).Return(&attackPayloadLoggingResponse, nil).Times(numberOfTimes)

		}

		updateAttackPayloadLogging = func(t *testing.T, updateAttackPayloadLogging appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest, client *appsec.Mock, payloadPath string, numberOfTimes int) {
			updateAttackPayloadLoggingResponse := appsec.UpdateAdvancedSettingsAttackPayloadLoggingResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, payloadPath), &updateAttackPayloadLoggingResponse)
			require.NoError(t, err)

			client.On("UpdateAdvancedSettingsAttackPayloadLogging",
				mock.Anything, updateAttackPayloadLogging,
			).Return(&updateAttackPayloadLoggingResponse, nil).Times(numberOfTimes)

		}

		removeAttackPayloadLogging = func(t *testing.T, removeAttackPayloadLogging appsec.RemoveAdvancedSettingsAttackPayloadLoggingRequest, client *appsec.Mock, payloadPath string, numberOfTimes int) {
			removeAttackPayloadLoggingResponse := appsec.RemoveAdvancedSettingsAttackPayloadLoggingResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, payloadPath), &removeAttackPayloadLoggingResponse)
			require.NoError(t, err)
			client.On("RemoveAdvancedSettingsAttackPayloadLogging",
				mock.Anything, removeAttackPayloadLogging,
			).Return(&removeAttackPayloadLoggingResponse, nil).Times(numberOfTimes)
		}
	)

	t.Run("match by AdvancedSettingsAttackPayloadLogging ID", func(t *testing.T) {
		client := &appsec.Mock{}
		configResponse := configVersion(t, 43253, client)
		payloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"
		attackPayloadLoggingRead(t, 43253, 7, "", client, payloadPath, 2)
		updateAdvancedSettingsAttackPayloadLoggingJSON := testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLogging.json")
		updateAttackPayloadLoggingRequest := appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "", JSONPayloadRaw: updateAdvancedSettingsAttackPayloadLoggingJSON}

		updatePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLogging.json"
		updateAttackPayloadLogging(t, updateAttackPayloadLoggingRequest, client, updatePayloadPath, 1)

		removeAttackPayloadLoggingRequest := appsec.RemoveAdvancedSettingsAttackPayloadLoggingRequest{
			ConfigID:     43253,
			Version:      7,
			PolicyID:     "",
			Enabled:      true,
			RequestBody:  appsec.AttackPayloadLoggingRequestBody{Type: appsec.AttackPayload},
			ResponseBody: appsec.AttackPayloadLoggingResponseBody{Type: appsec.AttackPayload},
		}

		removePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLogging.json"
		removeAttackPayloadLogging(t, removeAttackPayloadLoggingRequest, client, removePayloadPath, 1)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_attack_payload_logging.test", "id", "43253:"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by AdvancedSettingsAttackPayloadLogging disabled ID", func(t *testing.T) {
		client := &appsec.Mock{}
		configResponse := configVersion(t, 43253, client)
		payloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLoggingDisabled.json"
		attackPayloadLoggingRead(t, 43253, 7, "", client, payloadPath, 2)
		updateAdvancedSettingsAttackPayloadLoggingJSON := testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingDisabled.json")
		updateAttackPayloadLoggingRequest := appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "", JSONPayloadRaw: updateAdvancedSettingsAttackPayloadLoggingJSON}

		updatePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingDisabledResponse.json"
		updateAttackPayloadLogging(t, updateAttackPayloadLoggingRequest, client, updatePayloadPath, 1)

		removeAttackPayloadLoggingRequest := appsec.RemoveAdvancedSettingsAttackPayloadLoggingRequest{
			ConfigID:     43253,
			Version:      7,
			PolicyID:     "",
			Enabled:      true,
			RequestBody:  appsec.AttackPayloadLoggingRequestBody{Type: appsec.AttackPayload},
			ResponseBody: appsec.AttackPayloadLoggingResponseBody{Type: appsec.AttackPayload},
		}

		removePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLogging.json"
		removeAttackPayloadLogging(t, removeAttackPayloadLoggingRequest, client, removePayloadPath, 1)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/match_by_id_disabled.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_attack_payload_logging.test", "id", "43253:"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("import", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := configVersion(t, 43253, client)
		payloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"
		attackPayloadLoggingRead(t, configResponse.ID, configResponse.LatestVersion, "", client, payloadPath, 4)
		updateAdvancedSettingsAttackPayloadLoggingJSON := testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLogging.json")
		updateAttackPayloadLoggingRequest := appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "", JSONPayloadRaw: updateAdvancedSettingsAttackPayloadLoggingJSON}

		updatePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLogging.json"
		updateAttackPayloadLogging(t, updateAttackPayloadLoggingRequest, client, updatePayloadPath, 1)

		removeAttackPayloadLoggingRequest := appsec.RemoveAdvancedSettingsAttackPayloadLoggingRequest{
			ConfigID:     configResponse.ID,
			Version:      configResponse.LatestVersion,
			PolicyID:     "",
			Enabled:      true,
			RequestBody:  appsec.AttackPayloadLoggingRequestBody{Type: appsec.AttackPayload},
			ResponseBody: appsec.AttackPayloadLoggingResponseBody{Type: appsec.AttackPayload},
		}

		removePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLogging.json"
		removeAttackPayloadLogging(t, removeAttackPayloadLoggingRequest, client, removePayloadPath, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/match_by_id.tf"),
					},
					{
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateId:     "43253",
						ResourceName:      "akamai_appsec_advanced_settings_attack_payload_logging.test",
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("import policy", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := configVersion(t, 43253, client)
		payloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLoggingPolicy.json"
		attackPayloadLoggingRead(t, 43253, 7, "test_policy", client, payloadPath, 4)
		updateAdvancedSettingsAttackPayloadLoggingJSON := testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingPolicy.json")
		updateAttackPayloadLoggingRequest := appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", JSONPayloadRaw: updateAdvancedSettingsAttackPayloadLoggingJSON}

		updatePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingPolicy.json"
		updateAttackPayloadLogging(t, updateAttackPayloadLoggingRequest, client, updatePayloadPath, 1)

		removeAttackPayloadLoggingRequest := appsec.RemoveAdvancedSettingsAttackPayloadLoggingRequest{
			ConfigID:     43253,
			Version:      7,
			PolicyID:     "test_policy",
			Enabled:      true,
			RequestBody:  appsec.AttackPayloadLoggingRequestBody{Type: appsec.AttackPayload},
			ResponseBody: appsec.AttackPayloadLoggingResponseBody{Type: appsec.AttackPayload},
		}

		removePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLogging.json"
		removeAttackPayloadLogging(t, removeAttackPayloadLoggingRequest, client, removePayloadPath, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/update_by_policy_id.tf"),
					},
					{
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateId:     "43253:test_policy",
						ResourceName:      "akamai_appsec_advanced_settings_attack_payload_logging.policy",
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("match by AdvancedSettingsAttackPayloadPolicy ID", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := configVersion(t, 43253, client)
		payloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLoggingPolicy.json"

		attackPayloadLoggingRead(t, 43253, 7, "test_policy", client, payloadPath, 2)
		updateAdvancedSettingsAttackPayloadLoggingJSON := testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingPolicy.json")
		updateAttackPayloadLoggingRequest := appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", JSONPayloadRaw: updateAdvancedSettingsAttackPayloadLoggingJSON}

		updatePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingPolicy.json"
		updateAttackPayloadLogging(t, updateAttackPayloadLoggingRequest, client, updatePayloadPath, 1)

		removeAttackPayloadLoggingRequest := appsec.RemoveAdvancedSettingsAttackPayloadLoggingRequest{
			ConfigID:     43253,
			Version:      7,
			PolicyID:     "test_policy",
			Enabled:      true,
			RequestBody:  appsec.AttackPayloadLoggingRequestBody{Type: appsec.AttackPayload},
			ResponseBody: appsec.AttackPayloadLoggingResponseBody{Type: appsec.AttackPayload},
		}

		removePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingPolicy.json"
		removeAttackPayloadLogging(t, removeAttackPayloadLoggingRequest, client, removePayloadPath, 1)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/update_by_policy_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_attack_payload_logging.policy", "id", "43253:test_policy"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by AdvancedSettingsAttackPayloadPolicy ID drift", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := configVersion(t, 43253, client)
		payloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLoggingPolicy.json"

		attackPayloadLoggingRead(t, 43253, 7, "test_policy", client, payloadPath, 2)
		updateAdvancedSettingsAttackPayloadLoggingJSON := testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingPolicy.json")
		updateAttackPayloadLoggingRequest := appsec.UpdateAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", JSONPayloadRaw: updateAdvancedSettingsAttackPayloadLoggingJSON}

		updatePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingPolicy.json"
		updateAttackPayloadLogging(t, updateAttackPayloadLoggingRequest, client, updatePayloadPath, 1)

		removeAttackPayloadLoggingRequest := appsec.RemoveAdvancedSettingsAttackPayloadLoggingRequest{
			ConfigID:     43253,
			Version:      7,
			PolicyID:     "test_policy",
			Enabled:      true,
			RequestBody:  appsec.AttackPayloadLoggingRequestBody{Type: appsec.AttackPayload},
			ResponseBody: appsec.AttackPayloadLoggingResponseBody{Type: appsec.AttackPayload},
		}

		removePayloadPath := "testdata/TestResAdvancedSettingsAttackPayloadLogging/UpdateAdvancedSettingsAttackPayloadLoggingPolicy.json"
		removeAttackPayloadLogging(t, removeAttackPayloadLoggingRequest, client, removePayloadPath, 1)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsAttackPayloadLogging/update_by_policy_id.tf"),
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
