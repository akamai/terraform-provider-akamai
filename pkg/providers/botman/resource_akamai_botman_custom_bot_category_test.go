package botman

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceCustomBotCategory(t *testing.T) {
	t.Parallel()
	t.Run("ResourceCustomBotCategory", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createResponse := map[string]interface{}{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		client.BotMan.On("CreateCustomBotCategory",
			testutils.MockContext,
			botman.CreateCustomBotCategoryRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		client.BotMan.On("GetCustomBotCategory",
			testutils.MockContext,
			botman.GetCustomBotCategoryRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "updated_testValue3"}
		updateRequest := `{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7","testKey":"updated_testValue3"}`
		client.BotMan.On("UpdateCustomBotCategory",
			testutils.MockContext,
			botman.UpdateCustomBotCategoryRequest{
				ConfigID:    43253,
				Version:     15,
				CategoryID:  "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(updateRequest),
			},
		).Return(updateResponse, nil).Once()

		client.BotMan.On("GetCustomBotCategory",
			testutils.MockContext,
			botman.GetCustomBotCategoryRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		client.BotMan.On("RemoveCustomBotCategory",
			testutils.MockContext,
			botman.RemoveCustomBotCategoryRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(nil).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomBotCategory/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_custom_bot_category.test", "id", "43253:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_custom_bot_category.test", "custom_bot_category", expectedCreateJSON)),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceCustomBotCategory/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_custom_bot_category.test", "id", "43253:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_custom_bot_category.test", "custom_bot_category", expectedUpdateJSON)),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
