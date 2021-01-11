package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiApiRequestConstraints_res_basic(t *testing.T) {
	t.Run("match by ApiRequestConstraints ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateApiRequestConstraintsResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResApiRequestConstraints/ApiRequestConstraints.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetApiRequestConstraintsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResApiRequestConstraints/ApiRequestConstraints.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetApiRequestConstraints",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetApiRequestConstraintsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdateApiRequestConstraints",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateApiRequestConstraintsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApiID: 1, Action: "alert"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResApiRequestConstraints/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_api_request_constraints.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
