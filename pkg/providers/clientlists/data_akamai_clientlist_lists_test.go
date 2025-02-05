package clientlists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestClientList_data_all_lists(t *testing.T) {
	dataSourceName := "data.akamai_clientlist_lists.lists"
	tests := map[string]struct {
		params       clientlists.GetClientListsRequest
		config       string
		responseBody []byte
	}{
		"All lists": {
			params: clientlists.GetClientListsRequest{
				Type: []clientlists.ClientListType{},
			},
			config:       loadFixtureString("testData/TestDSClientList/match_all.tf"),
			responseBody: loadFixtureBytes("testData/TestDSClientList/ClientLists.json"),
		},
		"Filtered lists": {
			params: clientlists.GetClientListsRequest{
				Name: "test",
				Type: []clientlists.ClientListType{clientlists.IP, clientlists.GEO},
			},
			config:       loadFixtureString("testData/TestDSClientList/match_by_filters.tf"),
			responseBody: loadFixtureBytes("testData/TestDSClientList/ClientLists.json"),
		},
		"Empty content list": {
			params: clientlists.GetClientListsRequest{
				Type: []clientlists.ClientListType{},
			},
			config: loadFixtureString("testData/TestDSClientList/match_all.tf"),
			responseBody: []byte(`{
				"content": []
			}`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &clientlists.Mock{}
			clientListsResponse := clientlists.GetClientListsResponse{}
			err := json.Unmarshal(test.responseBody, &clientListsResponse)
			require.NoError(t, err)

			r, err := json.Marshal(clientListsResponse.Content)
			require.NoError(t, err)

			buf := &bytes.Buffer{}
			err = json.Indent(buf, r, "", "  ")
			require.NoError(t, err)

			clientListsResponseJSONString := buf.String()

			client.On("GetClientLists",
				testutils.MockContext,
				test.params,
			).Return(&clientListsResponse, nil)

			useClient(client, func() {
				checks := []resource.TestCheckFunc{
					resource.TestCheckResourceAttr(dataSourceName, "json", clientListsResponseJSONString),
					resource.TestCheckResourceAttr(dataSourceName, "list_ids.#", strconv.Itoa(len(clientListsResponse.Content))),
					resource.TestCheckResourceAttr(dataSourceName, "lists.#", strconv.Itoa(len(clientListsResponse.Content))),
				}
				for k, v := range clientListsResponse.Content {
					checks = append(checks, resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("list_ids.%d", k), v.ListID))
					checks = append(checks, resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("lists.%d.list_id", k), v.ListID))
				}

				if test.params.Name != "" {
					checks = append(checks, resource.TestCheckResourceAttr(dataSourceName, "name", test.params.Name))
				} else {
					checks = append(checks, resource.TestCheckNoResourceAttr(dataSourceName, "name"))
				}
				if len(test.params.Type) > 0 {
					checks = append(checks, resource.TestCheckResourceAttr(dataSourceName, "type.#", strconv.Itoa(len(test.params.Type))))
					checks = append(checks, resource.TestCheckResourceAttr(dataSourceName, "type.0", string(test.params.Type[1])))
					checks = append(checks, resource.TestCheckResourceAttr(dataSourceName, "type.1", string(test.params.Type[0])))
				} else {
					checks = append(checks, resource.TestCheckNoResourceAttr(dataSourceName, "type"))
				}

				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							Config: test.config,
							Check:  resource.ComposeAggregateTestCheckFunc(checks...),
						},
					},
				})
			})

			client.AssertExpectations(t)
		})
	}
}
