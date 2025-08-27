package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataChallengeAction(t *testing.T) {
	t.Run("DataChallengeAction", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.GetChallengeActionListResponse{
			ChallengeActions: []map[string]interface{}{
				{"actionId": "b85e3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
				{"actionId": "69acad64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
				{"actionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				{"actionId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
				{"actionId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
			},
		}
		expectedJSON := `
{
	"challengeActions":[
		{"actionId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"actionId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"actionId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"actionId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`
		mockedBotmanClient.On("GetChallengeActionList",
			testutils.MockContext,
			botman.GetChallengeActionListRequest{ConfigID: 43253, Version: 15},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataChallengeAction/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_challenge_action.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
	t.Run("DataChallengeAction filter by id", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.GetChallengeActionListResponse{
			ChallengeActions: []map[string]interface{}{
				{"actionId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			},
		}
		expectedJSON := `
{
	"challengeActions":[
		{"actionId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}
	]
}`
		mockedBotmanClient.On("GetChallengeActionList",
			testutils.MockContext,
			botman.GetChallengeActionListRequest{ConfigID: 43253, Version: 15, ActionID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataChallengeAction/filter_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_challenge_action.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
