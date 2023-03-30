package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceTransactionalEndpointProtection(t *testing.T) {
	t.Run("ResourceTransactionalEndpointProtection", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		createResponse := map[string]interface{}{"testKey": "testValue3"}
		createRequest := test.FixtureBytes("testdata/JsonPayload/create.json")
		mockedBotmanClient.On("UpdateTransactionalEndpointProtection",
			mock.Anything,
			botman.UpdateTransactionalEndpointProtectionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: createRequest,
			},
		).Return(createResponse, nil).Once()

		mockedBotmanClient.On("GetTransactionalEndpointProtection",
			mock.Anything,
			botman.GetTransactionalEndpointProtectionRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"testKey": "updated_testValue3"}
		updateRequest := test.FixtureBytes("testdata/JsonPayload/update.json")
		mockedBotmanClient.On("UpdateTransactionalEndpointProtection",
			mock.Anything,
			botman.UpdateTransactionalEndpointProtectionRequest{
				ConfigID:    43253,
				Version:     15,
				JsonPayload: updateRequest,
			},
		).Return(updateResponse, nil).Once()

		mockedBotmanClient.On("GetTransactionalEndpointProtection",
			mock.Anything,
			botman.GetTransactionalEndpointProtectionRequest{
				ConfigID: 43253,
				Version:  15,
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: test.Fixture("testdata/TestResourceTransactionalEndpointProtection/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_botman_transactional_endpoint_protection.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_botman_transactional_endpoint_protection.test", "transactional_endpoint_protection", expectedCreateJSON)),
					},
					{
						Config: test.Fixture("testdata/TestResourceTransactionalEndpointProtection/update.tf"),
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
