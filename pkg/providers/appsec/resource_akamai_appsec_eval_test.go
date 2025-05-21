package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEval_res_basic(t *testing.T) {
	t.Run("match by Eval ID", func(t *testing.T) {
		client := &appsec.Mock{}

		updateEvalResponse := appsec.UpdateEvalResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEval/EvalStart.json"), &updateEvalResponse)
		require.NoError(t, err)

		getEvalResponse := appsec.GetEvalResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEval/EvalStart.json"), &getEvalResponse)
		require.NoError(t, err)

		removeEvalResponse := appsec.RemoveEvalResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEval/EvalStop.json"), &removeEvalResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEval",
			testutils.MockContext,
			appsec.GetEvalRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Current: "", Eval: ""},
		).Return(&getEvalResponse, nil)

		client.On("UpdateEval",
			testutils.MockContext,
			appsec.UpdateEvalRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Current: "", Eval: "START"},
		).Return(&updateEvalResponse, nil)

		client.On("RemoveEval",
			testutils.MockContext,
			appsec.RemoveEvalRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Current: "", Eval: "STOP"},
		).Return(&removeEvalResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResEval/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval.test", "id", "43253:AAAA_81230"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResEval/update_by_id.tf"),
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
