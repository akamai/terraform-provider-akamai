package networklists

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAccAkamaiNetworkList_res_basic(t *testing.T) {

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
		networklists.CreateNetworkListRequest{Name: "Voyager Call Center Whitelist", Type: "IP", Description: "Notes about this network list", List: []string{"10.1.8.23", "10.3.5.67"}},
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

	t.Run("match by NetworkList ID", func(t *testing.T) {
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResNetworkList/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "name", "Voyager Call Center Whitelist"),
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "group_id", "64867"),
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "contract_id", "C-1FRYVV3"),
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

// Negative test case for checking changed groupID and contractID in config
func TestAccAkamaiNetworkListConfigChanged(t *testing.T) {

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
		networklists.CreateNetworkListRequest{Name: "Voyager Call Center Whitelist", Type: "IP", Description: "Notes about this network list", List: []string{"10.1.8.23", "10.3.5.67"}},
	).Return(&createResponse, nil)

	client.On("GetNetworkLists",
		mock.Anything,
		networklists.GetNetworkListsRequest{Name: "Voyager Call Center Whitelist", Type: "IP"},
	).Return(&crl, nil)

	client.On("GetNetworkList",
		mock.Anything,
		networklists.GetNetworkListRequest{UniqueID: "2275_VOYAGERCALLCENTERWHITELI"},
	).Return(&getResponse, nil).Times(4)

	client.On("RemoveNetworkList",
		mock.Anything,
		networklists.RemoveNetworkListRequest{UniqueID: "2275_VOYAGERCALLCENTERWHITELI"},
	).Return(&cd, nil)

	t.Run("changed contractID and groupID", func(t *testing.T) {
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResNetworkList/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "name", "Voyager Call Center Whitelist"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResNetworkList/changed_contract_id.tf"),
						ExpectError: regexp.MustCompile("contract_id value C-1FRYVV5 specified in configuration differs from resource ID's value C-1FRYVV3"),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResNetworkList/changed_group_id.tf"),
						ExpectError: regexp.MustCompile("group_id value 64865 specified in configuration differs from resource ID's value 64867"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
