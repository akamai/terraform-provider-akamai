package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalPenaltyBox_res_basic(t *testing.T) {
	t.Run("match by EvalPenaltyBox ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updResp := appsec.UpdatePenaltyBoxResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalPenaltyBox/PenaltyBox.json"), &updResp)
		require.NoError(t, err)

		getResp := appsec.GetPenaltyBoxResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalPenaltyBox/PenaltyBox.json"), &getResp)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEvalPenaltyBox",
			testutils.MockContext,
			appsec.GetPenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResp, nil)

		client.On("UpdateEvalPenaltyBox",
			testutils.MockContext,
			appsec.UpdatePenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "none", PenaltyBoxProtection: false},
		).Return(&updResp, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
