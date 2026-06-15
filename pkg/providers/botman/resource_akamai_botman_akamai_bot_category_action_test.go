package botman

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceAkamaiBotCategoryAction(t *testing.T) {
	t.Parallel()
	t.Run("ResourceAkamaiBotCategoryAction", func(t *testing.T) {
		t.Parallel()

		expectedCreateJSON := `{"testKey":"testValue3"}`
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`
		client := setupMockedAkamaiBotCategoryActionBotmanClient(false)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAkamaiBotCategoryAction/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_akamai_bot_category_action.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_akamai_bot_category_action.test", "akamai_bot_category_action", expectedCreateJSON)),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAkamaiBotCategoryAction/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_akamai_bot_category_action.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_akamai_bot_category_action.test", "akamai_bot_category_action", expectedUpdateJSON)),
				},
			},
		})
	})
	t.Run("ResourceAkamaiBotCategoryActionWhenDeletedFromRemoteWithoutCache", func(t *testing.T) {
		t.Parallel()

		client := setupMockedAkamaiBotCategoryActionBotmanClient(true)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAkamaiBotCategoryAction/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_akamai_bot_category_action.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_akamai_bot_category_action.test", "akamai_bot_category_action", expectedCreateJSON)),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAkamaiBotCategoryAction/update_when_resource_deleted_outside_TF_without_cache.tf"),
					ExpectError: regexp.MustCompile(`Akamai Bot Category with id \[cc9c3f89-e179-4892-89cf-d5e623ba9dc7] does not exist`),
					Check:       resource.TestCheckNoResourceAttr("akamai_botman_akamai_bot_category_action.test", "akamai_bot_category_action"),
				},
			},
		})
	})
	t.Run("ResourceAkamaiBotCategoryActionWhenDeletedFromRemoteWithCache", func(t *testing.T) {
		t.Parallel()

		expectedCreateJSON := `{"testKey":"testValue3"}`
		client := setupMockedAkamaiBotCategoryActionBotmanClient(true)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceAkamaiBotCategoryAction/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_akamai_bot_category_action.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_akamai_bot_category_action.test", "akamai_bot_category_action", expectedCreateJSON)),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceAkamaiBotCategoryAction/update_when_resource_deleted_outside_TF_with_cache.tf"),
					ExpectError: regexp.MustCompile(`Akamai Bot Category with id \[cc9c3f89-e179-4892-89cf-d5e623ba9dc7] does not exist`),
					Check:       resource.TestCheckNoResourceAttr("akamai_botman_akamai_bot_category_action.test", "akamai_bot_category_action"),
				},
			},
		})
	})
}

func setupMockedAkamaiBotCategoryActionBotmanClient(errScenario bool) *edgegrid.TestClient {
	client := edgegrid.NewTestClient()
	mockGetConfigVersion(client.APPSEC)
	createResponse := map[string]interface{}{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
	createRequest := `{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`
	client.BotMan.On("UpdateAkamaiBotCategoryAction",
		testutils.MockContext,
		botman.UpdateAkamaiBotCategoryActionRequest{
			ConfigID:         43253,
			Version:          15,
			SecurityPolicyID: "AAAA_81230",
			CategoryID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			JsonPayload:      json.RawMessage(compactJSON(createRequest)),
		},
	).Return(createResponse, nil).Once()

	client.BotMan.On("GetAkamaiBotCategoryAction",
		testutils.MockContext,
		botman.GetAkamaiBotCategoryActionRequest{
			ConfigID:         43253,
			Version:          15,
			SecurityPolicyID: "AAAA_81230",
			CategoryID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
	).Return(createResponse, nil).Times(2)

	if !errScenario {
		updateResponse := map[string]interface{}{"categoryId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "updated_testValue3"}
		updateRequest := `{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"updated_testValue3"}`
		client.BotMan.On("UpdateAkamaiBotCategoryAction",
			testutils.MockContext,
			botman.UpdateAkamaiBotCategoryActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				CategoryID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(compactJSON(updateRequest)),
			},
		).Return(updateResponse, nil).Once()

		client.BotMan.On("GetAkamaiBotCategoryAction",
			testutils.MockContext,
			botman.GetAkamaiBotCategoryActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				CategoryID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(updateResponse, nil).Times(2)
	} else {
		err := fmt.Errorf("%s", "Title: Not Found; Type: https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND; Detail: Akamai Bot Category with id [cc9c3f89-e179-4892-89cf-d5e623ba9dc7] does not exist")

		updateRequest1 := `{"categoryId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"updated_testValue4"}`

		client.BotMan.On("GetAkamaiBotCategoryAction",
			testutils.MockContext,
			botman.GetAkamaiBotCategoryActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				CategoryID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(nil, err).Times(1)

		client.BotMan.On("UpdateAkamaiBotCategoryAction",
			testutils.MockContext,
			botman.UpdateAkamaiBotCategoryActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				CategoryID:       "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(compactJSON(updateRequest1)),
			},
		).Return(nil, err).Once()

		akamaiBotCategoryActionListResponse := &botman.GetAkamaiBotCategoryActionListResponse{
			Actions: []map[string]interface{}{
				{"categoryId": "cc9c3f91-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				{"categoryId": "cc9c3f90-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue2"},
			},
		}

		client.BotMan.On("GetAkamaiBotCategoryActionList",
			testutils.MockContext,
			botman.GetAkamaiBotCategoryActionListRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(akamaiBotCategoryActionListResponse, nil).Once()
	}
	return client
}
