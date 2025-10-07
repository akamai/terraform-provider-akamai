package botman

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceBotDetectionAction(t *testing.T) {
	t.Run("ResourceBotDetectionAction", func(t *testing.T) {

		expectedCreateJSON := `{"testKey":"testValue3"}`
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`
		mockedBotmanClient := setupMockedBotDetectionActionBotmanClient(false)
		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
	})
	t.Run("ResourceBotDetectionActionDeletedFromRemoteWithoutCache", func(t *testing.T) {

		expectedCreateJSON := `{"testKey":"testValue3"}`
		mockedBotmanClient := setupMockedBotDetectionActionBotmanClient(true)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
	})
	t.Run("ResourceBotDetectionActionDeletedFromRemoteWithCache", func(t *testing.T) {

		expectedCreateJSON := `{"testKey":"testValue3"}`
		mockedBotmanClient := setupMockedBotDetectionActionBotmanClient(true)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
	})
}
func setupMockedBotDetectionActionBotmanClient(errScenario bool) *botman.Mock {
	mockedBotmanClient := &botman.Mock{}
	createResponse := map[string]interface{}{"detectionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
	createRequest := `{"detectionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}`
	mockedBotmanClient.On("UpdateBotDetectionAction",
		testutils.MockContext,
		botman.UpdateBotDetectionActionRequest{
			ConfigID:         43253,
			Version:          15,
			SecurityPolicyID: "AAAA_81230",
			DetectionID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			JsonPayload:      json.RawMessage(compactJSON(createRequest)),
		},
	).Return(createResponse, nil).Once()

	mockedBotmanClient.On("GetBotDetectionAction",
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
		mockedBotmanClient.On("UpdateBotDetectionAction",
			testutils.MockContext,
			botman.UpdateBotDetectionActionRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				DetectionID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(compactJSON(updateRequest)),
			},
		).Return(updateResponse, nil).Once()

		mockedBotmanClient.On("GetBotDetectionAction",
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
	mockedBotmanClient.On("UpdateBotDetectionAction",
		testutils.MockContext,
		botman.UpdateBotDetectionActionRequest{
			ConfigID:         43253,
			Version:          15,
			SecurityPolicyID: "AAAA_81230",
			DetectionID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			JsonPayload:      json.RawMessage(compactJSON(updateRequest1)),
		},
	).Return(nil, err).Once()

	mockedBotmanClient.On("GetBotDetectionAction",
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

	mockedBotmanClient.On("GetBotDetectionActionList",
		testutils.MockContext,
		botman.GetBotDetectionActionListRequest{
			ConfigID:         43253,
			Version:          15,
			SecurityPolicyID: "AAAA_81230",
		},
	).Return(botDetectionActionListResponse, nil).Once()

	return mockedBotmanClient
}
