package botman

import (
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataTransactionalEndpoint(t *testing.T) {
	t.Run("DataTransactionalEndpoint", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.GetTransactionalEndpointListResponse{
			Operations: []map[string]interface{}{
				{"operationId": "b85e3eaa-d334-466d-857e-33308ce416be", "testKey": "testValue1"},
				{"operationId": "69acad64-7459-4c1d-9bad-672600150127", "testKey": "testValue2"},
				{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
				{"operationId": "10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey": "testValue4"},
				{"operationId": "4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey": "testValue5"},
			},
		}
		expectedJSON := `
{
	"operations":[
		{"operationId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
		{"operationId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
		{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
		{"operationId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
		{"operationId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
	]
}`
		mockedBotmanClient.On("GetTransactionalEndpointList",
			testutils.MockContext,
			botman.GetTransactionalEndpointListRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230"},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataTransactionalEndpoint/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_transactional_endpoint.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
	t.Run("DataTransactionalEndpoint filter by id", func(t *testing.T) {

		mockedBotmanClient := &botman.Mock{}
		response := botman.GetTransactionalEndpointListResponse{
			Operations: []map[string]interface{}{
				{"operationId": "cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey": "testValue3"},
			},
		}
		expectedJSON := `
{
	"operations":[
		{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"}
	]
}`
		mockedBotmanClient.On("GetTransactionalEndpointList",
			testutils.MockContext,
			botman.GetTransactionalEndpointListRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230", OperationID: "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"},
		).Return(&response, nil)

		useClient(mockedBotmanClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataTransactionalEndpoint/filter_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_botman_transactional_endpoint.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedBotmanClient.AssertExpectations(t)
	})
}
