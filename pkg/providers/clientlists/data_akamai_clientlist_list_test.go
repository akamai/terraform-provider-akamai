package clientlists

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDataClientList(t *testing.T) {
	t.Parallel()

	const testDir = "testData/TestDSClientList/SingleClientList"
	baseChecker := test.NewStateChecker("data.akamai_clientlist_list.list")

	getClientListResponseIP := clientlists.GetClientListResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, getPath(testDir, "ip_client_list.json")), &getClientListResponseIP)
	require.NoError(t, err)

	getClientListResponseGEO := clientlists.GetClientListResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, getPath(testDir, "geo_client_list.json")), &getClientListResponseGEO)
	require.NoError(t, err)

	getClientListResponseASN := clientlists.GetClientListResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, getPath(testDir, "asn_client_list.json")), &getClientListResponseASN)
	require.NoError(t, err)

	getClientListResponseTLS := clientlists.GetClientListResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, getPath(testDir, "tls_fingerprint_client_list.json")), &getClientListResponseTLS)
	require.NoError(t, err)

	getClientListResponseFileHash := clientlists.GetClientListResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, getPath(testDir, "file_hash_client_list.json")), &getClientListResponseFileHash)
	require.NoError(t, err)

	getClientListResponseUsername := &clientlists.GetClientListResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, getPath(testDir, "user_client_list.json")), &getClientListResponseUsername)
	require.NoError(t, err)

	getClientListItemsResponseUsername := &clientlists.GetClientListItemsResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, getPath(testDir, "user_client_list_items.json")), &getClientListItemsResponseUsername)
	require.NoError(t, err)

	getClientListResponseUserID := clientlists.GetClientListResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, getPath(testDir, "user_client_list_user_id.json")), &getClientListResponseUserID)
	require.NoError(t, err)

	getClientListItemsResponseUserID := clientlists.GetClientListItemsResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, getPath(testDir, "user_client_list_items_user_id.json")), &getClientListItemsResponseUserID)
	require.NoError(t, err)

	tests := map[string]struct {
		listType clientlists.ClientListType
		init     func(*clientlists.Mock)
		steps    []resource.TestStep
	}{
		"IP client list": {
			listType: clientlists.IP,
			init: func(m *clientlists.Mock) {
				mockGetClientList(m, getClientListResponseIP, clientlists.GetClientListRequest{
					ListID:       "180991_TESTNLMIGRATION123",
					IncludeItems: true,
				}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(getPath(testDir, "ip_list.tf")),
					Check: baseChecker.
						CheckEqual("list.list_id", "180991_TESTNLMIGRATION123").
						CheckEqual("list.items.#", "4").
						CheckEqual("list.items.0.value", "2.3.4.5").
						CheckEqual("list.items.1.value", "7.7.7.7").
						CheckEqual("list.items.2.value", "8.8.8.8").
						CheckEqual("list.items.3.value", "9.9.9.9").
						CheckEqual("output_text", loadText(t, getPath(testDir, "ip_output_text.txt"))).
						CheckEqual("json", loadJSON(t, getPath(testDir, "ip_output_json.txt"))).
						Build(),
				},
			},
		},
		"GEO client list": {
			listType: clientlists.GEO,
			init: func(m *clientlists.Mock) {
				mockGetClientList(m, getClientListResponseGEO, clientlists.GetClientListRequest{
					ListID:       "115165_PAVITHRALISTGEO",
					IncludeItems: true,
				}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(getPath(testDir, "geo_list.tf")),
					Check: baseChecker.
						CheckEqual("list.list_id", "115165_PAVITHRALISTGEO").
						CheckEqual("list.items.#", "5").
						CheckEqual("list.items.0.value", "AF").
						CheckEqual("list.items.1.value", "AL").
						CheckEqual("list.items.2.value", "AO").
						CheckEqual("list.items.3.value", "DZ").
						CheckEqual("list.items.4.value", "IN").
						CheckEqual("output_text", loadText(t, getPath(testDir, "geo_output_text.txt"))).
						CheckEqual("json", loadJSON(t, getPath(testDir, "geo_output_json.txt"))).
						Build(),
				},
			},
		},
		"ASN client list": {
			listType: clientlists.ASN,
			init: func(m *clientlists.Mock) {
				mockGetClientList(m, getClientListResponseASN, clientlists.GetClientListRequest{
					ListID:       "164730_SECKSD28365ASNKONAQAAA",
					IncludeItems: true,
				}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(getPath(testDir, "asn_list.tf")),
					Check: baseChecker.
						CheckEqual("list.list_id", "164730_SECKSD28365ASNKONAQAAA").
						CheckEqual("list.items.#", "5").
						CheckEqual("list.items.0.value", "1001").
						CheckEqual("list.items.1.value", "2002").
						CheckEqual("list.items.2.value", "3003").
						CheckEqual("list.items.3.value", "4004").
						CheckEqual("list.items.4.value", "5005").
						CheckEqual("output_text", loadText(t, getPath(testDir, "asn_output_text.txt"))).
						CheckEqual("json", loadJSON(t, getPath(testDir, "asn_output_json.txt"))).
						Build(),
				},
			},
		},
		"TLS_FINGERPRINT client list": {
			listType: clientlists.TLSFingerprint,
			init: func(m *clientlists.Mock) {
				mockGetClientList(m, getClientListResponseTLS, clientlists.GetClientListRequest{
					ListID:       "183799_TESTLISTDIPESH",
					IncludeItems: true,
				}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(getPath(testDir, "tls_fingerprint_list.tf")),
					Check: baseChecker.
						CheckEqual("list.list_id", "183799_TESTLISTDIPESH").
						CheckEqual("list.items.#", "2").
						CheckEqual("list.items.0.value", "c18eaddafe6a3bba").
						CheckEqual("list.items.1.value", "cd08e31494f9531f560d64c695473da9").
						CheckEqual("output_text", loadText(t, getPath(testDir, "tls_fingerprint_output_text.txt"))).
						CheckEqual("json", loadJSON(t, getPath(testDir, "tls_fingerprint_output_json.txt"))).
						Build(),
				},
			},
		},
		"FILE_HASH client list": {
			listType: clientlists.FileHash,
			init: func(m *clientlists.Mock) {
				mockGetClientList(m, getClientListResponseFileHash, clientlists.GetClientListRequest{
					ListID:       "164579_SECKSD28365FILEHASHKONA",
					IncludeItems: true,
				}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(fmt.Sprintf("%s/file_hash_list.tf", testDir)),
					Check: baseChecker.
						CheckEqual("list.list_id", "164579_SECKSD28365FILEHASHKONA").
						CheckEqual("list.items.#", "2").
						CheckEqual("list.items.0.value", "65d860160bdc9b98abf72407e14ca40b609417de7939897d3b58d55787aaef69").
						CheckEqual("list.items.1.value", "f0456d7aed088e791e4610c3c2ad63afe46e2e777988fdbc9270f15ec9711b42").
						CheckEqual("output_text", loadText(t, getPath(testDir, "file_hash_output_text.txt"))).
						CheckEqual("json", loadJSON(t, getPath(testDir, "file_hash_output_json.txt"))).
						Build(),
				},
			},
		},
		"USER client list - show usernames enabled": {
			listType: clientlists.USER,
			init: func(m *clientlists.Mock) {
				mockGetClientList(m, *getClientListResponseUsername, clientlists.GetClientListRequest{
					ListID:       "193203_VPUSERTYPECLIENTLISTUS",
					IncludeItems: true,
				}, 3)
				mockGetClientListItems(m, *getClientListItemsResponseUsername, clientlists.GetClientListItemsRequest{
					ListID: "193203_VPUSERTYPECLIENTLISTUS",
				}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(getPath(testDir, "user_list.tf")),
					Check: baseChecker.
						CheckEqual("list.list_id", "193203_VPUSERTYPECLIENTLISTUS").
						CheckEqual("list.items.#", "5").
						CheckEqual("list.items.0.value", "03c9fa81-f3af-4e3f-a15a-93bf49407d5f (user4)").
						CheckEqual("list.items.1.value", "672a43b0-2841-42f1-aeee-e13569c31524 (user1)").
						CheckEqual("list.items.2.value", "7340cbfc-3dad-4c52-a653-7874bc2a79dc (sales@ubs.com)").
						CheckEqual("list.items.3.value", "8894f428-bedc-4ca2-8258-fa24d7740709 (user2)").
						CheckEqual("list.items.4.value", "d7505c0f-b1f6-4c02-93bf-a601577f5641 (user3)").
						CheckEqual("output_text", loadText(t, getPath(testDir, "user_output_text.txt"))).
						CheckEqual("json", loadJSON(t, getPath(testDir, "user_output_json.txt"))).
						Build(),
				},
			},
		},
		"USER client list - show usernames disabled": {
			listType: clientlists.USER,
			init: func(m *clientlists.Mock) {
				mockGetClientList(m, getClientListResponseUserID, clientlists.GetClientListRequest{
					ListID:       "193089_VPTERRAFORM2",
					IncludeItems: true,
				}, 3)
				mockGetClientListItems(m, getClientListItemsResponseUserID, clientlists.GetClientListItemsRequest{
					ListID: "193089_VPTERRAFORM2",
				}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString(getPath(testDir, "user_list_user_id.tf")),
					Check: baseChecker.
						CheckEqual("list.list_id", "193089_VPTERRAFORM2").
						CheckEqual("list.items.#", "2").
						CheckEqual("list.items.0.value", "8728690e-cd6c-42d1-94e8-6a9a9f326cb5").
						CheckEqual("list.items.1.value", "e00b5827-7105-4366-bc24-aa735fc18e4c").
						CheckEqual("output_text", loadText(t, getPath(testDir, "user_output_text_user_id.txt"))).
						CheckEqual("json", loadJSON(t, getPath(testDir, "user_output_json_user_id.txt"))).
						Build(),
				},
			},
		},
		"error response from GetClientList api": {
			listType: clientlists.IP,
			init: func(m *clientlists.Mock) {
				mockGetClientListFailure(m, clientlists.GetClientListRequest{
					ListID:       "180991_TESTNLMIGRATION123",
					IncludeItems: true,
				})
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(getPath(testDir, "ip_list.tf")),
					ExpectError: regexp.MustCompile("Error: get client list error"),
				},
			},
		},
		"error response from GetClientListItems api": {
			listType: clientlists.USER,
			init: func(m *clientlists.Mock) {
				err := json.Unmarshal(testutils.LoadFixtureBytes(t, getPath(testDir, "user_client_list.json")), &getClientListResponseUsername)
				require.NoError(t, err)
				mockGetClientList(m, *getClientListResponseUsername, clientlists.GetClientListRequest{
					ListID:       "193203_VPUSERTYPECLIENTLISTUS",
					IncludeItems: true,
				}, 1)
				mockGetClientListItemsFailure(m, clientlists.GetClientListItemsRequest{
					ListID: "193203_VPUSERTYPECLIENTLISTUS",
				})
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString(getPath(testDir, "user_list.tf")),
					ExpectError: regexp.MustCompile("Error: get client list items error"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &clientlists.Mock{}
			if test.init != nil {
				test.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})

			client.AssertExpectations(t)
		})
	}
}

