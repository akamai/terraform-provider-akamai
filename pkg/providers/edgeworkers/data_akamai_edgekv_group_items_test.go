package edgeworkers

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/edgeworkers"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestEdgeKVGroupItems(t *testing.T) {
	client := &edgeworkers.Mock{}
	client.Test(testutils.TattleT{T: t})

	items := map[string]string{
		"TestItem1": "TestValue1",
		"TestItem2": "TestValue2",
		"TestItem3": "TestValue3",
	}

	t.Run("happy path", func(t *testing.T) {
		client.On("ListItems", testutils.MockContext, edgeworkers.ListItemsRequest{
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				Network:     "staging",
				NamespaceID: "test_namespace",
				GroupID:     "TestGroup",
			},
		}).Return(&edgeworkers.ListItemsResponse{"TestItem1", "TestItem2", "TestItem3"}, nil).Times(3)

		for k, v := range items {
			mockGetItemReq(client, k, edgeworkers.Item(v))
		}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataEdgeKVGroupItems/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttrSet("data.akamai_edgekv_group_items.test", "id"),
							resource.TestCheckResourceAttr("data.akamai_edgekv_group_items.test", "id", "test_namespace:staging:TestGroup"),
							resource.TestCheckResourceAttr("data.akamai_edgekv_group_items.test", "items.%", "3"),
							resource.TestCheckResourceAttr("data.akamai_edgekv_group_items.test", "items.TestItem1", "TestValue1"),
							resource.TestCheckResourceAttr("data.akamai_edgekv_group_items.test", "items.TestItem2", "TestValue2"),
							resource.TestCheckResourceAttr("data.akamai_edgekv_group_items.test", "items.TestItem3", "TestValue3"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("missed required `namespace_name` field", func(t *testing.T) {
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataEdgeKVGroupItems/missed_namespace_name.tf"),
						ExpectError: regexp.MustCompile(`The argument "namespace_name" is required, but no definition was found.`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("missed required `network` field", func(t *testing.T) {
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataEdgeKVGroupItems/missed_network.tf"),
						ExpectError: regexp.MustCompile(`The argument "network" is required, but no definition was found.`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("missed required `group_name` field", func(t *testing.T) {
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataEdgeKVGroupItems/missed_group.tf"),
						ExpectError: regexp.MustCompile(`The argument "group_name" is required, but no definition was found.`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("incorrect `network` field", func(t *testing.T) {
		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestDataEdgeKVGroupItems/incorrect_network.tf"),
						ExpectError: regexp.MustCompile(`expected network to be one of \["staging" "production"], got incorrect_network`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func mockGetItemReq(client *edgeworkers.Mock, itemID string, itemValue edgeworkers.Item) *mock.Call {

	return client.On("GetItem", testutils.MockContext, edgeworkers.GetItemRequest{
		ItemID: itemID,
		ItemsRequestParams: edgeworkers.ItemsRequestParams{
			Network:     "staging",
			NamespaceID: "test_namespace",
			GroupID:     "TestGroup",
		},
	}).Return(&itemValue, nil).Times(3)
}
