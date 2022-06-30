package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiEval_res_basic(t *testing.T) {
	t.Run("match by Eval ID", func(t *testing.T) {
		client := &mockappsec{}

		updateEvalResponse := appsec.UpdateEvalResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResEval/EvalStart.json"), &updateEvalResponse)

		getEvalResponse := appsec.GetEvalResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResEval/EvalStart.json"), &getEvalResponse)

		removeEvalResponse := appsec.RemoveEvalResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResEval/EvalStop.json"), &removeEvalResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

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
