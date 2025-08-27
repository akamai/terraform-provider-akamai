package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceRecategorizedAkamaiDefinedBot(t *testing.T) {
	t.Run("ResourceRecategorizedAkamaiDefinedBot", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createResponse := botman.RecategorizedAkamaiDefinedBotResponse{BotID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", CategoryID: "87fb601b-4d30-4e0d-a74f-dc77e2b1bb74"}
		mockedBotmanClient.On("CreateRecategorizedAkamaiDefinedBot",
			testutils.MockContext,
			botman.CreateRecategorizedAkamaiDefinedBotRequest{
				ConfigID:   43253,
				Version:    15,
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "87fb601b-4d30-4e0d-a74f-dc77e2b1bb74",
			},
		).Return(&createResponse, nil).Once()

		mockedBotmanClient.On("GetRecategorizedAkamaiDefinedBot",
			testutils.MockContext,
			botman.GetRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(&createResponse, nil).Times(3)

		updateResponse := botman.RecategorizedAkamaiDefinedBotResponse{BotID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", CategoryID: "c43b638c-8f9a-4ea3-b1bd-3c82c96fefbf"}
		mockedBotmanClient.On("UpdateRecategorizedAkamaiDefinedBot",
			testutils.MockContext,
			botman.UpdateRecategorizedAkamaiDefinedBotRequest{
				ConfigID:   43253,
				Version:    15,
				BotID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				CategoryID: "c43b638c-8f9a-4ea3-b1bd-3c82c96fefbf",
			},
		).Return(&updateResponse, nil).Once()

		mockedBotmanClient.On("GetRecategorizedAkamaiDefinedBot",
			testutils.MockContext,
			botman.GetRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(&updateResponse, nil).Times(2)

		mockedBotmanClient.On("RemoveRecategorizedAkamaiDefinedBot",
			testutils.MockContext,
			botman.RemoveRecategorizedAkamaiDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(nil).Once()

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceRecategorizedAkamaiDefinedBot/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_recategorized_akamai_defined_bot.test", "id", "43253:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("akamai_botman_recategorized_akamai_defined_bot.test", "bot_id", "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("akamai_botman_recategorized_akamai_defined_bot.test", "category_id", "87fb601b-4d30-4e0d-a74f-dc77e2b1bb74")),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceRecategorizedAkamaiDefinedBot/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_recategorized_akamai_defined_bot.test", "id", "43253:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("akamai_botman_recategorized_akamai_defined_bot.test", "bot_id", "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("akamai_botman_recategorized_akamai_defined_bot.test", "category_id", "c43b638c-8f9a-4ea3-b1bd-3c82c96fefbf")),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
