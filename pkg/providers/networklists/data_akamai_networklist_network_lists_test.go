package networklists

import (
	"encoding/json"
	"testing"

	network "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccAkamaiNetworkList_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by NetworkList ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		networkListsResponse := network.GetNetworkListsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSNetworkList/NetworkList.json"), &networkListsResponse)
		require.NoError(t, err)

		client.NetworkLists.On("GetNetworkLists",
			testutils.MockContext,
			network.GetNetworkListsRequest{Name: "40996_ARTYLABWHITELIST", Type: "IP"},
		).Return(&networkListsResponse, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSNetworkList/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_networklist_network_lists.test", "id", "365_AKAMAITOREXITNODES"),
					),
				},
			},
		})

		client.NetworkLists.AssertExpectations(t)
	})
}

func TestAccAkamaiNetworkList_data_by_uniqueID(t *testing.T) {
	t.Parallel()
	t.Run("match by uniqueID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		networkListResponse := network.GetNetworkListResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSNetworkList/SingleNetworkList.json"), &networkListResponse)
		require.NoError(t, err)

		client.NetworkLists.On("GetNetworkList",
			testutils.MockContext,
			network.GetNetworkListRequest{UniqueID: "86093_AGEOLIST"},
		).Return(&networkListResponse, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
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

		client.NetworkLists.AssertExpectations(t)
	})

}
