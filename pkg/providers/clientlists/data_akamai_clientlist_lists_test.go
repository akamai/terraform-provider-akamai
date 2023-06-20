package clientlists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/clientlists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClientList_data_all_lists(t *testing.T) {
	t.Run("match all lists", func(t *testing.T) {
		const dataSourceName = "data.akamai_clientlist_lists.lists"
		client := &clientlists.Mock{}

		clientListsResponse := clientlists.GetClientListsResponse{}
		err := json.Unmarshal(loadFixtureBytes("testData/TestDSClientList/ClientLists.json"), &clientListsResponse)
		require.NoError(t, err)

		r, err := json.Marshal(clientListsResponse.Content)
		require.NoError(t, err)

		buf := &bytes.Buffer{}
		err = json.Indent(buf, r, "", "  ")
		require.NoError(t, err)

		clientListsResponseJSONString := buf.String()

		client.On("GetClientLists",
			mock.Anything,
			clientlists.GetClientListsRequest{},
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

			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testData/TestDSClientList/match_all.tf"),
						Check:  resource.ComposeAggregateTestCheckFunc(checks...),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
