package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataChallengeInterceptionRules(t *testing.T) {
	t.Run("DataChallengeInterceptionRules", func(t *testing.T) {

		mockedBotmanClient := &mockbotman{}
		response := map[string]interface{}{"testKey": "testValue3"}
		expectedJSON := `{"testKey":"testValue3"}`
		mockedBotmanClient.On("GetChallengeInterceptionRules",
			mock.Anything,
			botman.GetChallengeInterceptionRulesRequest{ConfigID: 43253, Version: 15},
		).Return(response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestDataChallengeInterceptionRules/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_challenge_interception_rules.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
