package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceChallengeInjectionRules(t *testing.T) {
	t.Run("ResourceChallengeInjectionRules", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := test.FixtureBytes("testdata/JsonPayload/create.json")
		mockedBotmanClient.On("UpdateChallengeInjectionRules",
			mock.Anything,
			botman.UpdateChallengeInjectionRulesRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedBotmanClient.On("GetChallengeInjectionRules",
			mock.Anything,
			botman.GetChallengeInjectionRulesRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := test.FixtureBytes("testdata/JsonPayload/update.json")
		mockedBotmanClient.On("UpdateChallengeInjectionRules",
			mock.Anything,
			botman.UpdateChallengeInjectionRulesRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: updateRequest,
			},
		).Return(updateResponse, nil).Once()

		mockedBotmanClient.On("GetChallengeInjectionRules",
			mock.Anything,
			botman.GetChallengeInjectionRulesRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestResourceChallengeInjectionRules/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_challenge_injection_rules.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_challenge_injection_rules.test", "challenge_injection_rules", expectedCreateJSON)),
					},
					{
						Config: test.Fixture("testdata/TestResourceChallengeInjectionRules/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_challenge_injection_rules.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_challenge_injection_rules.test", "challenge_injection_rules", expectedUpdateJSON)),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
