package accountprotection

import (
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceUserRiskResponseStrategy(t *testing.T) {
	t.Parallel()
	t.Run("TestResourceUserRiskResponseStrategy", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		client.AccountProtection.On("UpsertUserRiskResponseStrategy",
			testutils.MockContext,
			apr.UpsertUserRiskResponseStrategyRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		client.AccountProtection.On("GetUserRiskResponseStrategy",
			testutils.MockContext,
			apr.GetUserRiskResponseStrategyRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/update.json")
		client.AccountProtection.On("UpsertUserRiskResponseStrategy",
			testutils.MockContext,
			apr.UpsertUserRiskResponseStrategyRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: updateRequest,
			},
		).Return(updateResponse, nil).Once()

		client.AccountProtection.On("GetUserRiskResponseStrategy",
			testutils.MockContext,
			apr.GetUserRiskResponseStrategyRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceUserRiskResponseStrategy/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apr_user_risk_response_strategy.test", "id", "43253"),
						resource.TestCheckResourceAttr("akamai_apr_user_risk_response_strategy.test", "user_risk_response_strategy", expectedCreateJSON)),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceUserRiskResponseStrategy/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apr_user_risk_response_strategy.test", "id", "43253"),
						resource.TestCheckResourceAttr("akamai_apr_user_risk_response_strategy.test", "user_risk_response_strategy", expectedUpdateJSON)),
				},
			},
		})

		client.AccountProtection.AssertExpectations(t)
	})
}
