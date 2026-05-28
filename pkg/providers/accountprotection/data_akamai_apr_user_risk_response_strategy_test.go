package accountprotection

import (
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataUserRiskResponseStrategy(t *testing.T) {
	t.Parallel()
	t.Run("TestDataUserRiskResponseStrategy", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		response := map[string]interface{}{"testKey": "testValue3"}
		expectedJSON := `{"testKey":"testValue3"}`
		client.AccountProtection.On("GetUserRiskResponseStrategy",
			testutils.MockContext,
			apr.GetUserRiskResponseStrategyRequest{ConfigID: 43253, Version: 15},
		).Return(response, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataUserRiskResponseStrategy/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_apr_user_risk_response_strategy.test", "json", compactJSON(expectedJSON))),
				},
			},
		})

		client.AccountProtection.AssertExpectations(t)
	})
}
