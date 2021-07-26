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

		crnl := networklists.CreateNetworkListResponse{}
		expectJSCNL := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/NetworkList.json"))
		json.Unmarshal([]byte(expectJSCNL), &crnl)

		cu := networklists.UpdateNetworkListResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/NetworkLists.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := networklists.GetNetworkListResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/NetworkList.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crl := networklists.GetNetworkListsResponse{}
		expectJSL := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/NetworkLists.json"))
		json.Unmarshal([]byte(expectJSL), &crl)

		cd := networklists.RemoveNetworkListResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/Empty.json"))
		json.Unmarshal([]byte(expectJSD), &cd)

		client.On("CreateNetworkList",
			mock.Anything, // ctx is irrelevant for this test
			networklists.CreateNetworkListRequest{Name: "Voyager Call Center Whitelist", Type: "IP", Description: "Notes about this network list", List: []string{"10.1.8.23", "10.3.5.67"}},
		).Return(&crnl, nil)

		client.On("GetNetworkList",
			mock.Anything, // ctx is irrelevant for this test
			networklists.GetNetworkListRequest{UniqueID: "2275_VOYAGERCALLCENTERWHITELI"},
		).Return(&cr, nil)

		client.On("GetNetworkLists",
			mock.Anything, // ctx is irrelevant for this test
			networklists.GetNetworkListsRequest{Name: "Voyager Call Center Whitelist", Type: "IP"},
		).Return(&crl, nil)

		client.On("UpdateNetworkList",
			mock.Anything, // ctx is irrelevant for this test
			networklists.UpdateNetworkListRequest{Name: "Voyager Call Center Whitelist", Type: "IP", Description: "Notes about this network list", SyncPoint: 0, List: []string{"10.1.8.23", "10.3.5.67"}, UniqueID: "2275_VOYAGERCALLCENTERWHITELI"},
		).Return(&cu, nil)

		client.On("RemoveNetworkList",
			mock.Anything, // ctx is irrelevant for this test
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
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString("testdata/TestResNetworkList/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "name", "Voyager Call Center Whitelist"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAccAkamaiNetworkList_with_ACG(t *testing.T) {
	t.Run("access control group", func(t *testing.T) {
		client := &mocknetworklists{}

		crnl := networklists.CreateNetworkListResponse{}
		expectJSCNL := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/NetworkList.json"))
		json.Unmarshal([]byte(expectJSCNL), &crnl)

		cu := networklists.UpdateNetworkListResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/NetworkLists.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := networklists.GetNetworkListResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/NetworkList.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crl := networklists.GetNetworkListsResponse{}
		expectJSL := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/NetworkLists.json"))
		json.Unmarshal([]byte(expectJSL), &crl)

		cd := networklists.RemoveNetworkListResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResNetworkList/Empty.json"))
		json.Unmarshal([]byte(expectJSD), &cd)

		client.On("CreateNetworkList",
			mock.Anything, // ctx is irrelevant for this test
			networklists.CreateNetworkListRequest{Name: "Voyager Call Center Whitelist", Type: "IP", Description: "Notes about this network list", List: []string{"10.1.8.23", "10.3.5.67"}, AccessControlGroup: " Kona Site Defender - C-D5TW8R - C-D5TW8R.G31325"},
		).Return(&crnl, nil)

		client.On("GetNetworkList",
			mock.Anything, // ctx is irrelevant for this test
			networklists.GetNetworkListRequest{UniqueID: "2275_VOYAGERCALLCENTERWHITELI"},
		).Return(&cr, nil)

		client.On("GetNetworkLists",
			mock.Anything, // ctx is irrelevant for this test
			networklists.GetNetworkListsRequest{Name: "Voyager Call Center Whitelist", Type: "IP"},
		).Return(&crl, nil)

		client.On("UpdateNetworkList",
			mock.Anything, // ctx is irrelevant for this test
			networklists.UpdateNetworkListRequest{Name: "Voyager Call Center Whitelist", Type: "IP", Description: "Notes about this network list", SyncPoint: 0, List: []string{"10.1.8.23", "10.3.5.67"}, UniqueID: "2275_VOYAGERCALLCENTERWHITELI", Group: 31325, Contract: "C-D5TW8R"},
		).Return(&cu, nil)

		client.On("RemoveNetworkList",
			mock.Anything, // ctx is irrelevant for this test
			networklists.RemoveNetworkListRequest{UniqueID: "2275_VOYAGERCALLCENTERWHITELI"},
		).Return(&cd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResNetworkList/access_group.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "name", "Voyager Call Center Whitelist"),
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "access_control_group", " Kona Site Defender - C-D5TW8R - C-D5TW8R.G31325"),
						),
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString("testdata/TestResNetworkList/update_access_group.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "name", "Voyager Call Center Whitelist"),
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "group", "31325"),
							resource.TestCheckResourceAttr("akamai_networklist_network_list.test", "contract", "C-D5TW8R"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
