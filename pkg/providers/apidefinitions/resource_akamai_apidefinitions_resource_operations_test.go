package apidefinitions

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAPIResourceOperations(t *testing.T) {
	t.Parallel()

	checker := test.NewStateChecker("akamai_apidefinitions_resource_operations.e2")

	var tests = map[string]struct {
		init   func(*v0.Mock)
		steps  []resource.TestStep
		error  *regexp.Regexp
		checks resource.TestCheckFunc
	}{
		"create api resource operations": {
			init: func(mV0 *v0.Mock) {
				mockUpdateResourceOperation(mV0, "resource-operations-01.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-01.json", 1)
				mockDeleteResourceOperation(mV0, 1)
			},
			steps: []resource.TestStep{
				{
					Config: apiResourceOperationsConfig(),
				},
			},
			checks: checker.CheckEqual("api_id", "1").CheckEqual("resource_operations", readJSONFile("resource-operations-01.json")).Build(),
		},
		"create api resource operations with all fields": {
			init: func(mV0 *v0.Mock) {
				mockUpdateResourceOperation(mV0, "resource-operations-02.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-02.json", 1)
				mockDeleteResourceOperation(mV0, 1)
			},
			steps: []resource.TestStep{
				{
					Config:           apiResourceOperationsCfgWithAllFieldsFromFile(),
					ImportState:      true,
					ImportStateId:    "1:1",
					ResourceName:     "akamai_apidefinitions_resource_operations.e2",
					ImportStateCheck: test.NewImportChecker().CheckEqual("api_id", "1").CheckEqual("resource_operations", "{\n  \"operations\": {\n    \"/index.php*\": {\n      \"onlineshop\": {\n        \"method\": \"POST\",\n        \"purpose\": \"LOGIN\",\n        \"parameters\": {\n          \"username\": {\n            \"path\": [\n              \"root\",\n              \"email\"\n            ],\n            \"location\": \"REQUEST_BODY\"\n          }\n        },\n        \"successConditions\": [\n          {\n            \"headerName\": \"X-Success\",\n            \"positiveMatch\": true,\n            \"suppressFromClientResponse\": false,\n            \"type\": \"HEADER_VALUE\",\n            \"valueCase\": false,\n            \"valueWildcard\": false,\n            \"values\": [\n              \"201\"\n            ]\n          }\n        ]\n      },\n      \"onlineshop-get\": {\n        \"method\": \"GET\",\n        \"purpose\": \"SEARCH\",\n        \"successConditions\": [\n          {\n            \"headerName\": \"X-Success\",\n            \"positiveMatch\": true,\n            \"suppressFromClientResponse\": false,\n            \"type\": \"HEADER_VALUE\",\n            \"valueCase\": false,\n            \"valueWildcard\": false,\n            \"values\": [\n              \"201\"\n            ]\n          }\n        ]\n      }\n    },\n    \"/login\": {\n      \"purposeLoginGET\": {\n        \"method\": \"GET\",\n        \"purpose\": \"ACCOUNT_VERIFICATION\"\n      },\n      \"purposeLoginPOST\": {\n        \"method\": \"POST\",\n        \"purpose\": \"ACCOUNT_VERIFICATION\"\n      }\n    }\n  }\n}").Build()},
			},
			checks: checker.CheckEqual("api_id", "1").CheckEqual("resource_operations", readJSONFile("resource-operations-02.json")).Build(),
		},
		"delete api resource operations": {
			init: func(mV0 *v0.Mock) {
				mockUpdateResourceOperation(mV0, "resource-operations-delete.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-delete.json", 1)
				mockDeleteResourceOperation(mV0, 1)
			},
			steps: []resource.TestStep{
				{
					Config: deleteAPIResourceOperationsConfig(),
				},
			},
			checks: checker.CheckEqual("api_id", "1").CheckEqual("resource_operations", readJSONFile("resource-operations-delete.json")).Build(),
		},
		"update api resource operations with all fields": {
			init: func(mV0 *v0.Mock) {
				mockUpdateResourceOperation(mV0, "resource-operations-02.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-02.json", 2)
				mockUpdateResourceOperation(mV0, "resource-operations-03.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-03.json", 1)
				mockDeleteResourceOperation(mV0, 1)
			},
			steps: []resource.TestStep{
				{
					Config:           apiResourceOperationsCfgWithAllFieldsFromFile(),
					ImportState:      true,
					ImportStateId:    "1:1",
					ResourceName:     "akamai_apidefinitions_resource_operations.e2",
					ImportStateCheck: test.NewImportChecker().CheckEqual("api_id", "1").CheckEqual("resource_operations", "{\n  \"operations\": {\n    \"/index.php*\": {\n      \"onlineshop\": {\n        \"method\": \"POST\",\n        \"purpose\": \"LOGIN\",\n        \"parameters\": {\n          \"username\": {\n            \"path\": [\n              \"root\",\n              \"email\"\n            ],\n            \"location\": \"REQUEST_BODY\"\n          }\n        },\n        \"successConditions\": [\n          {\n            \"headerName\": \"X-Success\",\n            \"positiveMatch\": true,\n            \"suppressFromClientResponse\": false,\n            \"type\": \"HEADER_VALUE\",\n            \"valueCase\": false,\n            \"valueWildcard\": false,\n            \"values\": [\n              \"201\"\n            ]\n          }\n        ]\n      },\n      \"onlineshop-get\": {\n        \"method\": \"GET\",\n        \"purpose\": \"SEARCH\",\n        \"successConditions\": [\n          {\n            \"headerName\": \"X-Success\",\n            \"positiveMatch\": true,\n            \"suppressFromClientResponse\": false,\n            \"type\": \"HEADER_VALUE\",\n            \"valueCase\": false,\n            \"valueWildcard\": false,\n            \"values\": [\n              \"201\"\n            ]\n          }\n        ]\n      }\n    },\n    \"/login\": {\n      \"purposeLoginGET\": {\n        \"method\": \"GET\",\n        \"purpose\": \"ACCOUNT_VERIFICATION\"\n      },\n      \"purposeLoginPOST\": {\n        \"method\": \"POST\",\n        \"purpose\": \"ACCOUNT_VERIFICATION\"\n      }\n    }\n  }\n}").Build(),
				},
				{
					Config: updateAPIiResourceOperationsCfgWithAllFields(),
				},
			},
			checks: checker.CheckEqual("api_id", "1").CheckEqual("resource_operations", readJSONFile("resource-operations-03.json")).Build(),
		},
		"update api resource operations with all fields : 400 Bad Request": {
			init: func(mV0 *v0.Mock) {
				mockUpdateResourceOperation(mV0, "resource-operations-02.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-02.json", 2)
				mockUpdateResourceOperationFail(mV0, 1)
				mockGetResourceOperation(mV0, "resource-operations-03.json", 1)
				mockDeleteResourceOperation(mV0, 1)
			},
			steps: []resource.TestStep{
				{
					Config: apiResourceOperationsCfgWithAllFieldsFromFile(),
				},
				{
					Config:      updateAPIiResourceOperationsCfgWithAllFields(),
					ExpectError: regexp.MustCompile("Upsert Resource Operations Failed"),
				},
			},
			checks: checker.CheckEqual("api_id", "1").CheckEqual("resource_operations", readJSONFile("resource-operations-02.json")).Build(),
		},
		"import state resource operations ok": {
			init: func(mV0 *v0.Mock) {
				mockUpdateResourceOperation(mV0, "resource-operations-01.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-01.json", 2)
				mockDeleteResourceOperation(mV0, 1)
			},
			steps: []resource.TestStep{
				{
					Config:             apiResourceOperationsConfig(),
					ImportState:        true,
					ImportStateId:      "1:1",
					ResourceName:       "akamai_apidefinitions_resource_operations.e1",
					ImportStateCheck:   test.NewImportChecker().CheckEqual("api_id", "1").CheckEqual("version", "1").CheckEqual("resource_operations", "{\n  \"operations\": {\n    \"/index.php*\": {\n      \"testPurpose\": {\n        \"method\": \"POST\",\n        \"purpose\": \"LOGIN\"\n      }\n    }\n  }\n}").Build(),
					ImportStatePersist: true,
				},
			},
			checks: checker.CheckEqual("api_id", "1").CheckEqual("resource_operations", readJSONFile("resource-operations-01.json")).Build(),
		},
		"import - invalid id format": {
			steps: []resource.TestStep{
				{
					Config:             apiResourceOperationsConfig(),
					ImportState:        true,
					ImportStateId:      "12345",
					ResourceName:       "akamai_apidefinitions_resource_operations.e1",
					ImportStatePersist: true,
					ExpectError:        regexp.MustCompile("Error: ID '12345' incorrectly formatted: should be 'API_ID:VERSION'"),
				},
			},
		},
		"import - invalid id value": {
			steps: []resource.TestStep{
				{
					Config:             apiResourceOperationsConfig(),
					ImportState:        true,
					ImportStateId:      "abc:123",
					ResourceName:       "akamai_apidefinitions_resource_operations.e1",
					ImportStatePersist: true,
					ExpectError:        regexp.MustCompile("Error: invalid API ID 'abc'"),
				},
			},
		},
		"import - invalid version value": {
			steps: []resource.TestStep{
				{
					Config:             apiResourceOperationsConfig(),
					ImportState:        true,
					ImportStateId:      "1:abc",
					ResourceName:       "akamai_apidefinitions_resource_operations.e1",
					ImportStatePersist: true,
					ExpectError:        regexp.MustCompile("Error: invalid API version 'abc'"),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := &apidefinitions.Mock{}
			clientV0 := &v0.Mock{}
			if test.init != nil {
				test.init(clientV0)
			}
			useClient(client, clientV0, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func mockUpdateResourceOperation(clientV0 *v0.Mock, file string, times int) {
	data, err := os.ReadFile("testdata/resourceOperations/" + file)
	if err != nil {
		log.Printf("Warning: Could not read file %s: %v", file, err)
		return
	}
	response, err := deserializeAPIRequest(string(data))
	if err != nil {
		log.Printf("Warning: Deserialization error  %s: %v", file, err)
		return
	}
	clientV0.On("UpdateResourceOperation", mock.Anything, mock.Anything).
		Return(ptr.To(v0.UpdateResourceOperationResponse(*response)), nil).Times(times)
}

func mockUpdateResourceOperationFail(clientV0 *v0.Mock, times int) {
	clientV0.On("UpdateResourceOperation", mock.Anything, mock.Anything).
		Return(nil, &badRequestErrorResOperations).Times(times)
}

func mockGetResourceOperation(clientV0 *v0.Mock, file string, times int) {
	data, err := os.ReadFile("testdata/resourceOperations/" + file)
	if err != nil {
		log.Printf("Warning: Could not read file %s: %v", file, err)
		return
	}

	response := v0.GetResourceOperationResponse{}

	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		return
	}
	clientV0.On("GetResourceOperation", mock.Anything, mock.Anything).
		Return(ptr.To(response), nil).Times(times)
}

func mockDeleteResourceOperation(clientV0 *v0.Mock, times int) {
	response := v0.DeleteResourceOperationResponse{
		APIID:         1,
		VersionNumber: 2,
		Status:        200,
		Detail:        "Api resource operations for Endpoint is Deleted",
	}
	clientV0.On("DeleteResourceOperation", mock.Anything, mock.Anything).
		Return(ptr.To(v0.DeleteResourceOperationResponse(response)), nil).Times(times)
}

func readJSONFile(file string) string {
	data, err := os.ReadFile("testdata/resourceOperations/" + file)
	if err != nil {
		log.Printf("Warning: Could not read file %s: %v", file, err)
		return ""
	}
	return string(data)
}

func deserializeAPIRequest(body string) (*v0.UpdateResourceOperationResponse, error) {
	endpoint := v0.UpdateResourceOperationResponse{}

	err := json.Unmarshal([]byte(body), &endpoint)
	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}

func apiResourceOperationsConfig() string {
	return providerConfig + `
resource "akamai_apidefinitions_resource_operations" "e1" {
  api_id = 1
  version = 1
  resource_operations = jsonencode({
  "operations": {
    "/index.php*": {
      "testPurpose": {
        "method": "POST",
        "purpose": "LOGIN"
      }
    }
  }
})
}
`
}

func deleteAPIResourceOperationsConfig() string {
	return providerConfig + `
		resource "akamai_apidefinitions_resource_operations" "e3" {
		  api_id = 1
          version = 1
		  resource_operations = file("testdata/resourceOperations/resource-operations-delete.json")
		}`
}

func apiResourceOperationsCfgWithAllFieldsFromFile() string {
	return providerConfig + `
			resource "akamai_apidefinitions_resource_operations" "e2" {
			  api_id = 1
              version = 1
			  resource_operations = file("testdata/resourceOperations/resource-operations-02.json")
			}`
}

func updateAPIiResourceOperationsCfgWithAllFields() string {
	return providerConfig + `
			resource "akamai_apidefinitions_resource_operations" "e2" {
			  api_id = 1
              version = 1
			  resource_operations = file("testdata/resourceOperations/resource-operations-03.json")
			}`
}

var badRequestErrorResOperations = v0.Error{
	Type:     "/api-definitions/error-types/invalid-input-error",
	Title:    "Invalid input error",
	Detail:   "The request you submitted is invalid. Modify the request and try again.",
	Instance: "id_001",
	Status:   400,
	Severity: ptr.To("ERROR"),
	Errors: []v0.Error{
		{
			Type:     "/api-definitions/error-types/resource-path-operation-check",
			Title:    "resource-path-operation-check.title",
			Detail:   "resource-path-operation-check.detail",
			Severity: ptr.To("ERROR"),
			Field:    ptr.To("put.dto.resourceOperationsMap[/base].<map value>[test login].operationParameter"),
			RejectedValue: map[string]interface{}{
				"method":           "POST",
				"operationPurpose": "LOGIN",
			},
		},
		{
			Type:     "/api-definitions/error-types/resource-path-operation-check",
			Title:    "resource-path-operation-check.title",
			Detail:   "resource-path-operation-check.detail",
			Severity: ptr.To("ERROR"),
			Field:    ptr.To("put.dto.resourceOperationsMap[/base].<map value>[test login].operationParameter.username"),
			RejectedValue: map[string]interface{}{
				"method":           "POST",
				"operationPurpose": "LOGIN",
			},
		},
	},
}
