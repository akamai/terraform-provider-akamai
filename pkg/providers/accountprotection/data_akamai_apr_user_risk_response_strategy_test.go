package accountprotection

import (
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataUserRiskResponseStrategy(t *testing.T) {
	t.Run("TestDataUserRiskResponseStrategy", func(t *testing.T) {

		mockedAprClient := &apr.Mock{}
		response := map[string]interface{}{"testKey": "testValue3"}
		expectedJSON := `{"testKey":"testValue3"}`
		mockedAprClient.On("GetUserRiskResponseStrategy",
			testutils.MockContext,
			apr.GetUserRiskResponseStrategyRequest{ConfigID: 43253, Version: 15},
		).Return(response, nil)

		useClient(mockedAprClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataUserRiskResponseStrategy/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_apr_user_risk_response_strategy.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedAprClient.AssertExpectations(t)
	})
}
