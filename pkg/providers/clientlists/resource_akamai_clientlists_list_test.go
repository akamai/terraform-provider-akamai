package clientlists

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/clientlists"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceClientList(t *testing.T) {
	type listAttributes struct {
		ListID, Name, Notes, Type, ContractID string
		Tags                                  []string
		GroupID                               int
		Version                               int
		ItemsCount                            int
		Items                                 []clientlists.ListItemPayload
	}

	const testDir = "testData/TestResClientList"

	var (
		updateAPIError = "{\n      \"type\": \"https://problems.luna.akamaiapis.net/client-list-api/error-types/INVALID-INPUT-ERROR\",\n       \"status\": 400,\n       \"title\": \"Invalid Input Error\",\n       \"detail\": \"Validation failed: tags: tags exceeds max size of 5 entries\",\n       \"instance\": \"https://problems.luna.akamaiapis.net/client-list-api/error-instances/9ff3649993cb002b\",\n       \"\n }\n"

		mapItemsPayloadToContent = func(items []clientlists.ListItemPayload) []clientlists.ListItemContent {
			result := make([]clientlists.ListItemContent, 0, len(items))
			for _, v := range items {
				result = append(result, clientlists.ListItemContent{
					Value:          v.Value,
					Description:    v.Description,
					Tags:           v.Tags,
					ExpirationDate: v.ExpirationDate,
				})
			}
			return result
		}

		expectCreateList = func(t *testing.T, client *clientlists.Mock, req clientlists.CreateClientListRequest) *clientlists.CreateClientListResponse {

			createResponse := clientlists.CreateClientListResponse{
				ListContent: clientlists.ListContent{
					ListID:     "1_AB",
					Name:       req.Name,
					Notes:      req.Notes,
					Tags:       req.Tags,
					Type:       req.Type,
					Version:    1,
					ItemsCount: int64(len(req.Items)),
				},
				ContractID: req.ContractID,
				GroupID:    req.GroupID,
				GroupName:  "Group-Name",
				Items:      mapItemsPayloadToContent(req.Items),
			}

			client.On("CreateClientList", mock.Anything, req).Return(&createResponse, nil).Once()
			return &createResponse
		}

		expectUpdateList = func(t *testing.T, client *clientlists.Mock, listType clientlists.ClientListType, itemsCount int64, req clientlists.UpdateClientListRequest) *clientlists.UpdateClientListResponse {

			updateResponse := clientlists.UpdateClientListResponse{
				ListContent: clientlists.ListContent{
					ListID:     "1_AB",
					Name:       req.Name,
					Notes:      req.Notes,
					Tags:       req.Tags,
					Type:       listType,
					Version:    1,
					ItemsCount: itemsCount,
				},
			}

			client.On("UpdateClientList", mock.Anything, req).Return(&updateResponse, nil).Once()
			return &updateResponse
		}

		expectUpdateListItems = func(t *testing.T, client *clientlists.Mock, req clientlists.UpdateClientListItemsRequest) *clientlists.UpdateClientListItemsResponse {
			appended := make([]clientlists.ListItemContent, 0, len(req.Append))
			for _, v := range req.Append {
				appended = append(appended, clientlists.ListItemContent{
					Value:          v.Value,
					Description:    v.Description,
					Tags:           v.Tags,
					ExpirationDate: v.ExpirationDate,
				})
			}
			updated := make([]clientlists.ListItemContent, 0, len(req.Update))
			for _, v := range req.Update {
				updated = append(updated, clientlists.ListItemContent{
					Value:          v.Value,
					Description:    v.Description,
					Tags:           v.Tags,
					ExpirationDate: v.ExpirationDate,
				})
			}
			deleted := make([]clientlists.ListItemContent, 0, len(req.Delete))
			for _, v := range req.Delete {
				deleted = append(deleted, clientlists.ListItemContent{
					Value:          v.Value,
					Description:    v.Description,
					Tags:           v.Tags,
					ExpirationDate: v.ExpirationDate,
				})
			}

			updateResponse := clientlists.UpdateClientListItemsResponse{
				Appended: appended,
				Updated:  updated,
				Deleted:  deleted,
			}

			client.On("UpdateClientListItems", mock.Anything, req).Return(&updateResponse, nil).Once()
			return &updateResponse
		}

		expectReadList = func(t *testing.T, client *clientlists.Mock, list clientlists.ListContent, items []clientlists.ListItemContent, callTimes int) {
			clientListGetReq := clientlists.GetClientListRequest{
				ListID:       list.ListID,
				IncludeItems: true,
			}

			clientList := clientlists.GetClientListResponse{
				ListContent: list,
				Items:       items,
				ContractID:  "12_ABC",
				GroupID:     12,
			}
			client.On("GetClientList", mock.Anything, clientListGetReq).Return(&clientList, nil).Times(callTimes)
		}

		expectDeleteList = func(t *testing.T, client *clientlists.Mock, list clientlists.ListContent) {
			clientListDeleteReq := clientlists.DeleteClientListRequest{
				ListID: list.ListID,
			}
			client.On("DeleteClientList", mock.Anything, clientListDeleteReq).Return(nil).Once()
		}

		expectAPIErrorWithUpdateList = func(t *testing.T, client *clientlists.Mock, req clientlists.UpdateClientListRequest) {
			err := fmt.Errorf(updateAPIError)
			client.On("UpdateClientList", mock.Anything, req).Return(nil, err).Once()
		}

		checkAttributes = func(attrs listAttributes) resource.TestCheckFunc {
			resourceName := "akamai_clientlist_list.test_list"
			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "list_id", attrs.ListID),
				resource.TestCheckResourceAttr(resourceName, "name", attrs.Name),
				resource.TestCheckResourceAttr(resourceName, "notes", attrs.Notes),
				resource.TestCheckResourceAttr(resourceName, "type", attrs.Type),
				resource.TestCheckResourceAttr(resourceName, "tags.#", strconv.Itoa(len(attrs.Tags))),
				resource.TestCheckResourceAttr(resourceName, "contract_id", attrs.ContractID),
				resource.TestCheckResourceAttr(resourceName, "group_id", strconv.Itoa(attrs.GroupID)),
				resource.TestCheckResourceAttr(resourceName, "version", strconv.Itoa(attrs.Version)),
				resource.TestCheckResourceAttr(resourceName, "items_count", strconv.Itoa(attrs.ItemsCount)),
				resource.TestCheckResourceAttr(resourceName, "items.#", strconv.Itoa(len(attrs.Items))),
			}

			return resource.ComposeAggregateTestCheckFunc(checks...)
		}
	)

	t.Run("Create a new client list", func(t *testing.T) {
		client := new(clientlists.Mock)
		clientList := expectCreateList(t, client, clientlists.CreateClientListRequest{
			Name:       "List Name",
			Notes:      "List Notes",
			Tags:       []string{"a", "b"},
			Type:       clientlists.ASN,
			ContractID: "12_ABC",
			GroupID:    12,
			Items:      []clientlists.ListItemPayload{},
		})
		expectReadList(t, client, clientList.ListContent, []clientlists.ListItemContent{}, 2)
		expectDeleteList(t, client, clientList.ListContent)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_create.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name",
							Notes:      "List Notes",
							Tags:       []string{"a", "b"},
							Type:       "ASN",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: 0,
							Items:      []clientlists.ListItemPayload{},
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Update client list", func(t *testing.T) {
		client := new(clientlists.Mock)
		clientList := expectCreateList(t, client, clientlists.CreateClientListRequest{
			Name:       "List Name",
			Notes:      "List Notes",
			Tags:       []string{"a", "b"},
			Type:       clientlists.ASN,
			ContractID: "12_ABC",
			GroupID:    12,
			Items:      []clientlists.ListItemPayload{},
		})
		expectReadList(t, client, clientList.ListContent, []clientlists.ListItemContent{}, 3)
		updateResponse := expectUpdateList(t, client, clientlists.ASN, 0, clientlists.UpdateClientListRequest{
			UpdateClientList: clientlists.UpdateClientList{
				Name:  "List Name Updated",
				Notes: "List Notes Updated",
				Tags:  []string{"a", "c", "d"},
			},
			ListID: clientList.ListID,
		})
		expectReadList(t, client, updateResponse.ListContent, []clientlists.ListItemContent{}, 2)
		expectDeleteList(t, client, clientList.ListContent)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_create.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name",
							Notes:      "List Notes",
							Tags:       []string{"a", "b"},
							Type:       "ASN",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: 0,
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_update.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name Updated",
							Notes:      "List Notes Updated",
							Tags:       []string{"a", "c", "d"},
							Type:       "ASN",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: 0,
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Update client list not expected when empty tags list removed", func(t *testing.T) {
		client := new(clientlists.Mock)
		clientList := expectCreateList(t, client, clientlists.CreateClientListRequest{
			Name:       "List Name",
			Notes:      "List Notes",
			Tags:       []string{},
			Type:       clientlists.IP,
			ContractID: "12_ABC",
			GroupID:    12,
			Items:      []clientlists.ListItemPayload{},
		})
		expectReadList(t, client, clientList.ListContent, []clientlists.ListItemContent{}, 4)
		expectDeleteList(t, client, clientList.ListContent)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_create_empty_tags.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name",
							Notes:      "List Notes",
							Type:       "IP",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: 0,
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_update_remove_tags.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name",
							Notes:      "List Notes",
							Tags:       []string{},
							Type:       "IP",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: 0,
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Update client list returns an API error", func(t *testing.T) {
		client := new(clientlists.Mock)
		clientList := expectCreateList(t, client, clientlists.CreateClientListRequest{
			Name:       "List Name",
			Notes:      "List Notes",
			Tags:       []string{"a", "b"},
			Type:       clientlists.ASN,
			ContractID: "12_ABC",
			GroupID:    12,
			Items:      []clientlists.ListItemPayload{},
		})
		expectReadList(t, client, clientList.ListContent, []clientlists.ListItemContent{}, 3)

		expectAPIErrorWithUpdateList(t, client, clientlists.UpdateClientListRequest{
			UpdateClientList: clientlists.UpdateClientList{
				Name:  "List Name Updated",
				Notes: "List Notes Updated",
				Tags:  []string{"a", "c", "d"},
			},
			ListID: clientList.ListID,
		})
		expectDeleteList(t, client, clientList.ListContent)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_create.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name",
							Notes:      "List Notes",
							Tags:       []string{"a", "b"},
							Type:       "ASN",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: 0,
						}),
					},
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/list_update.tf", testDir)),
						ExpectError: regexp.MustCompile(updateAPIError),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Create a new client list with items", func(t *testing.T) {
		client := new(clientlists.Mock)
		items := append([]clientlists.ListItemPayload{},
			clientlists.ListItemPayload{
				Value:       "1",
				Description: "Item 1 Desc",
				Tags:        []string{"item1Tag2", "item1Tag1"},
			},
			clientlists.ListItemPayload{
				Value:          "123",
				ExpirationDate: "2026-12-26T01:00:00+00:00",
				Tags:           []string{},
			},
			clientlists.ListItemPayload{
				Value:       "12",
				Description: "Item 12 Desc",
				Tags:        []string{"item12Tag1", "item12Tag2"},
			})

		clientList := expectCreateList(t, client, clientlists.CreateClientListRequest{
			Name:       "List Name",
			Notes:      "List Notes",
			Tags:       []string{"a", "b"},
			Type:       clientlists.ASN,
			ContractID: "12_ABC",
			GroupID:    12,
			Items:      items,
		})
		expectReadList(t, client, clientList.ListContent, mapItemsPayloadToContent(items), 2)
		expectDeleteList(t, client, clientList.ListContent)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_and_items_create.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name",
							Notes:      "List Notes",
							Tags:       []string{"a", "b"},
							Type:       "ASN",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: len(items),
							Items:      items,
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Update client list items and list", func(t *testing.T) {
		client := new(clientlists.Mock)
		items := append([]clientlists.ListItemPayload{},
			clientlists.ListItemPayload{
				Value:       "1",
				Description: "Item 1 Desc",
				Tags:        []string{"item1Tag2", "item1Tag1"},
			},
			clientlists.ListItemPayload{
				Value:          "123",
				ExpirationDate: "2026-12-26T01:00:00+00:00",
				Tags:           []string{},
			},
			clientlists.ListItemPayload{
				Value:       "12",
				Description: "Item 12 Desc",
				Tags:        []string{"item12Tag1", "item12Tag2"},
			})
		updatedItems := append([]clientlists.ListItemPayload{},
			clientlists.ListItemPayload{
				Value:       "1",
				Description: "Item 1 Desc",
				Tags:        []string{"item1Tag2", "item1Tag1"},
			},
			clientlists.ListItemPayload{
				Value:       "12",
				Description: "Item 12 Desc Updated",
				Tags:        []string{"item12Tag1", "item12Tag2"},
			},
			clientlists.ListItemPayload{
				Value:       "1234",
				Description: "Item 1234 Desc",
				Tags:        []string{"1234Tag"},
			})

		clientList := expectCreateList(t, client, clientlists.CreateClientListRequest{
			Name:       "List Name",
			Notes:      "List Notes",
			Tags:       []string{"a", "b"},
			Type:       clientlists.ASN,
			ContractID: "12_ABC",
			GroupID:    12,
			Items:      items,
		})
		expectReadList(t, client, clientList.ListContent, mapItemsPayloadToContent(items), 4)
		updateResponse := expectUpdateList(t, client, clientlists.ASN, 3, clientlists.UpdateClientListRequest{
			UpdateClientList: clientlists.UpdateClientList{
				Name:  "List Name Updated",
				Notes: "List Notes Updated",
				Tags:  []string{"a", "c", "d"},
			},
			ListID: clientList.ListID,
		})
		expectUpdateListItems(t, client, clientlists.UpdateClientListItemsRequest{
			ListID: clientList.ListID,
			UpdateClientListItems: clientlists.UpdateClientListItems{
				Append: []clientlists.ListItemPayload{
					{
						Value:       "1234",
						Description: "Item 1234 Desc",
						Tags:        []string{"1234Tag"},
					},
				},
				Update: []clientlists.ListItemPayload{
					{
						Value:       "12",
						Description: "Item 12 Desc Updated",
						Tags:        []string{"item12Tag1", "item12Tag2"},
					},
				},
				Delete: []clientlists.ListItemPayload{
					{
						Value: "123",
					},
				},
			},
		})
		expectReadList(t, client, updateResponse.ListContent, mapItemsPayloadToContent(updatedItems), 2)
		expectDeleteList(t, client, clientList.ListContent)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_and_items_create.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name",
							Notes:      "List Notes",
							Tags:       []string{"a", "b"},
							Type:       "ASN",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: 3,
							Items:      items,
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_and_items_update.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name Updated",
							Notes:      "List Notes Updated",
							Tags:       []string{"a", "c", "d"},
							Type:       "ASN",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: 3,
							Items:      updatedItems,
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Update client list items only", func(t *testing.T) {
		client := new(clientlists.Mock)
		items := append([]clientlists.ListItemPayload{},
			clientlists.ListItemPayload{
				Value:       "1",
				Description: "Item 1 Desc",
				Tags:        []string{"item1Tag2", "item1Tag1"},
			},
			clientlists.ListItemPayload{
				Value:          "123",
				ExpirationDate: "2026-12-26T01:00:00+00:00",
				Tags:           []string{},
			},
			clientlists.ListItemPayload{
				Value:       "12",
				Description: "Item 12 Desc",
				Tags:        []string{"item12Tag1", "item12Tag2"},
			})
		updatedItems := append([]clientlists.ListItemPayload{},
			clientlists.ListItemPayload{
				Value:       "1",
				Description: "Item 1 Desc",
				Tags:        []string{"item1Tag2", "item1Tag1"},
			},
			clientlists.ListItemPayload{
				Value:       "12",
				Description: "Item 12 Desc Updated",
				Tags:        []string{"item12Tag1", "item12Tag2"},
			},
			clientlists.ListItemPayload{
				Value:       "1234",
				Description: "Item 1234 Desc",
				Tags:        []string{"1234Tag"},
			})

		clientList := expectCreateList(t, client, clientlists.CreateClientListRequest{
			Name:       "List Name",
			Notes:      "List Notes",
			Tags:       []string{"a", "b"},
			Type:       clientlists.ASN,
			ContractID: "12_ABC",
			GroupID:    12,
			Items:      items,
		})
		expectReadList(t, client, clientList.ListContent, mapItemsPayloadToContent(items), 4)
		expectUpdateListItems(t, client, clientlists.UpdateClientListItemsRequest{
			ListID: clientList.ListID,
			UpdateClientListItems: clientlists.UpdateClientListItems{
				Append: []clientlists.ListItemPayload{
					{
						Value:       "1234",
						Description: "Item 1234 Desc",
						Tags:        []string{"1234Tag"},
					},
				},
				Update: []clientlists.ListItemPayload{
					{
						Value:       "12",
						Description: "Item 12 Desc Updated",
						Tags:        []string{"item12Tag1", "item12Tag2"},
					},
				},
				Delete: []clientlists.ListItemPayload{
					{
						Value: "123",
					},
				},
			},
		})
		expectReadList(t, client, clientList.ListContent, mapItemsPayloadToContent(updatedItems), 2)
		expectDeleteList(t, client, clientList.ListContent)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_and_items_create.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name",
							Notes:      "List Notes",
							Tags:       []string{"a", "b"},
							Type:       "ASN",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: 3,
							Items:      items,
						}),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_items_only_update.tf", testDir)),
						Check: checkAttributes(listAttributes{
							ListID:     clientList.ListID,
							Name:       "List Name",
							Notes:      "List Notes",
							Tags:       []string{"a", "b"},
							Type:       "ASN",
							ContractID: "12_ABC",
							GroupID:    12,
							Version:    1,
							ItemsCount: 3,
							Items:      updatedItems,
						}),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Update items set new computed version", func(t *testing.T) {
		client := new(clientlists.Mock)
		items := append([]clientlists.ListItemPayload{},
			clientlists.ListItemPayload{
				Value:       "1",
				Description: "Item 1 Desc",
				Tags:        []string{"item1Tag2", "item1Tag1"},
			})
		updatedItems := []clientlists.ListItemPayload{}

		clientList := expectCreateList(t, client, clientlists.CreateClientListRequest{
			Name:       "List Name",
			Notes:      "List Notes",
			Tags:       []string{"a", "b"},
			Type:       clientlists.ASN,
			ContractID: "12_ABC",
			GroupID:    12,
			Items:      items,
		})
		expectReadList(t, client, clientList.ListContent, mapItemsPayloadToContent(items), 4)
		expectUpdateListItems(t, client, clientlists.UpdateClientListItemsRequest{
			ListID: clientList.ListID,
			UpdateClientListItems: clientlists.UpdateClientListItems{
				Append: []clientlists.ListItemPayload{},
				Update: []clientlists.ListItemPayload{},
				Delete: []clientlists.ListItemPayload{{Value: "1"}},
			},
		})
		// Update version
		updatedClientList := clientList.ListContent
		updatedClientList.Version = 2

		expectReadList(t, client, updatedClientList, mapItemsPayloadToContent(updatedItems), 2)
		expectDeleteList(t, client, clientList.ListContent)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_and_items_create_one_item.tf", testDir)),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("version", "1"),
						),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_items_update_compute_version.tf", testDir)),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("version", "2"),
						),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Update items NOT set new computed version", func(t *testing.T) {
		client := new(clientlists.Mock)
		items := append([]clientlists.ListItemPayload{},
			clientlists.ListItemPayload{
				Value:       "1",
				Description: "Item 1 Desc",
				Tags:        []string{"item1Tag2", "item1Tag1"},
			})
		updatedItems := append([]clientlists.ListItemPayload{},
			clientlists.ListItemPayload{
				Value: "1",
			})

		clientList := expectCreateList(t, client, clientlists.CreateClientListRequest{
			Name:       "List Name",
			Notes:      "List Notes",
			Tags:       []string{"a", "b"},
			Type:       clientlists.ASN,
			ContractID: "12_ABC",
			GroupID:    12,
			Items:      items,
		})
		expectReadList(t, client, clientList.ListContent, mapItemsPayloadToContent(items), 4)
		expectUpdateListItems(t, client, clientlists.UpdateClientListItemsRequest{
			ListID: clientList.ListID,
			UpdateClientListItems: clientlists.UpdateClientListItems{
				Append: []clientlists.ListItemPayload{},
				Update: []clientlists.ListItemPayload{
					{
						Value: "1",
						Tags:  []string{},
					},
				},
				Delete: []clientlists.ListItemPayload{},
			},
		})
		// Fake version update
		updatedClientList := clientList.ListContent
		updatedClientList.Version = 2
		expectReadList(t, client, updatedClientList, mapItemsPayloadToContent(updatedItems), 2)
		expectDeleteList(t, client, clientList.ListContent)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_and_items_create_one_item.tf", testDir)),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("version", "1"),
						),
					},
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_items_update_not_compute_version.tf", testDir)),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("version", "1"),
						),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Create list with duplicate items fails", func(t *testing.T) {
		client := new(clientlists.Mock)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString(fmt.Sprintf("%s/list_and_duplicate_items_create.tf", testDir)),
						ExpectError: regexp.MustCompile("Error: 'Items' collection contains duplicate values for 'value' field. Duplicate value: 12"),
					},
				},
			})
		})
		client.AssertExpectations(t)
	})

	t.Run("Import clientlist resource", func(t *testing.T) {
		client := new(clientlists.Mock)

		clientList := expectCreateList(t, client, clientlists.CreateClientListRequest{
			Name:       "List Name",
			Notes:      "List Notes",
			Tags:       []string{"a", "b"},
			Type:       clientlists.ASN,
			ContractID: "12_ABC",
			GroupID:    12,
			Items:      []clientlists.ListItemPayload{},
		})
		expectReadList(t, client, clientList.ListContent, []clientlists.ListItemContent{}, 3)
		expectDeleteList(t, client, clientList.ListContent)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString(fmt.Sprintf("%s/list_create.tf", testDir)),
					},
					{
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateId:     "1_AB",
						ResourceName:      "akamai_clientlist_list.test_list",
					},
				},
			})
		})
		client.AssertExpectations(t)
	})
}
