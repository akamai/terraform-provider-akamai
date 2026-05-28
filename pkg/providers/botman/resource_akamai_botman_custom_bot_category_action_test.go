package botman

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceCustomBotCategoryAction(t *testing.T) {
	t.Parallel()
	t.Run("ResourceCustomBotCategoryAction", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client)
		createResponse := map[string]interface{}{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
		createRequest := `{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`
		client.BotMan.On("UpdateCustomBotCategoryAction",
			testutils.MockContext,
			botman.UpdateCustomBotCategoryActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				CategoryID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(compactJSON(createRequest)),
			},
		).Return(createResponse, nil).Once()

		client.BotMan.On("GetCustomBotCategoryAction",
			testutils.MockContext,
			botman.GetCustomBotCategoryActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				CategoryID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "updated_testValue3"}
		updateRequest := `{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"updated_testValue3"}`
		client.BotMan.On("UpdateCustomBotCategoryAction",
			testutils.MockContext,
			botman.UpdateCustomBotCategoryActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				CategoryID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(compactJSON(updateRequest)),
			},
		).Return(updateResponse, nil).Once()

		client.BotMan.On("GetCustomBotCategoryAction",
			testutils.MockContext,
			botman.GetCustomBotCategoryActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				CategoryID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomBotCategoryAction/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_action.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_action.test", "custom_bot_category_action", expectedCreateJSON)),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomBotCategoryAction/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_action.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_custom_bot_category_action.test", "custom_bot_category_action", expectedUpdateJSON)),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
