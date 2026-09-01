package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalPenaltyBoxConditions_res_basic(t *testing.T) {
	t.Parallel()
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

		evalPenaltyBoxConditionsRead = func(configId int, version int, policyId string, client *appsec.Mock, path string) {
			evalPenaltyBoxConditionsResponse := appsec.GetPenaltyBoxConditionsResponse{}
			err := json.Unmarshal(testutils.LoadFixtureBytes(t, path), &evalPenaltyBoxConditionsResponse)
			require.NoError(t, err)

			client.On("GetEvalPenaltyBoxConditions",
				testutils.MockContext,
				appsec.GetPenaltyBoxConditionsRequest{ConfigID: configId, Version: version, PolicyID: policyId},
			).Return(&evalPenaltyBoxConditionsResponse, nil)
		}

		evalPenaltyBoxConditionsUpdate = func(evalPenaltyBoxConditionsUpdateReq appsec.UpdatePenaltyBoxConditionsRequest, client *appsec.Mock) {
			evalPenaltyBoxConditionsResponse := appsec.UpdatePenaltyBoxConditionsResponse{}

			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalPenaltyBoxConditions/PenaltyBoxConditions.json"), &evalPenaltyBoxConditionsResponse)
			require.NoError(t, err)

			client.On("UpdateEvalPenaltyBoxConditions",
				testutils.MockContext,
				evalPenaltyBoxConditionsUpdateReq,
			).Return(&evalPenaltyBoxConditionsResponse, nil).Once()
		}

		evalPenaltyBoxConditionsDelete = func(evalPenaltyBoxConditionsUpdateReq appsec.UpdatePenaltyBoxConditionsRequest, client *appsec.Mock) {
			evalPenaltyBoxConditionsDeleteResponse := appsec.UpdatePenaltyBoxConditionsResponse{}

			err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalPenaltyBoxConditions/PenaltyBoxConditionsEmpty.json"), &evalPenaltyBoxConditionsDeleteResponse)
			require.NoError(t, err)

			client.On("UpdateEvalPenaltyBoxConditions",
				testutils.MockContext,
				evalPenaltyBoxConditionsUpdateReq,
			).Return(&evalPenaltyBoxConditionsDeleteResponse, nil)
		}
	)

	t.Run("match by EvalPenaltyBoxConditions ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		configResponse := configVersion(43253, client.APPSEC)

		// eval penalty box condition read test
		evalPenaltyBoxConditionsRead(43253, 7, "AAAA_81230", client.APPSEC, "testdata/TestResEvalPenaltyBoxConditions/PenaltyBoxConditions.json")

		// eval Penalty Box conditions update test
		evalPenaltyBoxConditionsUpdateReq := appsec.PenaltyBoxConditionsPayload{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalPenaltyBoxConditions/PenaltyBoxConditions.json"), &evalPenaltyBoxConditionsUpdateReq)
		require.NoError(t, err)

		updatePenaltyBoxConditionsReq := appsec.UpdatePenaltyBoxConditionsRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "AAAA_81230", ConditionsPayload: evalPenaltyBoxConditionsUpdateReq}
		evalPenaltyBoxConditionsUpdate(updatePenaltyBoxConditionsReq, client.APPSEC)

		// eval Penalty box conditions delete test
		evalPenaltyBoxConditionsDeleteReq := appsec.PenaltyBoxConditionsPayload{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalPenaltyBoxConditions/PenaltyBoxConditionsEmpty.json"), &evalPenaltyBoxConditionsDeleteReq)
		require.NoError(t, err)

		removeEvalPenaltyBoxConditionsReq := appsec.UpdatePenaltyBoxConditionsRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "AAAA_81230", ConditionsPayload: evalPenaltyBoxConditionsDeleteReq}
		evalPenaltyBoxConditionsDelete(removeEvalPenaltyBoxConditionsReq, client.APPSEC)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBoxConditions/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box_conditions.test", "id", "43253:AAAA_81230"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

	t.Run("match by EvalPenaltyBoxConditions ID for Delete case", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()
		configResponse := configVersion(43253, client.APPSEC)

		// eval penalty box condition read test
		evalPenaltyBoxConditionsRead(43253, 7, "AAAA", client.APPSEC, "testdata/TestResEvalPenaltyBoxConditions/PenaltyBoxConditionsEmpty.json")

		// eval Penalty box conditions delete test
		evalPenaltyBoxConditionsDeleteReq := appsec.PenaltyBoxConditionsPayload{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata//TestResEvalPenaltyBoxConditions/PenaltyBoxConditionsEmpty.json"), &evalPenaltyBoxConditionsDeleteReq)
		require.NoError(t, err)

		removeEvalPenaltyBoxConditionsReq := appsec.UpdatePenaltyBoxConditionsRequest{ConfigID: configResponse.ID, Version: configResponse.LatestVersion, PolicyID: "AAAA", ConditionsPayload: evalPenaltyBoxConditionsDeleteReq}
		evalPenaltyBoxConditionsDelete(removeEvalPenaltyBoxConditionsReq, client.APPSEC)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBoxConditions/match_by_id_for_delete.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box_conditions.delete_condition", "id", "43253:AAAA"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})
}
