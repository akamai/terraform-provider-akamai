package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAccAkamaiNetworkList_res_basic(t *testing.T) {
	t.Run("match by NetworkList ID", func(t *testing.T) {
		client := &networklists.Mock{}

		createResponse := networklists.CreateNetworkListResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkList/NetworkList.json"), &createResponse)
		require.NoError(t, err)

		crl := networklists.GetNetworkListsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkList/NetworkLists.json"), &crl)
		require.NoError(t, err)

		getResponse := networklists.GetNetworkListResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkList/NetworkList.json"), &getResponse)
		require.NoError(t, err)

		updateResponse := networklists.UpdateNetworkListResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkList/NetworkListUpdated.json"), &updateResponse)
		require.NoError(t, err)

		getResponseAfterUpdate := networklists.GetNetworkListResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkList/NetworkListUpdated.json"), &getResponseAfterUpdate)
		require.NoError(t, err)

		cd := networklists.RemoveNetworkListResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkList/empty.json"), &cd)
		require.NoError(t, err)

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
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResNetworkList/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "name", "Voyager Call Center Whitelist"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResNetworkList/update_by_id.tf"),
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
