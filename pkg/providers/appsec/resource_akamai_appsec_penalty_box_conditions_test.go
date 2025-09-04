package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiPenaltyBoxConditions_res_basic(t *testing.T) {
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

		penaltyBoxConditionsRead = func(configId int, version int, policyId string, client *appsec.Mock, path string) {
			penaltyBoxConditionsResponse := appsec.GetPenaltyBoxConditionsResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, path), &penaltyBoxConditionsResponse)
			require.NoError(t, err)

			client.On("GetPenaltyBoxConditions",
				testutils.MockContext,
				appsec.GetPenaltyBoxConditionsRequest{ConfigID: configId, Version: version, PolicyID: policyId},
			).Return(&penaltyBoxConditionsResponse, nil)
		}

		penaltyBoxConditionsUpdate = func(penaltyBoxConditionsUpdateReq appsec.UpdatePenaltyBoxConditionsRequest, client *appsec.Mock) {
			penaltyBoxConditionsResponse := appsec.UpdatePenaltyBoxConditionsResponse{}

			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResPenaltyBoxConditions/PenaltyBoxConditions.json"), &penaltyBoxConditionsResponse)
			require.NoError(t, err)

			client.On("UpdatePenaltyBoxConditions",
				testutils.MockContext,
				penaltyBoxConditionsUpdateReq,
			).Return(&penaltyBoxConditionsResponse, nil).Once()
		}

		penaltyBoxConditionsDelete = func(penaltyBoxConditionsUpdateReq appsec.UpdatePenaltyBoxConditionsRequest, client *appsec.Mock) {
			penaltyBoxConditionsDeleteResponse := appsec.UpdatePenaltyBoxConditionsResponse{}

			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResPenaltyBoxConditions/PenaltyBoxConditionsEmpty.json"), &penaltyBoxConditionsDeleteResponse)
			require.NoError(t, err)

			client.On("UpdatePenaltyBoxConditions",
				testutils.MockContext,
				penaltyBoxConditionsUpdateReq,
			).Return(&penaltyBoxConditionsDeleteResponse, nil)
		}
	)

	t.Run("match by PenaltyBoxConditions ID", func(t *testing.T) {
		client := &appsec.Mock{}
		configResponse := configVersion(43253, client)

		penaltyBoxConditionsRead(43253, 7, "AAAA_81230", client, "testdata/TestResPenaltyBoxConditions/PenaltyBoxConditions.json")

		penaltyBoxConditionsUpdateReq := appsec.PenaltyBoxConditionsPayload{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResPenaltyBoxConditions/PenaltyBoxConditions.json"), &penaltyBoxConditionsUpdateReq)
		require.NoError(t, err)

		updatePenaltyBoxConditionsReq := appsec.UpdatePenaltyBoxConditionsRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "AAAA_81230", ConditionsPayload: penaltyBoxConditionsUpdateReq}
		penaltyBoxConditionsUpdate(updatePenaltyBoxConditionsReq, client)

		penaltyBoxConditionsDeleteReq := appsec.PenaltyBoxConditionsPayload{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResPenaltyBoxConditions/PenaltyBoxConditionsEmpty.json"), &penaltyBoxConditionsDeleteReq)
		require.NoError(t, err)

		removePenaltyBoxConditionsReq := appsec.UpdatePenaltyBoxConditionsRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "AAAA_81230", ConditionsPayload: penaltyBoxConditionsDeleteReq}
		penaltyBoxConditionsDelete(removePenaltyBoxConditionsReq, client)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResPenaltyBoxConditions/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_penalty_box_conditions.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by PenaltyBoxConditions ID for Delete case", func(t *testing.T) {
		client := &appsec.Mock{}
		configResponse := configVersion(43253, client)

		penaltyBoxConditionsRead(43253, 7, "AAAA", client, "testdata/TestResPenaltyBoxConditions/PenaltyBoxConditionsEmpty.json")

		penaltyBoxConditionsDeleteReq := appsec.PenaltyBoxConditionsPayload{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResPenaltyBoxConditions/PenaltyBoxConditionsEmpty.json"), &penaltyBoxConditionsDeleteReq)
		require.NoError(t, err)

		removePenaltyBoxConditionsReq := appsec.UpdatePenaltyBoxConditionsRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "AAAA", ConditionsPayload: penaltyBoxConditionsDeleteReq}
		penaltyBoxConditionsDelete(removePenaltyBoxConditionsReq, client)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResPenaltyBoxConditions/match_by_id_for_delete.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_penalty_box_conditions.delete_condition", "id", "43253:AAAA"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
