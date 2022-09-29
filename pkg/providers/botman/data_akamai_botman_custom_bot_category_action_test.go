package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataCustomBotCategoryAction(t *testing.T) {
	t.Run("DataCustomBotCategoryAction", func(t *testing.T) {

		mockedBotmanClient := &mockbotman{}
		response := botman.GetCustomBotCategoryActionListResponse{
			Actions: []map[string]interface{}{
				{"categoryId": "b85e3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
				{"categoryId": "69acad64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
				{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				{"categoryId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
				{"categoryId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
			},
		}
		expectedJSON := `
{
	"actions":[
		{"categoryId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"categoryId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"categoryId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"categoryId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`
		mockedBotmanClient.On("GetCustomBotCategoryActionList",
			mock.Anything,
			botman.GetCustomBotCategoryActionListRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230"},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataCustomBotCategoryAction/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_custom_bot_category_action.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
	t.Run("DataCustomBotCategoryAction filter by id", func(t *testing.T) {

		mockedBotmanClient := &mockbotman{}
		response := botman.GetCustomBotCategoryActionListResponse{
			Actions: []map[string]interface{}{
				{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			},
		}
		expectedJSON := `
{
	"actions":[
		{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}
	]
}`
		mockedBotmanClient.On("GetCustomBotCategoryActionList",
			mock.Anything,
			botman.GetCustomBotCategoryActionListRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230", CategoryID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataCustomBotCategoryAction/filter_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_custom_bot_category_action.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
