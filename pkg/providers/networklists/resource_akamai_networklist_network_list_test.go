package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/networklists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiNetworkList_res_basic(t *testing.T) {
	t.Run("match by NetworkList ID", func(t *testing.T) {
		client := &mocknetworklists{}

		cu := networklists.UpdateNetworkListResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/NetworkList.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := networklists.GetNetworkListResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/NetworkList.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetNetworkList",
			mock.Anything, // ctx is irrelevant for this test
			networklists.GetNetworkListRequest{Name: "Test"},
		).Return(&cr, nil)

		client.On("UpdateNetworkList",
			mock.Anything, // ctx is irrelevant for this test
			networklists.UpdateNetworkListRequest{Name: "Test"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResNetworkList/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "name", "Martin Network List"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResNetworkList/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "name", "Martin Network List"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
