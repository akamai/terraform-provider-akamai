package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataAkamaiBotCategory(t *testing.T) {
	t.Run("DataAkamaiBotCategory", func(t *testing.T) {
		mockedBotmanClient := &botman.Mock{}

		response := botman.GetAkamaiBotCategoryListResponse{
			Categories: []map[string]interface{}{
				{"categoryId": "b85e3eaa-d334-466d-857e-33308ce416be", "categoryName": "Test name 1", "testKey": "testValue1"},
				{"categoryId": "69acad64-7459-4c1d-9bad-672600150127", "categoryName": "Test name 2", "testKey": "testValue2"},
				{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "categoryName": "Test name 3", "testKey": "testValue3"},
				{"categoryId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "categoryName": "Test name 4", "testKey": "testValue4"},
				{"categoryId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "categoryName": "Test name 5", "testKey": "testValue5"},
			},
		}
		expectedJSON := `
{
	"categories":[
		{"categoryId":"b85e3eaa-d334-466d-857e-33308ce416be", "categoryName": "Test name 1", "testKey":"testValue1"},
		{"categoryId":"69acad64-7459-4c1d-9bad-672600150127", "categoryName": "Test name 2", "testKey":"testValue2"},
		{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "categoryName": "Test name 3", "testKey":"testValue3"},
		{"categoryId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "categoryName": "Test name 4", "testKey":"testValue4"},
		{"categoryId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "categoryName": "Test name 5", "testKey":"testValue5"}
	]
}`
		mockedBotmanClient.On("GetAkamaiBotCategoryList",
			mock.Anything,
			botman.GetAkamaiBotCategoryListRequest{},
		).Return(&response, nil)
		useClient(mockedBotmanClient, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataAkamaiBotCategory/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_akamai_bot_category.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
	t.Run("DataAkamaiBotCategory filter by name", func(t *testing.T) {
		mockedBotmanClient := &botman.Mock{}

		response := botman.GetAkamaiBotCategoryListResponse{
			Categories: []map[string]interface{}{
				{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "categoryName": "Test name 3", "testKey": "testValue3"},
			},
		}
		expectedJSON := `
{
	"categories":[
		{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "categoryName": "Test name 3", "testKey":"testValue3"}
	]
}`
		mockedBotmanClient.On("GetAkamaiBotCategoryList",
			mock.Anything,
			botman.GetAkamaiBotCategoryListRequest{CategoryName: "Test name 3"},
		).Return(&response, nil)
		useClient(mockedBotmanClient, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataAkamaiBotCategory/filter_by_name.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_akamai_bot_category.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
