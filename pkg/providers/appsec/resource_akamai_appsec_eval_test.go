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

		cu := appsec.UpdateEvalResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResEval/EvalStart.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetEvalResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResEval/EvalStart.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveEvalResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResEval/EvalStop.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		client.On("GetEval",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetEvalRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Current: "", Eval: "START"},
		).Return(&cr, nil)

		client.On("UpdateEval",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateEvalRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Current: "", Eval: "START"},
		).Return(&cu, nil)

		client.On("RemoveEval",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveEvalRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Current: "", Eval: "STOP"},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEval/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval.test", "id", "43253"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResEval/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
