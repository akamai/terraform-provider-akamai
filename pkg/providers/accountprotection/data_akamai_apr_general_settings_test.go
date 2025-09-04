package accountprotection

import (
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataGeneralSettings(t *testing.T) {
	t.Run("TestDataGeneralSettings", func(t *testing.T) {

		mockedAprClient := &apr.Mock{}
		response := map[string]interface{}{"testKey": "testValue3"}
		expectedJSON := `{"testKey":"testValue3"}`
		mockedAprClient.On("GetGeneralSettings",
			testutils.MockContext,
			apr.GetGeneralSettingsRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230"},
		).Return(response, nil)

		useClient(mockedAprClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataAprGeneralSettings/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_apr_general_settings.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedAprClient.AssertExpectations(t)
	})
}
