package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataChallengeInjectionRules(t *testing.T) {
	t.Run("DataChallengeInjectionRules", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := map[string]interface{}{"testKey": "testValue3"}
		expectedJSON := `{"testKey":"testValue3"}`
		mockedBotmanClient.On("GetChallengeInjectionRules",
			testutils.MockContext,
			botman.GetChallengeInjectionRulesRequest{ConfigID: 43253, Version: 15},
		).Return(response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataChallengeInjectionRules/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_challenge_injection_rules.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
