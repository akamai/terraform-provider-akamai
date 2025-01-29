package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataBotDetection(t *testing.T) {
	t.Run("DataBotDetection", func(t *testing.T) {
		mockedBotmanClient := &botman.Mock{}

		response := botman.GetBotDetectionListResponse{
			Detections: []map[string]interface{}{
				{"detectionId": "b85e3eaa-d334-466d-857e-33308ce416be", "detectionName": "Test name 1", "testKey": "testValue1"},
				{"detectionId": "69acad64-7459-4c1d-9bad-672600150127", "detectionName": "Test name 2", "testKey": "testValue2"},
				{"detectionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "detectionName": "Test name 3", "testKey": "testValue3"},
				{"detectionId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "detectionName": "Test name 4", "testKey": "testValue4"},
				{"detectionId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "detectionName": "Test name 5", "testKey": "testValue5"},
			},
		}
		expectedJSON := `
{
	"detections":[
		{"detectionId":"b85e3eaa-d334-466d-857e-33308ce416be", "detectionName": "Test name 1", "testKey":"testValue1"},
		{"detectionId":"69acad64-7459-4c1d-9bad-672600150127", "detectionName": "Test name 2", "testKey":"testValue2"},
		{"detectionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "detectionName": "Test name 3", "testKey":"testValue3"},
		{"detectionId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "detectionName": "Test name 4", "testKey":"testValue4"},
		{"detectionId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "detectionName": "Test name 5", "testKey":"testValue5"}
	]
}`
		mockedBotmanClient.On("GetBotDetectionList",
			testutils.MockContext,
			botman.GetBotDetectionListRequest{},
		).Return(&response, nil)
		useClient(mockedBotmanClient, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataBotDetection/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_bot_detection.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
	t.Run("DataBotDetection filter by BotName", func(t *testing.T) {
		mockedBotmanClient := &botman.Mock{}

		response := botman.GetBotDetectionListResponse{
			Detections: []map[string]interface{}{
				{"detectionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "detectionName": "Test name 3", "testKey": "testValue3"},
			},
		}
		expectedJSON := `
{
	"detections":[
		{"detectionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "detectionName": "Test name 3", "testKey":"testValue3"}
	]
}`
		mockedBotmanClient.On("GetBotDetectionList",
			testutils.MockContext,
			botman.GetBotDetectionListRequest{DetectionName: "Test name 3"},
		).Return(&response, nil)
		useClient(mockedBotmanClient, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataBotDetection/filter_by_name.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_bot_detection.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
