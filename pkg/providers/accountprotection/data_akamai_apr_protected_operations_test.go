package accountprotection

import (
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataAprProtectedOperations(t *testing.T) {
	t.Run("TestDataAprProtectedOperations", func(t *testing.T) {

		mockedAprClient := &apr.Mock{}
		response := apr.ListProtectedOperationsResponse{
			Metadata: apr.Metadata{ConfigID: 43253, ConfigVersion: 15, SecurityPolicyID: "AAAA_81230"},
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
					"metadata": {
						"configId" : 43253,
						"configVersion" : 15,
						"securityPolicyId" : "AAAA_81230"
					},
					"operations":[
						{"operationId":"b85e3eaa-d334-466d-857e-33308ce416be", "testKey":"testValue1"},
						{"operationId":"69acad64-7459-4c1d-9bad-672600150127", "testKey":"testValue2"},
						{"operationId":"cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "testKey":"testValue3"},
						{"operationId":"10c54ea3-e3cb-4fc0-b0e0-fa3658aebd7b", "testKey":"testValue4"},
						{"operationId":"4d64d85a-a07f-485a-bbac-24c60658a1b8", "testKey":"testValue5"}
					]
				}`
		mockedAprClient.On("ListProtectedOperations",
			testutils.MockContext,
			apr.ListProtectedOperationsRequest{ConfigID: 43253, Version: 15, SecurityPolicyID: "AAAA_81230"},
		).Return(&response, nil)

		useClient(mockedAprClient, func() {

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDataProtectedOperations/basic.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_apr_protected_operations.test", "json", compactJSON(expectedJSON))),
					},
				},
			})
		})

		mockedAprClient.AssertExpectations(t)
	})
}
