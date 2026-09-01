package botman

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceTransactionalEndpoint(t *testing.T) {
	t.Parallel()
	t.Run("ResourceTransactionalEndpoint", func(t *testing.T) {
		t.Parallel()

		client := edgegrid.NewTestClient()
		mockGetConfigVersion(client.APPSEC)
		createResponse := map[string]interface{}{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"}
		createRequest := `{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7","testKey":"testValue3"}`
		client.BotMan.On("CreateTransactionalEndpoint",
			testutils.MockContext,
			botman.CreateTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				JsonPayload:      json.RawMessage(createRequest),
			},
		).Return(createResponse, nil).Once()

		client.BotMan.On("GetTransactionalEndpoint",
			testutils.MockContext,
			botman.GetTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(createResponse, nil).Times(3)
		expectedCreateJSON := `{"testKey":"testValue3"}`

		updateResponse := map[string]interface{}{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "updated_testValue3"}
		updateRequest := `{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7","testKey":"updated_testValue3"}`
		client.BotMan.On("UpdateTransactionalEndpoint",
			testutils.MockContext,
			botman.UpdateTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
				JsonPayload:      json.RawMessage(updateRequest),
			},
		).Return(updateResponse, nil).Once()

		client.BotMan.On("GetTransactionalEndpoint",
			testutils.MockContext,
			botman.GetTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(updateResponse, nil).Times(2)
		expectedUpdateJSON := `{"testKey":"updated_testValue3"}`

		client.BotMan.On("RemoveTransactionalEndpoint",
			testutils.MockContext,
			botman.RemoveTransactionalEndpointRequest{
				ConfigID:         43253,
				Version:          15,
				SecurityPolicyID: "AAAA_81230",
				OperationID:      "cc9c3f89-e179-4892-89cf-d5e623ba9dc7",
			},
		).Return(nil).Once()

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceTransactionalEndpoint/create.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_transactional_endpoint.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_transactional_endpoint.test", "transactional_endpoint", expectedCreateJSON)),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResourceTransactionalEndpoint/update.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_botman_transactional_endpoint.test", "id", "43253:AAAA_81230:cc9c3f89-e179-4892-89cf-d5e623ba9dc7"),
						resource.TestCheckResourceAttr("akamai_botman_transactional_endpoint.test", "transactional_endpoint", expectedUpdateJSON)),
				},
			},
		})

		client.BotMan.AssertExpectations(t)
	})
}
