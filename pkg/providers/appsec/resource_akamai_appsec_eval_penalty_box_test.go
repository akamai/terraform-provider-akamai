package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalPenaltyBox_res_basic(t *testing.T) {
	t.Run("match by EvalPenaltyBox ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updResp := appsec.UpdatePenaltyBoxResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResEvalPenaltyBox/PenaltyBox.json"), &updResp)
		require.NoError(t, err)

		getResp := appsec.GetPenaltyBoxResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResEvalPenaltyBox/PenaltyBox.json"), &getResp)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEvalPenaltyBox",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetPenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResp, nil)

		client.On("UpdateEvalPenaltyBox",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdatePenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "none", PenaltyBoxProtection: false},
		).Return(&updResp, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEvalPenaltyBox/match_by_id.tf"),
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
