package botman

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceServeAlternateAction(t *testing.T) {
	t.Run("ResourceServeAlternateAction", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createResponse := map[string]interface{}{"actionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		mockedBotmanClient.On("CreateServeAlternateAction",
			testutils.MockContext,
			botman.CreateServeAlternateActionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedBotmanClient.On("GetServeAlternateAction",
			testutils.MockContext,
			botman.GetServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  15,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"actionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "updated_testValue3"}
		updateRequest := `{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7","testKey":"updated_testValue3"}`
		mockedBotmanClient.On("UpdateServeAlternateAction",
			testutils.MockContext,
			botman.UpdateServeAlternateActionRequest{
				ConfigID:    43253,
				Version:     15,
				ActionID:    "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload: json.RawMessage(updateRequest),
			},
		).Return(updateResponse, nil).Once()

		mockedBotmanClient.On("GetServeAlternateAction",
			testutils.MockContext,
			botman.GetServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  15,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		mockedBotmanClient.On("RemoveServeAlternateAction",
			testutils.MockContext,
			botman.RemoveServeAlternateActionRequest{
				ConfigID: 43253,
				Version:  15,
				ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(nil).Once()

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceServeAlternateAction/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_serve_alternate_action.test", "id", "43253:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("akamai_botman_serve_alternate_action.test", "serve_alternate_action", expectedCreateJSON)),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceServeAlternateAction/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_serve_alternate_action.test", "id", "43253:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
							resource.TestCheckResourceAttr("akamai_botman_serve_alternate_action.test", "serve_alternate_action", expectedUpdateJSON)),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
