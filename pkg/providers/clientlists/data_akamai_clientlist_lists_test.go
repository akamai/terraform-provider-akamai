package clientlists

import (
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDataClientLists(t *testing.T) {
	allListsResponse := clientlists.GetClientListsResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testData/TestDSClientList/ClientLists.json"), &allListsResponse)
	require.NoError(t, err)

	emptyListsResponse := clientlists.GetClientListsResponse{
		Content: []clientlists.ClientList{},
	}

	baseChecker := test.NewStateChecker("data.akamai_clientlist_lists.lists")

	tests := map[string]struct {
		init  func(*clientlists.Mock)
		steps []resource.TestStep
	}{
		"happy path - all lists": {
			init: func(m *clientlists.Mock) {
				mockGetClientLists(m, allListsResponse, clientlists.GetClientListsRequest{}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testData/TestDSClientList/match_all.tf"),
					Check: baseChecker.
						CheckEqual("list_ids.#", "10").
						CheckEqual("lists.#", "10").
						CheckEqual("lists.0.list_id", "91596_AUDITLOGSTESTLIST").
						CheckEqual("lists.0.name", "AUDIT LOGS - TEST LIST").
						CheckEqual("lists.0.type", "IP").
						Build(),
				},
			},
		},
		"happy path - filtered lists": {
			init: func(m *clientlists.Mock) {
				mockGetClientLists(m, allListsResponse, clientlists.GetClientListsRequest{
					Name: "test",
					Type: []clientlists.ClientListType{clientlists.GEO, clientlists.IP},
				}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testData/TestDSClientList/match_by_filters.tf"),
					Check: baseChecker.
						CheckEqual("name", "test").
						CheckEqual("type.#", "2").
						CheckEqual("type.0", "GEO").
						CheckEqual("type.1", "IP").
						CheckEqual("list_ids.#", "10").
						CheckEqual("lists.#", "10").
						Build(),
				},
			},
		},
		"happy path - user type lists": {
			init: func(m *clientlists.Mock) {
				mockGetClientLists(m, allListsResponse, clientlists.GetClientListsRequest{
					Name: "test",
					Type: []clientlists.ClientListType{clientlists.USER},
				}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testData/TestDSClientList/user_type.tf"),
					Check: baseChecker.
						CheckEqual("name", "test").
						CheckEqual("type.#", "1").
						CheckEqual("type.0", string(clientlists.USER)).
						CheckEqual("list_ids.#", "10").
						CheckEqual("lists.#", "10").
						Build(),
				},
			},
		},
		"happy path - empty content list": {
			init: func(m *clientlists.Mock) {
				mockGetClientLists(m, emptyListsResponse, clientlists.GetClientListsRequest{}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testData/TestDSClientList/match_all.tf"),
					Check: baseChecker.
						CheckEqual("list_ids.#", "0").
						CheckEqual("lists.#", "0").
						CheckEqual("json", "[]").
						Build(),
				},
			},
		},
		"error response from GetClientLists api": {
			init: func(m *clientlists.Mock) {
				mockGetClientListsFailure(m, clientlists.GetClientListsRequest{})
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testData/TestDSClientList/match_all.tf"),
					ExpectError: regexp.MustCompile("get client lists error"),
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

func mockGetClientLists(m *clientlists.Mock, response clientlists.GetClientListsResponse, request clientlists.GetClientListsRequest, times int) {
	m.On("GetClientLists", mock.Anything, request).
		Return(&response, nil).Times(times)
}

func mockGetClientListsFailure(m *clientlists.Mock, request clientlists.GetClientListsRequest) {
	err := errors.New("failed to get client lists")
	m.On("GetClientLists", mock.Anything, request).
		Return(nil, err).Once()
}
