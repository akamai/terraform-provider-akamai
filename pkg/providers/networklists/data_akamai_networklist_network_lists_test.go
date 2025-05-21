package networklists

import (
	"encoding/json"
	"testing"

	network "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccAkamaiNetworkList_data_basic(t *testing.T) {
	t.Run("match by NetworkList ID", func(t *testing.T) {
		client := &network.Mock{}

		networkListsResponse := network.GetNetworkListsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSNetworkList/NetworkList.json"), &networkListsResponse)
		require.NoError(t, err)

		client.On("GetNetworkLists",
			testutils.MockContext,
			network.GetNetworkListsRequest{Name: "40996_ARTYLABWHITELIST", Type: "IP"},
		).Return(&networkListsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSNetworkList/match_by_id.tf"),
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
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSNetworkList/SingleNetworkList.json"), &networkListResponse)
		require.NoError(t, err)

		client.On("GetNetworkList",
			testutils.MockContext,
			network.GetNetworkListRequest{UniqueID: "86093_AGEOLIST"},
		).Return(&networkListResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSNetworkList/match_by_uniqueid.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_networklist_network_lists.test", "id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("data.akamai_networklist_network_lists.test", "contract_id", "3-4168BG"),
							resource.TestCheckResourceAttr("data.akamai_networklist_network_lists.test", "group_id", "17240"),
							resource.TestCheckResourceAttr("data.akamai_networklist_network_lists.test", "sync_point", "1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
