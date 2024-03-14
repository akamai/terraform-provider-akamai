package appsec

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalPenaltyBoxConditions_data_basic(t *testing.T) {
	t.Run("match by EvalPenaltyBoxCondition ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		penaltyBoxCondition := appsec.GetPenaltyBoxConditionsResponse{}
		penaltyBoxConditionBytes := testutils.LoadFixtureBytes(t, "testdata/TestDSEvalPenaltyBoxConditions/PenaltyBoxConditions.json")
		var penaltyBoxConditionsJSON bytes.Buffer
		err = json.Compact(&penaltyBoxConditionsJSON, []byte(penaltyBoxConditionBytes))
		require.NoError(t, err)
		err = json.Unmarshal(penaltyBoxConditionBytes, &penaltyBoxCondition)
		require.NoError(t, err)

		expectedOutputText := "\n+---------------------------------+\n| evalPenaltyBoxConditionsDS      |\n+--------------------+------------+\n| CONDITIONSOPERATOR | CONDITIONS |\n+--------------------+------------+\n| AND                | True       |\n+--------------------+------------+\n"
		client.On("GetEvalPenaltyBoxConditions",
			mock.Anything,
			appsec.GetPenaltyBoxConditionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&penaltyBoxCondition, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSEvalPenaltyBoxConditions/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_eval_penalty_box_conditions.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("data.akamai_appsec_eval_penalty_box_conditions.test", "json", penaltyBoxConditionsJSON.String()),
							resource.TestCheckResourceAttr("data.akamai_appsec_eval_penalty_box_conditions.test", "output_text", expectedOutputText),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
