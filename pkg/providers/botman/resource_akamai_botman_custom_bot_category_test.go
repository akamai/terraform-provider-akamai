package botman

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceCustomBotCategory(t *testing.T) {
	t.Run("ResourceCustomBotCategory", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createResponse := map[string]interface{}{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		mockedBotmanClient.On("CreateCustomBotCategory",
			testutils.MockContext,
			botman.CreateCustomBotCategoryRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedBotmanClient.On("GetCustomBotCategory",
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
		mockedBotmanClient.On("UpdateCustomBotCategory",
			testutils.MockContext,
			botman.UpdateCustomBotCategoryRequest{
				ConfigID:    43253,
				Version:     15,
				CategoryID:  "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(updateRequest),
			},
		).Return(updateResponse, nil).Once()

		mockedBotmanClient.On("GetCustomBotCategory",
			testutils.MockContext,
			botman.GetCustomBotCategoryRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		mockedBotmanClient.On("RemoveCustomBotCategory",
			testutils.MockContext,
			botman.RemoveCustomBotCategoryRequest{
				ConfigID:   43253,
				Version:    15,
				CategoryID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(nil).Once()

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
