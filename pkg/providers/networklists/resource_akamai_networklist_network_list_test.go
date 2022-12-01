package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/networklists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiNetworkList_res_basic(t *testing.T) {
	t.Run("match by NetworkList ID", func(t *testing.T) {
		client := &networklists.Mock{}

		createResponse := networklists.CreateNetworkListResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResNetworkList/NetworkList.json"), &createResponse)

		crl := networklists.GetNetworkListsResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResNetworkList/NetworkLists.json"), &crl)

		getResponse := networklists.GetNetworkListResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResNetworkList/NetworkList.json"), &getResponse)

		updateResponse := networklists.UpdateNetworkListResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResNetworkList/NetworkListUpdated.json"), &updateResponse)

		getResponseAfterUpdate := networklists.GetNetworkListResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResNetworkList/NetworkListUpdated.json"), &getResponseAfterUpdate)

		cd := networklists.RemoveNetworkListResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResNetworkList/empty.json"), &cd)

		client.On("CreateNetworkList",
			mock.Anything,
			networklists.CreateNetworkListRequest{Name: "Voyager Call Center Whitelist", Type: "IP", Description: "Notes about this network list", List: []string{"10.1.8.23", "10.3.5.67"}, ContractID: "C-1FRYVV3", GroupID: 64867},
		).Return(&createResponse, nil)

		client.On("GetNetworkLists",
			mock.Anything,
			networklists.GetNetworkListsRequest{Name: "Voyager Call Center Whitelist", Type: "IP"},
		).Return(&crl, nil)

		client.On("GetNetworkList",
			mock.Anything,
			networklists.GetNetworkListRequest{UniqueID: "2275_VOYAGERCALLCENTERWHITELI"},
		).Return(&getResponse, nil).Times(3)

		client.On("UpdateNetworkList",
			mock.Anything,
			networklists.UpdateNetworkListRequest{Name: "Voyager Call Center Whitelist", Type: "IP", Description: "New notes about this network list", SyncPoint: 0, List: []string{"10.1.8.23", "10.3.5.67"}, UniqueID: "2275_VOYAGERCALLCENTERWHITELI", ContractID: "C-1FRYVV3", GroupID: 64867},
		).Return(&updateResponse, nil)

		client.On("GetNetworkList",
			mock.Anything,
			networklists.GetNetworkListRequest{UniqueID: "2275_VOYAGERCALLCENTERWHITELI"},
		).Return(&getResponseAfterUpdate, nil).Times(3)

		client.On("RemoveNetworkList",
			mock.Anything,
			networklists.RemoveNetworkListRequest{UniqueID: "2275_VOYAGERCALLCENTERWHITELI"},
		).Return(&cd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResNetworkList/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "name", "Voyager Call Center Whitelist"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResNetworkList/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "name", "Voyager Call Center Whitelist"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
