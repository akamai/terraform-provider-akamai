package apidefinitions

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type testData struct {
	response string
}

func TestAPIResourceOperations(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		configPath string
		init       func(*testing.T, *apidefinitions.Mock, *v0.Mock, testData)
		mockData   testData
		steps      []resource.TestStep
		error      *regexp.Regexp
	}{
		"create api resource operations": {
			init: func(t *testing.T, m *apidefinitions.Mock, mV0 *v0.Mock, resourceData testData) {
				mockListEndpointVersions(m, 3)
				mockUpdateResourceOperation(mV0, "resource-operations-01.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-01.json", 1)
				mockDeleteResourceOperation(mV0, 1)
			},
			mockData: testData{
				response: "resource-operations-01.json",
			},
			steps: []resource.TestStep{
				{
					Config: apiResourceOperationsConfig(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e1", "endpoint_id", "1"),
					),
				},
			}},
		"create api resource operations with all fields": {
			init: func(t *testing.T, m *apidefinitions.Mock, mV0 *v0.Mock, resourceData testData) {
				mockListEndpointVersions(m, 3)
				mockUpdateResourceOperation(mV0, "resource-operations-02.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-02.json", 1)
				mockDeleteResourceOperation(mV0, 1)
			},
			mockData: testData{
				response: "resource-operations-02.json",
			},
			steps: []resource.TestStep{
				{
					Config: apiResourceOperationsCfgWithAllFields(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e2", "endpoint_id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e2", "resource_operations", readJSONFile("resource-operations-02.json")),
					),
				},
			}},
		"delete api resource operations": {
			init: func(t *testing.T, m *apidefinitions.Mock, mV0 *v0.Mock, resourceData testData) {
				mockListEndpointVersions(m, 3)
				mockUpdateResourceOperation(mV0, "resource-operations-delete.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-delete.json", 1)
				mockDeleteResourceOperation(mV0, 1)
			},
			mockData: testData{
				response: "resource-operations-delete.json",
			},
			steps: []resource.TestStep{
				{
					Config: deleteAPIResourceOperationsConfig(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e3", "endpoint_id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e3", "resource_operations", readJSONFile("resource-operations-delete.json")),
					),
				},
			}},
		"update api resource operations with all fields": {
			init: func(t *testing.T, m *apidefinitions.Mock, mV0 *v0.Mock, resourceData testData) {
				mockListEndpointVersions(m, 6)
				mockUpdateResourceOperation(mV0, "resource-operations-02.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-02.json", 2)
				mockUpdateResourceOperation(mV0, "resource-operations-03.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-03.json", 1)
				mockDeleteResourceOperation(mV0, 1)
			},
			mockData: testData{
				response: "resource-operations-02.json",
			},
			steps: []resource.TestStep{
				{
					Config: apiResourceOperationsCfgWithAllFields(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e2", "endpoint_id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e2", "resource_operations", readJSONFile("resource-operations-02.json")),
					),
				},
				{
					Config: updateAPIiResourceOperationsCfgWithAllFields(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e2", "endpoint_id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e2", "resource_operations", readJSONFile("resource-operations-03.json")),
					),
				},
			}},
		"update api resource operations with all fields : 400 Bad Request": {
			init: func(t *testing.T, m *apidefinitions.Mock, mV0 *v0.Mock, resourceData testData) {
				mockListEndpointVersions(m, 5)
				mockUpdateResourceOperation(mV0, "resource-operations-02.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-02.json", 2)
				mockUpdateResourceOperationFail(mV0)
				mockGetResourceOperation(mV0, "resource-operations-03.json", 1)
				mockDeleteResourceOperation(mV0, 1)
			},
			mockData: testData{
				response: "resource-operations-02.json",
			},
			steps: []resource.TestStep{
				{
					Config: apiResourceOperationsCfgWithAllFields(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e2", "endpoint_id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_resource_operations.e2", "resource_operations", readJSONFile("resource-operations-02.json")),
					),
				},
				{
					Config:      updateAPIiResourceOperationsCfgWithAllFields(),
					ExpectError: regexp.MustCompile("Create Resource Operations Failed"),
				},
			}},
		"import state resource operations ok": {
			init: func(t *testing.T, m *apidefinitions.Mock, mV0 *v0.Mock, resourceData testData) {
				mockListEndpointVersions(m, 2)
				mockUpdateResourceOperation(mV0, "resource-operations-01.json", 1)
				mockGetResourceOperation(mV0, "resource-operations-01.json", 2)
				mockDeleteResourceOperation(mV0, 1)
			},
			mockData: testData{
				response: "resource-operations-01.json",
			},
			steps: []resource.TestStep{
				{
					Config:        apiResourceOperationsConfig(),
					ImportState:   true,
					ImportStateId: "1:1",
					ResourceName:  "akamai_apidefinitions_resource_operations.e1",
					ImportStateCheck: func(states []*terraform.InstanceState) error {
						state := states[0].Attributes
						assert.Equal(t, "1", state["endpoint_id"])
						return nil
					},
					ImportStatePersist: true,
				},
			}},
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
					ExpectError:        regexp.MustCompile("Error: invalid API id 'abc'"),
				},
			},
		},
		"import - invalid version value": {
			steps: []resource.TestStep{
				{
					Config:             apiResourceOperationsConfig(),
					ImportState:        true,
					ImportStateId:      "12345:abc",
					ResourceName:       "akamai_apidefinitions_resource_operations.e1",
					ImportStatePersist: true,
					ExpectError:        regexp.MustCompile("Error: invalid API version 'abc'"),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &apidefinitions.Mock{}
			clientV0 := &v0.Mock{}
			if test.init != nil {
				test.init(t, client, clientV0, test.mockData)
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
	data, _ := os.ReadFile("testdata/resourceOperations/" + file)
	response, _ := deserializeAPIRequest(string(data))
	clientV0.On("UpdateResourceOperation", mock.Anything, mock.Anything).
		Return(ptr.To(v0.UpdateResourceOperationResponse(*response)), nil).Times(times)
}

func mockUpdateResourceOperationFail(clientV0 *v0.Mock) {
	clientV0.On("UpdateResourceOperation", mock.Anything, mock.Anything).
		Return(nil, &badRequestError)
}
func mockGetResourceOperation(clientV0 *v0.Mock, file string, times int) {
	data, _ := os.ReadFile("testdata/resourceOperations/" + file)
	response := v0.GetResourceOperationResponse{}

	err := json.Unmarshal([]byte(data), &response)
	if err != nil {
		return
	}
	clientV0.On("GetResourceOperation", mock.Anything, mock.Anything).
		Return(ptr.To(response), nil).Times(times)
}

func mockDeleteResourceOperation(clientV0 *v0.Mock, times int) {
	response := v0.DeleteResourceOperationResponse{
		APIEndpointID: 1,
		VersionNumber: 2,
		Status:        200,
		Detail:        "Api resource operations for Endpoint is Deleted",
	}
	clientV0.On("DeleteResourceOperation", mock.Anything, mock.Anything).
		Return(ptr.To(v0.DeleteResourceOperationResponse(response)), nil).Times(times)
}

func readJSONFile(file string) string {
	data, _ := os.ReadFile("testdata/resourceOperations/" + file)

	response := v0.GetResourceOperationResponse{}
	// unmarshal the input json file to struct
	err := json.Unmarshal([]byte(data), &response)

	if err != nil {
		return ""
	}

	// marshal the struct to json string
	jsonString, err := json.Marshal(response)

	if err != nil {
		return ""
	}

	// normalize the json string response
	jsonFiole, err := normalizeJSON(string(jsonString))

	if err != nil {
		return ""
	}

	if err != nil {
		return ""
	}

	return jsonFiole
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
	return providerConfig + fmt.Sprintf(`
resource "akamai_apidefinitions_resource_operations" "e1" {
  endpoint_id = 1
  resource_operations = jsonencode({
  "resourceOperationsMap": {
    "/index.php*": {
      "testPurpose": {
        "resourcePath": "/resources/123",
        "method": "POST",
        "operationName": "testpurpose",
        "operationPurpose": "LOGIN"
      }
    }
  }
})
}
`)
}

func deleteAPIResourceOperationsConfig() string {
	return providerConfig + fmt.Sprintf(`
resource "akamai_apidefinitions_resource_operations" "e3" {
  endpoint_id = 1
  resource_operations = jsonencode({
  "resourceOperationsMap": {}
 })
}
`)

}

func apiResourceOperationsCfgWithAllFields() string {
	return providerConfig + fmt.Sprintf(`
resource "akamai_apidefinitions_resource_operations" "e2" {
  endpoint_id = 1
  resource_operations = jsonencode({
  "resourceOperationsMap": {
    "/index.php*": {
      "onlineshop": {
        "resourcePath": "/index.php*",
        "operationName": "onlineshop",
        "operationPurpose": "LOGIN",
        "method": "POST",
        "operationParameter": {
          "username": {
            "parameterPath": [
              "root",
              "email"
            ],
            "parameterLocation": "REQUEST_BODY"
          }
        },
        "successConditions": [
          {
            "headerName": "X-Success",
            "positiveMatch": true,
            "suppressFromClientResponse": false,
            "type": "HEADER_VALUE",
            "valueCase": false,
            "valueWildcard": false,
            "values": [
              "201"
            ]
          }
        ]
      },
      "onlineshop-get": {
        "resourcePath": "/index.php*",
        "operationName": "onlineshop-get",
        "operationPurpose": "SEARCH",
        "method": "GET",
        "successConditions": [
          {
            "headerName": "X-Success",
            "positiveMatch": true,
            "suppressFromClientResponse": false,
            "type": "HEADER_VALUE",
            "valueCase": false,
            "valueWildcard": false,
            "values": [
              "201"
            ]
          }
        ]
      }
    },
    "/login": {
      "purposeLoginGET": {
        "resourcePath": "/login",
        "operationName": "purposeLoginGET",
        "operationPurpose": "ACCOUNT_VERIFICATION",
        "method": "GET"
      },
      "purposeLoginPOST": {
        "resourcePath": "/login",
        "operationName": "purposeLoginPOST",
        "operationPurpose": "ACCOUNT_VERIFICATION",
        "method": "POST"
      }
    }
  }
})
}
`)
}

func updateAPIiResourceOperationsCfgWithAllFields() string {
	return providerConfig + fmt.Sprintf(`
resource "akamai_apidefinitions_resource_operations" "e2" {
  endpoint_id = 1
  resource_operations = jsonencode({
  "resourceOperationsMap": {
    "/index.php*": {
      "onlineshop": {
        "resourcePath": "/index.php*",
        "operationName": "onlineshop",
        "operationPurpose": "ACCOUNT_CREATION",
        "method": "POST",
        "operationParameter": {
          "username": {
            "parameterPath": [
              "root123",
              "email"
            ],
            "parameterLocation": "REQUEST_BODY"
          }
        },
        "successConditions": [
          {
            "headerName": "X-Success",
            "positiveMatch": true,
            "suppressFromClientResponse": false,
            "type": "HEADER_VALUE",
            "valueCase": false,
            "valueWildcard": false,
            "values": [
              "201"
            ]
          }
        ]
      },
      "onlineshop-get": {
        "resourcePath": "/index.php*",
        "operationName": "onlineshop-get",
        "operationPurpose": "SEARCH",
        "method": "GET",
        "successConditions": [
          {
            "headerName": "X-Success",
            "positiveMatch": true,
            "suppressFromClientResponse": false,
            "type": "HEADER_VALUE",
            "valueCase": false,
            "valueWildcard": false,
            "values": [
              "201"
            ]
          }
        ]
      }
    },
    "/login": {
      "purposeLoginGET": {
        "resourcePath": "/login",
        "operationName": "purposeLoginGET",
        "operationPurpose": "ACCOUNT_VERIFICATION",
        "method": "GET"
      },
      "purposeLoginPOST": {
        "resourcePath": "/login",
        "operationName": "purposeLoginPOST",
        "operationPurpose": "ACCOUNT_VERIFICATION",
        "method": "POST"
      }
    }
  }
})
}
`)
}
