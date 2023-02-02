package networklists

import (
	"encoding/json"
	"testing"

	network "github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/networklists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAccAkamaiNetworkList_data_basic(t *testing.T) {
	t.Run("match by NetworkList ID", func(t *testing.T) {
		client := &network.Mock{}

		networkListsResponse := network.GetNetworkListsResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestDSNetworkList/NetworkList.json"), &networkListsResponse)
		require.NoError(t, err)

		client.On("GetNetworkLists",
			mock.Anything,
			network.GetNetworkListsRequest{Name: "40996_ARTYLABWHITELIST", Type: "IP"},
		).Return(&networkListsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSNetworkList/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_networklist_network_lists.test", "id", "365_AKAMAITOREXITNODES"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestAccAkamaiNetworkList_data_by_uniqueID(t *testing.T) {
	t.Run("match by uniqueID", func(t *testing.T) {
		client := &network.Mock{}

		networkListResponse := network.GetNetworkListResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestDSNetworkList/SingleNetworkList.json"), &networkListResponse)
		require.NoError(t, err)

		client.On("GetNetworkList",
			mock.Anything,
			network.GetNetworkListRequest{UniqueID: "86093_AGEOLIST"},
		).Return(&networkListResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSNetworkList/match_by_uniqueid.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_networklist_network_lists.test", "id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("data.akamai_networklist_network_lists.test", "contract_id", "3-4168BG"),
							resource.TestCheckResourceAttr("data.akamai_networklist_network_lists.test", "group_id", "17240"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
