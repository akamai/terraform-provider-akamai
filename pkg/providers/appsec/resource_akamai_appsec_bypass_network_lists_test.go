package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiBypassNetworkLists_res_basic(t *testing.T) {
	t.Run("match by BypassNetworkLists ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateBypassNetworkListsResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResBypassNetworkLists/BypassNetworkLists.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetBypassNetworkListsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResBypassNetworkLists/BypassNetworkLists.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveBypassNetworkListsResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResBypassNetworkLists/BypassNetworkLists.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		client.On("GetBypassNetworkLists",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetBypassNetworkListsRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateBypassNetworkLists",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateBypassNetworkListsRequest{ConfigID: 43253, Version: 7, NetworkLists: []string{"888518_ACDDCKERS", "1304427_AAXXBBLIST"}},
		).Return(&cu, nil)

		client.On("RemoveBypassNetworkLists",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveBypassNetworkListsRequest{ConfigID: 43253, Version: 7, NetworkLists: []string(nil)},
		).Return(&crd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResBypassNetworkLists/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_bypass_network_lists.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
