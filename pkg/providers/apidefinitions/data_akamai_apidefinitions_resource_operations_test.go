package apidefinitions

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceOperationsDataSource(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		APIID int64
		init  func(*apidefinitions.Mock)
		steps []resource.TestStep
		error *regexp.Regexp
	}{
		"200 - OK Search using no resource path and name": {
			init: func(m *apidefinitions.Mock) {
				mockSearchResourceOperations(m, 3)
			},
			steps: []resource.TestStep{
				{
					Config: resourceOperationsDataSourceConfig(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_resource_operations.test", "resource_operations", readJSONFileWithoutIndentation("resource-operations-data-200.json")),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_resource_operations.test", "version", "4"),
					),
				},
			},
		},
		"200 - OK Search using resource path and name": {
			init: func(m *apidefinitions.Mock) {
				mockSearchResourceOperations(m, 3)
			},
			steps: []resource.TestStep{
				{
					Config: resourceOperationsSearchDataSourceConfig(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_resource_operations.test", "resource_operations", readJSONFileWithoutIndentation("resource-operations-data-200.json")),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_resource_operations.test", "version", "4"),
					),
				},
			},
		},
		"200 - OK Search using only resource path": {
			init: func(m *apidefinitions.Mock) {
				mockSearchResourceOperations(m, 3)
			},
			steps: []resource.TestStep{
				{
					Config: resourceOperationsWithOnlyResPathDataSourceConfig(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_resource_operations.test", "resource_operations", readJSONFileWithoutIndentation("resource-operations-data-200.json")),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_resource_operations.test", "version", "4"),
					),
				},
			},
		},
		"500 - error": {
			init: func(m *apidefinitions.Mock) {
				mockFailedSearchResourceOperations(m)
			},
			steps: []resource.TestStep{
				{
					Config:      resourceOperationsDataSourceConfig(),
					ExpectError: regexp.MustCompile("Error retrieving resource operations"),
				},
			},
		},
		"check schema - missing required attributes": {
			steps: []resource.TestStep{
				{
					Config:      failedResourceOperationsDataSourceConfig(),
					ExpectError: regexp.MustCompile("The argument \"api_id\" is required|The argument \"version\" is required"),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &apidefinitions.Mock{}
			if test.init != nil {
				test.init(client)
			}
			useClient(client, nil, func() {
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

func resourceOperationsDataSourceConfig() string {
	return providerConfig + `
data "akamai_apidefinitions_resource_operations" "test" {
	api_id = 503130
}
`
}

func resourceOperationsSearchDataSourceConfig() string {
	return providerConfig + `
data "akamai_apidefinitions_resource_operations" "test" {
	api_id = 503130
	resource_path = "/index.php"
	resource_name = "Login Form Data"
}
`
}

func resourceOperationsWithOnlyResPathDataSourceConfig() string {
	return providerConfig + `
data "akamai_apidefinitions_resource_operations" "test" {
	api_id = 503130
	resource_path = "/index.php"
}
`
}

func failedResourceOperationsDataSourceConfig() string {
	return providerConfig + `
data "akamai_apidefinitions_resource_operations" "test" {
	// Intentionally omitting required attributes to test validation
}
`
}

func mockSearchResourceOperations(client *apidefinitions.Mock, times int) {
	file := "resource-operations-data.json"
	data, err := os.ReadFile("testdata/resourceOperations/" + file)
	if err != nil {
		log.Printf("Warning: Could not read file %s: %v", file, err)
		return
	}

	res, err := deserializeSearchResourceOperationsResponse(string(data))
	if err != nil {
		log.Printf("Warning: Deserialization error  %s: %v", file, err)
		return
	}

	client.On("SearchResourceOperations", mock.Anything, mock.Anything, mock.Anything).
		Return(res, nil).
		Times(times)
}

func mockFailedSearchResourceOperations(client *apidefinitions.Mock) {
	client.On("SearchResourceOperations", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, &serverError).
		Once()
}

func deserializeSearchResourceOperationsResponse(body string) (*apidefinitions.SearchResourceOperationsResponse, error) {
	deserializedResponse := apidefinitions.SearchResourceOperationsResponse{}

	err := json.Unmarshal([]byte(body), &deserializedResponse)
	if err != nil {
		return nil, err
	}

	return &deserializedResponse, nil
}

func readJSONFileWithoutIndentation(file string) string {
	data, err := os.ReadFile("testdata/resourceOperations/" + file)
	if err != nil {
		log.Printf("Warning: Could not read file %s: %v", file, err)
		return ""
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, data); err != nil {
		log.Printf("Warning: Could not compact JSON %s: %v", file, err)
		return string(data)
	}

	return buf.String()
}

var serverError = apidefinitions.Error{
	Status: 500,
	Detail: "Internal server error",
}
