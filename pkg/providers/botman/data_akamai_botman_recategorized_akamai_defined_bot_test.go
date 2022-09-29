package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataRecategorizedAkamaiDefinedBot(t *testing.T) {
	t.Run("DataRecategorizedAkamaiDefinedBot", func(t *testing.T) {

		mockedBotmanClient := &mockbotman{}
		response := botman.GetRecategorizedAkamaiDefinedBotListResponse{
			Bots: []botman.RecategorizedAkamaiDefinedBotResponse{
				{BotID: "b85e3eaa-d334-466d-857e-33308ce416be", CategoryID: "39cbadc6-c5ef-42d1-9290-7895f24316ad"},
				{BotID: "69acad64-7459-4c1d-9bad-672600150127", CategoryID: "5eb700c8-275d-4866-a271-b6fa25e1fdc5"},
				{BotID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e"},
				{BotID: "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", CategoryID: "87a152a9-8af0-4c4f-9c37-a895fe7ca6b4"},
				{BotID: "4d64d85a-a07f-485a-bbac-24c60658a1b8", CategoryID: "b61a3017-bff4-41b0-9396-be378d4f07c1"},
			},
		}
		expectedJSON := `
{
	"recategorizedBots": [
		{"botId":"b85e3eaa-d334-466d-857e-33308ce416be", "customBotCategoryId":"39cbadc6-c5ef-42d1-9290-7895f24316ad"},
		{"botId":"69acad64-7459-4c1d-9bad-672600150127", "customBotCategoryId":"5eb700c8-275d-4866-a271-b6fa25e1fdc5"},
		{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "customBotCategoryId":"0d38d0fe-b05d-42f6-a58f-bc98c821793e"},
		{"botId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "customBotCategoryId":"87a152a9-8af0-4c4f-9c37-a895fe7ca6b4"},
		{"botId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "customBotCategoryId":"b61a3017-bff4-41b0-9396-be378d4f07c1"}
	]
}`
		mockedBotmanClient.On("GetRecategorizedAkamaiDefinedBotList",
			mock.Anything,
			botman.GetRecategorizedAkamaiDefinedBotListRequest{ConfigID: 43253, Version: 15},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataRecategorizedAkamaiDefinedBot/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_recategorized_akamai_defined_bot.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
	t.Run("DataRecategorizedAkamaiDefinedBot filter by id", func(t *testing.T) {

		mockedBotmanClient := &mockbotman{}
		response := botman.GetRecategorizedAkamaiDefinedBotListResponse{
			Bots: []botman.RecategorizedAkamaiDefinedBotResponse{
				{BotID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", CategoryID: "0d38d0fe-b05d-42f6-a58f-bc98c821793e"},
			},
		}
		expectedJSON := `
{
	"recategorizedBots": [
		{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "customBotCategoryId":"0d38d0fe-b05d-42f6-a58f-bc98c821793e"}
	]
}`
		mockedBotmanClient.On("GetRecategorizedAkamaiDefinedBotList",
			mock.Anything,
			botman.GetRecategorizedAkamaiDefinedBotListRequest{ConfigID: 43253, Version: 15, BotID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataRecategorizedAkamaiDefinedBot/filter_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_recategorized_akamai_defined_bot.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
