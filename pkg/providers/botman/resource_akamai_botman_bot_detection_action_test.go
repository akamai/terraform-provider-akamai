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

func TestResourceBotDetectionAction(t *testing.T) {
	t.Parallel()
	t.Run("ResourceBotDetectionAction", func(t *testing.T) {
		t.Parallel()

		expectedCreateJSON := `{"testKey":"testValue3"}`
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`
		client := setupMockedBotDetectionActionBotmanClient(false)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceBotDetectionAction/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_bot_detection_action.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_bot_detection_action.test", "bot_detection_action", expectedCreateJSON)),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceBotDetectionAction/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_bot_detection_action.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_bot_detection_action.test", "bot_detection_action", expectedUpdateJSON)),
				},
			},
		})
	})
	t.Run("ResourceBotDetectionActionDeletedFromRemoteWithoutCache", func(t *testing.T) {
		t.Parallel()

		expectedCreateJSON := `{"testKey":"testValue3"}`
		client := setupMockedBotDetectionActionBotmanClient(true)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceBotDetectionAction/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_bot_detection_action.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_bot_detection_action.test", "bot_detection_action", expectedCreateJSON)),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceBotDetectionAction/update_when_resource_deleted_outside_TF_without_cache.tf"),
					ExpectError: regexp.MustCompile(`Bot detection with id \[cc9c3f89-e179-4892-89cf-d5e623ba9dc7] does not exist`),
				},
			},
		})
	})
	t.Run("ResourceBotDetectionActionDeletedFromRemoteWithCache", func(t *testing.T) {
		t.Parallel()

		expectedCreateJSON := `{"testKey":"testValue3"}`
		client := setupMockedBotDetectionActionBotmanClient(true)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceBotDetectionAction/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_bot_detection_action.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_bot_detection_action.test", "bot_detection_action", expectedCreateJSON)),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResourceBotDetectionAction/update_when_resource_deleted_outside_TF_with_cache.tf"),
					ExpectError: regexp.MustCompile(`Bot detection with id \[cc9c3f89-e179-4892-89cf-d5e623ba9dc7] does not exist`),
				},
			},
		})
	})
}
func setupMockedBotDetectionActionBotmanClient(errScenario bool) *edgegrid.TestClient {
	client := edgegrid.NewTestClient()
	mockGetConfigVersion(client.APPSEC)
	createResponse := map[string]interface{}{"detectionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
	createRequest := `{"detectionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`
	client.BotMan.On("UpdateBotDetectionAction",
		testutils.MockContext,
		botman.UpdateBotDetectionActionRequest{
			ConfigID:         43253,
			Version:          15,
			SecurityPolicyID: "AAAA_81230",
			DetectionID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			JsonPayload:      json.RawMessage(compactJSON(createRequest)),
		},
	).Return(createResponse, nil).Once()

	client.BotMan.On("GetBotDetectionAction",
		testutils.MockContext,
		botman.GetBotDetectionActionRequest{
			ConfigID:         43253,
			Version:          15,
			SecurityPolicyID: "AAAA_81230",
			DetectionID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
	).Return(createResponse, nil).Times(2)

	if !errScenario {
		updateResponse := map[string]interface{}{"detectionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "updated_testValue3"}
		updateRequest := `{"detectionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"updated_testValue3"}`
		client.BotMan.On("UpdateBotDetectionAction",
			testutils.MockContext,
			botman.UpdateBotDetectionActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				DetectionID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(compactJSON(updateRequest)),
			},
		).Return(updateResponse, nil).Once()

		client.BotMan.On("GetBotDetectionAction",
			testutils.MockContext,
			botman.GetBotDetectionActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				DetectionID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(updateResponse, nil).Times(2)
	}
	updateRequest1 := `{"detectionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"updated_testValue4"}`
	err := fmt.Errorf("%s", "Title: Not Found; Type: https://problems.luna.akamaiapis.net/appsec/error-types/NOT-FOUND; Detail: Bot detection with id [cc9c3f89-e179-4892-89cf-d5e623ba9dc7] does not exist")
	client.BotMan.On("UpdateBotDetectionAction",
		testutils.MockContext,
		botman.UpdateBotDetectionActionRequest{
			ConfigID:         43253,
			Version:          15,
			SecurityPolicyID: "AAAA_81230",
			DetectionID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			JsonPayload:      json.RawMessage(compactJSON(updateRequest1)),
		},
	).Return(nil, err).Once()

	client.BotMan.On("GetBotDetectionAction",
		testutils.MockContext,
		botman.GetBotDetectionActionRequest{
			ConfigID:         43253,
			Version:          15,
			SecurityPolicyID: "AAAA_81230",
			DetectionID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
		},
	).Return(nil, err).Times(2)

	botDetectionActionListResponse := &botman.GetBotDetectionActionListResponse{
		Actions: []map[string]interface{}{
			{"detectionId": "cc9c3f91-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			{"detectionId": "cc9c3f90-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue2"},
		},
	}

	client.BotMan.On("GetBotDetectionActionList",
		testutils.MockContext,
		botman.GetBotDetectionActionListRequest{
			ConfigID:         43253,
			Version:          15,
			SecurityPolicyID: "AAAA_81230",
		},
	).Return(botDetectionActionListResponse, nil).Once()

	return client
}
