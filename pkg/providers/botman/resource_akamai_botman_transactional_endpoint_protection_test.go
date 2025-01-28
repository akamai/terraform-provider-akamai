package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceTransactionalEndpointProtection(t *testing.T) {
	t.Run("ResourceTransactionalEndpointProtection", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/create.json")
		mockedBotmanClient.On("UpdateTransactionalEndpointProtection",
			testutils.MockContext,
			botman.UpdateTransactionalEndpointProtectionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedBotmanClient.On("GetTransactionalEndpointProtection",
			testutils.MockContext,
			botman.GetTransactionalEndpointProtectionRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := testutils.LoadFixtureBytes(t, "testdata/JsonPayload/update.json")
		mockedBotmanClient.On("UpdateTransactionalEndpointProtection",
			testutils.MockContext,
			botman.UpdateTransactionalEndpointProtectionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: updateRequest,
			},
		).Return(updateResponse, nil).Once()

		mockedBotmanClient.On("GetTransactionalEndpointProtection",
			testutils.MockContext,
			botman.GetTransactionalEndpointProtectionRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
