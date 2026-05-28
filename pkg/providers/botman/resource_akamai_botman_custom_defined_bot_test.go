package botman

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceCustomDefinedBot(t *testing.T) {
	t.Parallel()
	t.Run("ResourceCustomDefinedBot", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client)
		createResponse := map[string]interface{}{"botId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		client.BotMan.On("CreateCustomDefinedBot",
			testutils.MockContext,
			botman.CreateCustomDefinedBotRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		client.BotMan.On("GetCustomDefinedBot",
			testutils.MockContext,
			botman.GetCustomDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"botId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "updated_testValue3"}
		updateRequest := `{"botId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7","testKey":"updated_testValue3"}`
		client.BotMan.On("UpdateCustomDefinedBot",
			testutils.MockContext,
			botman.UpdateCustomDefinedBotRequest{
				ConfigID:    43253,
				Version:     15,
				BotID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(updateRequest),
			},
		).Return(updateResponse, nil).Once()

		client.BotMan.On("GetCustomDefinedBot",
			testutils.MockContext,
			botman.GetCustomDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		client.BotMan.On("RemoveCustomDefinedBot",
			testutils.MockContext,
			botman.RemoveCustomDefinedBotRequest{
				ConfigID: 43253,
				Version:  15,
				BotID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(nil).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomDefinedBot/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_custom_defined_bot.test", "id", "43253:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_custom_defined_bot.test", "custom_defined_bot", expectedCreateJSON)),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomDefinedBot/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_custom_defined_bot.test", "id", "43253:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_custom_defined_bot.test", "custom_defined_bot", expectedUpdateJSON)),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
