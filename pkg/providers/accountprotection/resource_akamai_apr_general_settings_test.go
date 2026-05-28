package accountprotection

import (
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceAprGeneralSettings(t *testing.T) {
	t.Parallel()
	t.Run("TestResourceAprGeneralSettings", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client)
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		client.AccountProtection.On("UpsertGeneralSettings",
			testutils.MockContext,
			apr.UpsertGeneralSettingsRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      createRequest,
			},
		).Return(createResponse, nil).Once()

		client.AccountProtection.On("GetGeneralSettings",
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
		client.AccountProtection.On("UpsertGeneralSettings",
			testutils.MockContext,
			apr.UpsertGeneralSettingsRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      updateRequest,
			},
		).Return(updateResponse, nil).Once()

		client.AccountProtection.On("GetGeneralSettings",
			testutils.MockContext,
			apr.GetGeneralSettingsRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
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

		client.AccountProtection.AssertExpectations(t)
	})
}
