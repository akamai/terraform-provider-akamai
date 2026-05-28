package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceBotCategoryException(t *testing.T) {
	t.Parallel()
	t.Run("ResourceBotCategoryException", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client)
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		client.BotMan.On("UpdateBotCategoryException",
			testutils.MockContext,
			botman.UpdateBotCategoryExceptionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      createRequest,
			},
		).Return(createResponse, nil).Once()

		client.BotMan.On("GetBotCategoryException",
			testutils.MockContext,
			botman.GetBotCategoryExceptionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/update.json")
		client.BotMan.On("UpdateBotCategoryException",
			testutils.MockContext,
			botman.UpdateBotCategoryExceptionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      updateRequest,
			},
		).Return(updateResponse, nil).Once()

		client.BotMan.On("GetBotCategoryException",
			testutils.MockContext,
			botman.GetBotCategoryExceptionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceBotCategoryException/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_bot_category_exception.test", "id", "43253:AAAA_81230"),
						resource.TestCheckResourceAttr("akamai_botman_bot_category_exception.test", "bot_category_exception", expectedCreateJSON)),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceBotCategoryException/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_bot_category_exception.test", "id", "43253:AAAA_81230"),
						resource.TestCheckResourceAttr("akamai_botman_bot_category_exception.test", "bot_category_exception", expectedUpdateJSON)),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
