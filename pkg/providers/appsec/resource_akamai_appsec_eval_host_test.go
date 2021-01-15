package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiEvalHost_res_basic(t *testing.T) {
	t.Run("match by EvalHost ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateEvalHostResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResEvalHost/EvalHost.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetEvalHostResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResEvalHost/EvalHost.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveEvalHostResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResEvalHost/EvalHost.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		client.On("GetEvalHost",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetEvalHostRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateEvalHost",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateEvalHostRequest{ConfigID: 43253, Version: 7, Hostnames: []string{"example.com"}},
		).Return(&cu, nil)

		client.On("RemoveEvalHost",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveEvalHostRequest{ConfigID: 43253, Version: 7, Hostnames: []string(nil)},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEvalHost/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval_hostnames.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
