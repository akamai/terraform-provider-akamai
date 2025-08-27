package accountprotection

import (
	"testing"

	apr "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/accountprotection"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceProtectedOperations(t *testing.T) {
	const (
		configID         = 43253
		configVersion    = 15
		securityPolicyID = "AAAA_81230"
		operationID      = "b85e3eaa-d334-466d-857e-33308ce416be"

		createPayloadJSON = `{
			"operations": [
				{
					"operationId": "b85e3eaa-d334-466d-857e-33308ce416be",
					"testKey": "testValue"
				}
			]
		}`

		updatePayloadJSON = `{"testKey": "testValueUpdated"}`
	)

	t.Run("happy path", func(t *testing.T) {
		readResponse1 := apr.ListProtectedOperationsResponse{
			Metadata: apr.Metadata{ConfigID: configID, ConfigVersion: configVersion, SecurityPolicyID: securityPolicyID},
			Operations: []map[string]any{
				{"operationId": operationID, "testKey": "testValue"},
			},
		}
		readResponse2 := apr.ListProtectedOperationsResponse{
			Metadata: apr.Metadata{ConfigID: configID, ConfigVersion: configVersion, SecurityPolicyID: securityPolicyID},
			Operations: []map[string]any{
				{"operationId": operationID, "testKey": "testValue"},
			},
		}
		readResponse3 := apr.ListProtectedOperationsResponse{
			Metadata: apr.Metadata{ConfigID: configID, ConfigVersion: configVersion, SecurityPolicyID: securityPolicyID},
			Operations: []map[string]any{
				{"operationId": operationID, "testKey": "testValue"},
			},
		}
		readResponse4 := apr.ListProtectedOperationsResponse{
			Metadata: apr.Metadata{ConfigID: configID, ConfigVersion: configVersion, SecurityPolicyID: securityPolicyID},
			Operations: []map[string]any{
				{"operationId": operationID, "testKey": "testValueUpdated"},
			},
		}
		readResponse5 := apr.ListProtectedOperationsResponse{
			Metadata: apr.Metadata{ConfigID: configID, ConfigVersion: configVersion, SecurityPolicyID: securityPolicyID},
			Operations: []map[string]any{
				{"operationId": operationID, "testKey": "testValueUpdated"},
			},
		}

		clientMock := &apr.Mock{}
		clientMock.On("CreateProtectedOperations",
			testutils.MockContext,
			apr.CreateProtectedOperationsRequest{
				ConfigID:         configID,
				Version:          configVersion,
				SecurityPolicyID: securityPolicyID,
				JsonPayload:      []byte(compactJSON(createPayloadJSON)),
			},
		).Return(&readResponse1, nil).Once()
		clientMock.On("GetProtectedOperationByID",
			testutils.MockContext,
			apr.GetProtectedOperationByIDRequest{
				ConfigID:         configID,
				Version:          configVersion,
				SecurityPolicyID: securityPolicyID,
				OperationID:      operationID,
			},
		).Return(&readResponse1, nil).Once()
		clientMock.On("GetProtectedOperationByID",
			testutils.MockContext,
			apr.GetProtectedOperationByIDRequest{
				ConfigID:         configID,
				Version:          configVersion,
				SecurityPolicyID: securityPolicyID,
				OperationID:      operationID,
			},
		).Return(&readResponse2, nil).Once()
		clientMock.On("GetProtectedOperationByID",
			testutils.MockContext,
			apr.GetProtectedOperationByIDRequest{
				ConfigID:         configID,
				Version:          configVersion,
				SecurityPolicyID: securityPolicyID,
				OperationID:      operationID,
			},
		).Return(&readResponse3, nil).Once()
		clientMock.On("GetProtectedOperationByID",
			testutils.MockContext,
			apr.GetProtectedOperationByIDRequest{
				ConfigID:         configID,
				Version:          configVersion,
				SecurityPolicyID: securityPolicyID,
				OperationID:      operationID,
			},
		).Return(&readResponse4, nil).Once()
		clientMock.On("GetProtectedOperationByID",
			testutils.MockContext,
			apr.GetProtectedOperationByIDRequest{
				ConfigID:         configID,
				Version:          configVersion,
				SecurityPolicyID: securityPolicyID,
				OperationID:      operationID,
			},
		).Return(&readResponse5, nil).Once()
		clientMock.On("UpdateProtectedOperation",
			testutils.MockContext,
			apr.UpdateProtectedOperationRequest{
				ConfigID:         configID,
				Version:          configVersion,
				SecurityPolicyID: securityPolicyID,
				OperationID:      operationID,
				JsonPayload:      []byte(compactJSON(updatePayloadJSON)),
			},
		).Return(map[string]any{}, nil).Once()
		clientMock.On("RemoveProtectedOperation",
			testutils.MockContext,
			apr.RemoveProtectedOperationRequest{
				ConfigID:         configID,
				Version:          configVersion,
				SecurityPolicyID: securityPolicyID,
				OperationID:      operationID,
			},
		).Return(nil).Once()

		useClient(clientMock, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceProtectedOperations/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_apr_protected_operations.test", "protected_operation", `{"testKey":"testValue"}`),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResourceProtectedOperations/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_apr_protected_operations.test", "protected_operation", `{"testKey":"testValueUpdated"}`),
						),
					},
				},
			})
		})
	})
}
