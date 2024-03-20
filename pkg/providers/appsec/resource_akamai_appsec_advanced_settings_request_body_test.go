package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsRequestBodyResConfig(t *testing.T) {
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

		requestBodyRead = func(t *testing.T, configId int, version int, policyId string, client *appsec.Mock, numberOfTimes int) {
			requestBodyResponse := appsec.GetAdvancedSettingsRequestBodyResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json"), &requestBodyResponse)
			require.NoError(t, err)

			client.On("GetAdvancedSettingsRequestBody",
				mock.Anything,
				appsec.GetAdvancedSettingsRequestBodyRequest{ConfigID: configId, Version: version, PolicyID: policyId},
			).Return(&requestBodyResponse, nil).Times(numberOfTimes)

		}

		updateRequestBody = func(t *testing.T, updateRequestBody appsec.UpdateAdvancedSettingsRequestBodyRequest, client *appsec.Mock, numberOfTimes int) {
			updateRequestBodyResponse := appsec.UpdateAdvancedSettingsRequestBodyResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json"), &updateRequestBodyResponse)
			require.NoError(t, err)

			client.On("UpdateAdvancedSettingsRequestBody",
				mock.Anything, updateRequestBody,
			).Return(&updateRequestBodyResponse, nil).Times(numberOfTimes)

		}

		removeRequestBody = func(t *testing.T, updateRequestBody appsec.RemoveAdvancedSettingsRequestBodyRequest, client *appsec.Mock, numberOfTimes int) {
			removeRequestBodyResponse := appsec.RemoveAdvancedSettingsRequestBodyResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json"), &removeRequestBodyResponse)
			require.NoError(t, err)

			client.On("RemoveAdvancedSettingsRequestBody",
				mock.Anything, updateRequestBody,
			).Return(&removeRequestBodyResponse, nil).Times(numberOfTimes)

		}
	)

	t.Run("match by AdvancedSettingsRequestBody ID", func(t *testing.T) {
		client := &appsec.Mock{}
		configResponse := configVersion(t, 43253, client)

		requestBodyRead(t, 43253, 7, "", client, 2)
		updateRequestBodyRequest := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: 7, PolicyID: "", RequestBodyInspectionLimitInKB: appsec.Limit16KB}

		updateRequestBody(t, updateRequestBodyRequest, client, 1)

		removeRequestBodyRequest := appsec.RemoveAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: 7, PolicyID: "", RequestBodyInspectionLimitInKB: appsec.Default}

		removeRequestBody(t, removeRequestBodyRequest, client, 1)
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsRequestBody/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_request_body.test", "id", "43253:"),
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

		requestBodyRead(t, configResponse.ID, configResponse.LatestVersion, "", client, 4)
		updateRequestBodyRequest := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "", RequestBodyInspectionLimitInKB: appsec.Limit16KB}

		updateRequestBody(t, updateRequestBodyRequest, client, 1)

		removeRequestBodyRequest := appsec.RemoveAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "", RequestBodyInspectionLimitInKB: appsec.Default}

		removeRequestBody(t, removeRequestBodyRequest, client, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsRequestBody/match_by_id.tf"),
					},
					{
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateId:     "43253",
						ResourceName:      "akamai_appsec_advanced_settings_request_body.test",
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("import policy", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := configVersion(t, 43253, client)

		requestBodyRead(t, configResponse.ID, configResponse.LatestVersion, "test_policy", client, 4)

		updateRequestBodyRequest := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Limit16KB}

		updateRequestBody(t, updateRequestBodyRequest, client, 1)

		removeRequestBodyRequest := appsec.RemoveAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Default}

		removeRequestBody(t, removeRequestBodyRequest, client, 1)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsRequestBody/update_by_policy_id.tf"),
					},
					{
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateId:     "43253:test_policy",
						ResourceName:      "akamai_appsec_advanced_settings_request_body.policy",
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("match by AdvancedSettingsRequestBodyPolicy ID", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := configVersion(t, 43253, client)

		requestBodyRead(t, configResponse.ID, configResponse.LatestVersion, "test_policy", client, 2)

		updateRequestBodyRequest := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Limit16KB}

		updateRequestBody(t, updateRequestBodyRequest, client, 1)

		removeRequestBodyRequest := appsec.RemoveAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Default}

		removeRequestBody(t, removeRequestBodyRequest, client, 1)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsRequestBody/update_by_policy_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_request_body.policy", "id", "43253:test_policy"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
