package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceChallengeInterceptionRules(t *testing.T) {
	t.Run("ResourceChallengeInterceptionRules", func(t *testing.T) {

		mockedBotmanClient := &mockbotman{}
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := test.FixtureBytes("testdata/JsonPayload/create.json")
		mockedBotmanClient.On("UpdateChallengeInterceptionRules",
			mock.Anything,
			botman.UpdateChallengeInterceptionRulesRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedBotmanClient.On("GetChallengeInterceptionRules",
			mock.Anything,
			botman.GetChallengeInterceptionRulesRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := test.FixtureBytes("testdata/JsonPayload/update.json")
		mockedBotmanClient.On("UpdateChallengeInterceptionRules",
			mock.Anything,
			botman.UpdateChallengeInterceptionRulesRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: updateRequest,
			},
		).Return(updateResponse, nil).Once()

		mockedBotmanClient.On("GetChallengeInterceptionRules",
			mock.Anything,
			botman.GetChallengeInterceptionRulesRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestResourceChallengeInterceptionRules/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_challenge_interception_rules.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_challenge_interception_rules.test", "challenge_interception_rules", expectedCreateJSON)),
					},
					{
						Config: test.Fixture("testdata/TestResourceChallengeInterceptionRules/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_challenge_interception_rules.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_challenge_interception_rules.test", "challenge_interception_rules", expectedUpdateJSON)),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
