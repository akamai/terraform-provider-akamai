package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsRequestBodyResConfig(t *testing.T) {
	var (
		configVersion = func(configId int, client *appsec.Mock) appsec.GetConfigurationResponse {
			configResponse := appsec.GetConfigurationResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &configResponse)
			require.NoError(t, err)

			client.On("GetConfiguration",
				testutils.MockContext,
				appsec.GetConfigurationRequest{ConfigID: configId},
			).Return(&configResponse, nil)

			return configResponse
		}

		requestBodyRead = func(configId int, version int, policyId string, client *appsec.Mock, numberOfTimes int, filePath string) {
			requestBodyResponse := appsec.GetAdvancedSettingsRequestBodyResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, filePath), &requestBodyResponse)
			require.NoError(t, err)

			client.On("GetAdvancedSettingsRequestBody",
				testutils.MockContext,
				appsec.GetAdvancedSettingsRequestBodyRequest{ConfigID: configId, Version: version, PolicyID: policyId},
			).Return(&requestBodyResponse, nil).Times(numberOfTimes)

		}

		updateRequestBody = func(updateRequestBody appsec.UpdateAdvancedSettingsRequestBodyRequest, client *appsec.Mock, numberOfTimes int, filePath string) {
			updateRequestBodyResponse := appsec.UpdateAdvancedSettingsRequestBodyResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, filePath), &updateRequestBodyResponse)
			require.NoError(t, err)

			client.On("UpdateAdvancedSettingsRequestBody",
				testutils.MockContext, updateRequestBody,
			).Return(&updateRequestBodyResponse, nil).Times(numberOfTimes)

		}

		removeRequestBody = func(updateRequestBody appsec.RemoveAdvancedSettingsRequestBodyRequest, client *appsec.Mock, numberOfTimes int, filePath string) {
			removeRequestBodyResponse := appsec.RemoveAdvancedSettingsRequestBodyResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, filePath), &removeRequestBodyResponse)
			require.NoError(t, err)

			client.On("RemoveAdvancedSettingsRequestBody",
				testutils.MockContext, updateRequestBody,
			).Return(&removeRequestBodyResponse, nil).Times(numberOfTimes)

		}
	)

	t.Run("match by AdvancedSettingsRequestBody ID", func(t *testing.T) {
		client := &appsec.Mock{}
		configResponse := configVersion(43253, client)

		requestBodyRead(43253, 7, "", client, 2, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json")

		updateRequestBodyRequest := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: 7, PolicyID: "", RequestBodyInspectionLimitInKB: appsec.Limit16KB}

		updateRequestBody(updateRequestBodyRequest, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json")
		removeRequestBodyRequest := appsec.RemoveAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: 7, PolicyID: "", RequestBodyInspectionLimitInKB: appsec.Default}

		removeRequestBody(removeRequestBodyRequest, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json")
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

		configResponse := configVersion(43253, client)

		requestBodyRead(configResponse.ID, configResponse.LatestVersion, "", client, 4, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json")

		updateRequestBodyRequest := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "", RequestBodyInspectionLimitInKB: appsec.Limit16KB}

		updateRequestBody(updateRequestBodyRequest, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json")

		removeRequestBodyRequest := appsec.RemoveAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "", RequestBodyInspectionLimitInKB: appsec.Default}

		removeRequestBody(removeRequestBodyRequest, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json")

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsRequestBody/match_by_id.tf"),
					},
					{
						ImportState:             true,
						ImportStateVerify:       true,
						ImportStateVerifyIgnore: []string{"request_body_inspection_limit_override"},
						ImportStateId:           "43253",
						ResourceName:            "akamai_appsec_advanced_settings_request_body.test",
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
	t.Run("import policy", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := configVersion(43253, client)

		requestBodyRead(configResponse.ID, configResponse.LatestVersion, "test_policy", client, 4, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicy.json")

		updateRequestBodyRequest := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Limit16KB, RequestBodyInspectionLimitOverride: true}

		updateRequestBody(updateRequestBodyRequest, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json")

		removeRequestBodyRequest := appsec.RemoveAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Default, RequestBodyInspectionLimitOverride: false}

		removeRequestBody(removeRequestBodyRequest, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyDisabled.json")

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

		configResponse := configVersion(43253, client)

		requestBodyRead(configResponse.ID, configResponse.LatestVersion, "test_policy", client, 5, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicy.json")

		updateRequestBodyRequest := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Limit16KB, RequestBodyInspectionLimitOverride: true}

		updateRequestBody(updateRequestBodyRequest, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody.json")

		updateRequestBodyRequestWithVal := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Limit32KB, RequestBodyInspectionLimitOverride: true}

		updateRequestBody(updateRequestBodyRequestWithVal, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBody32.json")

		removeRequestBodyRequest := appsec.RemoveAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Default, RequestBodyInspectionLimitOverride: false}

		removeRequestBody(removeRequestBodyRequest, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyDisabled.json")

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
					{
						Config:             testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsRequestBody/update_by_policy_32.tf"),
						ExpectNonEmptyPlan: true,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_request_body.policy", "id", "43253:test_policy"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
	t.Run("match by AdvancedSettingsRequestBodyPolicyIDDisable", func(t *testing.T) {
		client := &appsec.Mock{}

		configResponse := configVersion(43253, client)

		requestBodyRead(configResponse.ID, configResponse.LatestVersion, "test_policy", client, 5, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicy.json")

		// create
		updateRequestBodyRequest := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Limit16KB, RequestBodyInspectionLimitOverride: true}

		updateRequestBody(updateRequestBodyRequest, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyPolicy.json")

		//update
		updateRequestBodyRequestDisable := appsec.UpdateAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Limit32KB, RequestBodyInspectionLimitOverride: false}

		updateRequestBody(updateRequestBodyRequestDisable, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyDisabled.json")

		//delete
		removeRequestBodyRequest := appsec.RemoveAdvancedSettingsRequestBodyRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "test_policy", RequestBodyInspectionLimitInKB: appsec.Default, RequestBodyInspectionLimitOverride: false}

		removeRequestBody(removeRequestBodyRequest, client, 1, "testdata/TestResAdvancedSettingsRequestBody/AdvancedSettingsRequestBodyDisabled.json")

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsRequestBody/update_by_policy_id.tf"),
					},
					{
						Config:             testutils.LoadFixtureString(t, "testdata/TestResAdvancedSettingsRequestBody/update_by_policy_id_disable.tf"),
						ExpectNonEmptyPlan: true,
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
