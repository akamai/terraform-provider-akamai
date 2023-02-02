package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEval_res_basic(t *testing.T) {
	t.Run("match by Eval ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateEvalResponse := appsec.UpdateEvalResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResEval/EvalStart.json"), &updateEvalResponse)
		require.NoError(t, err)

		getEvalResponse := appsec.GetEvalResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResEval/EvalStart.json"), &getEvalResponse)
		require.NoError(t, err)

		removeEvalResponse := appsec.RemoveEvalResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResEval/EvalStop.json"), &removeEvalResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEval",
			mock.Anything,
			appsec.GetEvalRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Current: "", Eval: ""},
		).Return(&getEvalResponse, nil)

		client.On("UpdateEval",
			mock.Anything,
			appsec.UpdateEvalRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Current: "", Eval: "START"},
		).Return(&updateEvalResponse, nil)

		client.On("RemoveEval",
			mock.Anything,
			appsec.RemoveEvalRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Current: "", Eval: "STOP"},
		).Return(&removeEvalResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEval/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval.test", "id", "43253:AAAA_81230"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResEval/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