func mockGetClientList(m *clientlists.Mock, response clientlists.GetClientListResponse, request clientlists.GetClientListRequest, times int) {
	m.On("GetClientList", mock.Anything, request).
		Return(&response, nil).Times(times)
}

func mockGetClientListFailure(m *clientlists.Mock, request clientlists.GetClientListRequest) {
	err := errors.New("get client list error")
	m.On("GetClientList", mock.Anything, request).
		Return(nil, err).Once()
}

func mockGetClientListItems(m *clientlists.Mock, response clientlists.GetClientListItemsResponse, request clientlists.GetClientListItemsRequest, times int) {
	m.On("GetClientListItems", mock.Anything, request).
		Return(&response, nil).Times(times)
}

func mockGetClientListItemsFailure(m *clientlists.Mock, request clientlists.GetClientListItemsRequest) {
	err := errors.New("get client list items error")
	m.On("GetClientListItems", mock.Anything, request).
		Return(nil, err).Once()
}

func loadText(t *testing.T, path string) string {
	return testutils.LoadFixtureString(t, path)
}

func loadJSON(t *testing.T, path string) string {
	return strings.TrimSuffix(testutils.LoadFixtureString(t, path), "\n")
}

func getPath(testDir string, fileName string) string {
	return fmt.Sprintf("%s/%s", testDir, fileName)
}
