package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceTransactionalEndpointProtection(t *testing.T) {
	t.Parallel()
	t.Run("ResourceTransactionalEndpointProtection", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		client.BotMan.On("UpdateTransactionalEndpointProtection",
			testutils.MockContext,
			botman.UpdateTransactionalEndpointProtectionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		client.BotMan.On("GetTransactionalEndpointProtection",
			testutils.MockContext,
			botman.GetTransactionalEndpointProtectionRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/update.json")
		client.BotMan.On("UpdateTransactionalEndpointProtection",
			testutils.MockContext,
			botman.UpdateTransactionalEndpointProtectionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: updateRequest,
			},
		).Return(updateResponse, nil).Once()

		client.BotMan.On("GetTransactionalEndpointProtection",
			testutils.MockContext,
			botman.GetTransactionalEndpointProtectionRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceTransactionalEndpointProtection/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_transactional_endpoint_protection.test", "id", "43253"),
						resource.TestCheckResourceAttr("akamai_botman_transactional_endpoint_protection.test", "transactional_endpoint_protection", expectedCreateJSON)),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceTransactionalEndpointProtection/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_transactional_endpoint_protection.test", "id", "43253"),
						resource.TestCheckResourceAttr("akamai_botman_transactional_endpoint_protection.test", "transactional_endpoint_protection", expectedUpdateJSON)),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
