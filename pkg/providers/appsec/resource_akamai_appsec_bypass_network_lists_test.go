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

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		cu := appsec.UpdateWAPBypassNetworkListsResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResBypassNetworkLists/BypassNetworkLists.json")), &cu)

		cr := appsec.GetWAPBypassNetworkListsResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResBypassNetworkLists/BypassNetworkLists.json")), &cr)

		crd := appsec.RemoveWAPBypassNetworkListsResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResBypassNetworkLists/BypassNetworkLists.json")), &crd)

		client.On("GetWAPBypassNetworkLists",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetWAPBypassNetworkListsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdateWAPBypassNetworkLists",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateWAPBypassNetworkListsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", NetworkLists: []string{"1304427_AAXXBBLIST", "888518_ACDDCKERS"}},
		).Return(&cu, nil)

		client.On("RemoveWAPBypassNetworkLists",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveWAPBypassNetworkListsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", NetworkLists: []string{}},
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
