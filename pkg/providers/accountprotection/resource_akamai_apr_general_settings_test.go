package accountprotection

import (
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceAprGeneralSettings(t *testing.T) {
	t.Run("TestResourceAprGeneralSettings", func(t *testing.T) {

		mockedAprClient := &apr.Mock{}
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		mockedAprClient.On("UpsertGeneralSettings",
			testutils.MockContext,
			apr.UpsertGeneralSettingsRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedAprClient.On("GetGeneralSettings",
			testutils.MockContext,
			apr.GetGeneralSettingsRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/update.json")
		mockedAprClient.On("UpsertGeneralSettings",
			testutils.MockContext,
			apr.UpsertGeneralSettingsRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      updateRequest,
			},
		).Return(updateResponse, nil).Once()

		mockedAprClient.On("GetGeneralSettings",
			testutils.MockContext,
			apr.GetGeneralSettingsRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		useClient(mockedAprClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceAprGeneralSettings/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_apr_general_settings.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_apr_general_settings.test", "general_settings", expectedCreateJSON)),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceAprGeneralSettings/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_apr_general_settings.test", "id", "43253:AAAA_81230"),
							resource.TestCheckResourceAttr("akamai_apr_general_settings.test", "general_settings", expectedUpdateJSON)),
					},
				},
			})
		})

		mockedAprClient.AssertExpectations(t)
	})
}
