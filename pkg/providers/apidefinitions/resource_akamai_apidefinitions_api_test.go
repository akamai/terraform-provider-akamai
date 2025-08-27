package apidefinitions

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type data struct {
	response string
}

func TestAPIResource(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		configPath string
		init       func(*apidefinitions.Mock, *v0.Mock)
		mockData   data
		steps      []resource.TestStep
		error      *regexp.Regexp
	}{
		"create endpoint - required only": {
			init: func(m *apidefinitions.Mock, mV0 *v0.Mock) {
				mockRegisterAPI(mV0, "required-only-response.json")
				mockReadAPIResource(m, mV0, "required-only-response.json", 1)
				mockDeleteEndpoint(m)
			},
			mockData: data{
				response: "required-only-response.json",
			},
			steps: []resource.TestStep{
				{
					Config: endpointResourceConfig("api-configuration-required-only.json"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "api", readJSON("api-configuration-required-only.json")),
					),
				},
			},
		},
		"create endpoint - with resources": {
			init: func(m *apidefinitions.Mock, mV0 *v0.Mock) {
				mockRegisterAPI(mV0, "with-resources-response.json")
				mockReadActiveAPIResource(m, mV0, "with-resources-response.json", 1)
				mockDeleteEndpoint(m)
			},
			mockData: data{
				response: "with-resources-response.json",
			},
			steps: []resource.TestStep{
				{
					Config: endpointResourceConfig("api-configuration-with-resources.json"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "api", readJSON("api-configuration-with-resources.json")),
					),
				},
			},
		},
		"create endpoint - 400 Bad Request": {
			init: func(_ *apidefinitions.Mock, mV0 *v0.Mock) {
				mockRegisterAPIFail(mV0)
			},
			mockData: data{
				response: "required-only-response.json",
			},
			steps: []resource.TestStep{
				{
					Config:      endpointResourceConfig("api-configuration-required-only.json"),
					ExpectError: regexp.MustCompile("Create API Failed"),
				},
			},
		},
		"update endpoint": {
			init: func(m *apidefinitions.Mock, mV0 *v0.Mock) {
				mockRegisterAPI(mV0, "required-only-response.json")
				mockReadActiveAPIResource(m, mV0, "required-only-response.json", 2)
				mockUpdateAPIVersion(m, mV0, "with-resources-response.json", false)
				mockReadActiveAPIResource(m, mV0, "with-resources-response.json", 1)
				mockDeleteEndpoint(m)
			},
			steps: []resource.TestStep{
				{
					Config: endpointResourceConfig("api-configuration-required-only.json"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "api", readJSON("api-configuration-required-only.json")),
					),
				},
				{
					Config: endpointResourceConfig("api-configuration-with-resources.json"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "api", readJSON("api-configuration-with-resources.json")),
					),
				},
			},
		},
		"update endpoint - version is locked": {
			init: func(m *apidefinitions.Mock, mV0 *v0.Mock) {
				mockRegisterAPI(mV0, "required-only-response.json")
				mockReadAPIResource(m, mV0, "required-only-response.json", 2)
				mockUpdateAPIVersion(m, mV0, "with-resources-response.json", true)
				mockReadAPIResource(m, mV0, "with-resources-response.json", 1)
				mockDeleteEndpoint(m)
			},
			steps: []resource.TestStep{
				{
					Config: endpointResourceConfig("api-configuration-required-only.json"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "api", readJSON("api-configuration-required-only.json")),
					),
				},
				{
					Config: endpointResourceConfig("api-configuration-with-resources.json"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "latest_version", "2"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "api", readJSON("api-configuration-with-resources.json")),
					),
				},
			},
		},
		"update endpoint - 400 Bad Request": {
			init: func(m *apidefinitions.Mock, mV0 *v0.Mock) {
				mockRegisterAPI(mV0, "required-only-response.json")
				mockReadActiveAPIResource(m, mV0, "required-only-response.json", 2)
				mockUpdateAPIVersionFail(m, mV0)
				mockDeleteEndpoint(m)
			},
			steps: []resource.TestStep{
				{
					Config: endpointResourceConfig("api-configuration-required-only.json"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "id", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "latest_version", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "api", readJSON("api-configuration-required-only.json")),
					),
				},
				{
					Config:      endpointResourceConfig("api-configuration-with-resources.json"),
					ExpectError: regexp.MustCompile("Update API Failed"),
				},
			},
		},
		"refresh endpoint - update staging and production version": {
			init: func(m *apidefinitions.Mock, mV0 *v0.Mock) {
				mockRegisterAPI(mV0, "required-only-response.json")
				mockReadActiveAPIResource(m, mV0, "required-only-response.json", 3)
				mockDestroyActiveAPIResource(m, 2)
			},
			mockData: data{
				response: "required-only-response.json",
			},
			steps: []resource.TestStep{
				{
					Config: endpointResourceConfig("api-configuration-required-only.json"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "id", "1"),
					),
				},
				{
					RefreshState: true,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "staging_version", "1"),
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "production_version", "1"),
					),
				},
			},
		},
		"delete endpoint - previously activated": {
			init: func(m *apidefinitions.Mock, mV0 *v0.Mock) {
				mockRegisterAPI(mV0, "required-only-response.json")
				mockReadActiveAPIResource(m, mV0, "required-only-response.json", 1)
				mockDestroyActiveAPIResource(m, 2)
			},
			mockData: data{
				response: "required-only-response.json",
			},
			steps: []resource.TestStep{
				{
					Config: endpointResourceConfig("api-configuration-required-only.json"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_apidefinitions_api.e1", "id", "1"),
					),
				},
			},
		},
		"check schema - missing required attributes": {
			steps: []resource.TestStep{
				{
					Config:      endpointResourceConfigWithEmptyAPI(),
					ExpectError: regexp.MustCompile("Attribute api cannot be empty"),
				},
				{
					Config:      endpointResourceConfigWithEmptyContractID(),
					ExpectError: regexp.MustCompile("Attribute contract_id string length must be at least 1, got: 0"),
				},
				{
					Config:      endpointResourceConfigWithoutGroupID(),
					ExpectError: regexp.MustCompile("The argument \"group_id\" is required, but no definition was found"),
				},
			},
		},
		"check schema - unexpected fields": {
			steps: []resource.TestStep{
				{
					Config:      endpointResourceConfig("api-configuration-with-unexpected-fields.json"),
					ExpectError: regexp.MustCompile("Error: Invalid JSON provided"),
				},
			},
		},
		"import - ok": {
			init: func(m *apidefinitions.Mock, mV0 *v0.Mock) {
				mockImportState(m, mV0)
				mockReadActiveAPIResource(m, mV0, "required-only-response.json", 1)
				mockDestroyActiveAPIResource(m, 2)
			},
			steps: []resource.TestStep{
				{
					Config:        endpointResourceConfig("api-configuration-required-only.json"),
					ImportState:   true,
					ImportStateId: "1:1",
					ResourceName:  "akamai_apidefinitions_api.e1",
					ImportStateCheck: func(states []*terraform.InstanceState) error {
						state := states[0].Attributes
						assert.Equal(t, "1", state["id"])
						assert.Equal(t, "1", state["latest_version"])
						assert.Equal(t, "Contract-1", state["contract_id"])
						assert.Equal(t, "1", state["group_id"])
						return nil
					},
					ImportStatePersist: true,
				},
			},
		},
		"import - invalid id format": {
			steps: []resource.TestStep{
				{
					Config:             endpointResourceConfig("api-configuration-required-only.json"),
					ImportState:        true,
					ImportStateId:      "12345",
					ResourceName:       "akamai_apidefinitions_api.e1",
					ImportStatePersist: true,
					ExpectError:        regexp.MustCompile("Error: ID '12345' incorrectly formatted: should be 'API_ID:VERSION'"),
				},
			},
		},
		"import - invalid id value": {
			steps: []resource.TestStep{
				{
					Config:             endpointResourceConfig("api-configuration-required-only.json"),
					ImportState:        true,
					ImportStateId:      "abc:123",
					ResourceName:       "akamai_apidefinitions_api.e1",
					ImportStatePersist: true,
					ExpectError:        regexp.MustCompile("Error: invalid API id 'abc'"),
				},
			},
		},
		"import - invalid version value": {
			steps: []resource.TestStep{
				{
					Config:             endpointResourceConfig("api-configuration-required-only.json"),
					ImportState:        true,
					ImportStateId:      "12345:abc",
					ResourceName:       "akamai_apidefinitions_api.e1",
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
				test.init(client, clientV0)
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

func mockImportState(m *apidefinitions.Mock, mV0 *v0.Mock) {
	mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive), 1)
	mockGetAPIVersion(mV0, "required-only-response.json", 1)
}

func mockDestroyActiveAPIResource(m *apidefinitions.Mock, times int) {
	mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive), 1)
	mockDeactivateVersion(m, times)
	mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusDeactivated), 1)
	mockHideEndpoint(m)
}

func mockDeleteEndpoint(client *apidefinitions.Mock) *mock.Call {
	mockGetEndpointWithActivationStatus(client, nil, nil, 1)
	return client.On("DeleteEndpoint", mock.Anything, apidefinitions.DeleteEndpointRequest{APIEndpointID: 1}).
		Return(nil).
		Once()
}

func mockRegisterAPI(client *v0.Mock, file string) {
	client.On("RegisterAPI", mock.Anything, mock.Anything).
		Return(ptr.To(v0.RegisterAPIResponse(*toState(file))), nil).
		Once()
}

func mockRegisterAPIFail(client *v0.Mock) {
	client.On("RegisterAPI", mock.Anything, mock.Anything).
		Return(nil, &badRequestError).
		Once()
}

func mockUpdateAPIVersion(client *apidefinitions.Mock, clientV0 *v0.Mock, file string, locked bool) {
	mockGetEndpointVersion(client, locked)
	clientV0.On("UpdateAPIVersion", mock.Anything, mock.Anything).
		Return(ptr.To(v0.UpdateAPIVersionResponse(*toState(file))), nil).
		Once()
	if locked {
		mockCloneEndpointVersion(client)
	}
}

func mockUpdateAPIVersionFail(client *apidefinitions.Mock, clientV0 *v0.Mock) {
	mockGetEndpointVersion(client, false)
	clientV0.On("UpdateAPIVersion", mock.Anything, mock.Anything).
		Return(nil, &badRequestError)
}

func mockGetAPIVersion(clientV0 *v0.Mock, file string, times int) {
	data, _ := os.ReadFile("testdata/endpoint/" + file)
	response, _ := deserializeAPI(string(data))
	clientV0.On("GetAPIVersion", mock.Anything, v0.GetAPIVersionRequest{ID: 1, Version: 1}).
		Return(ptr.To(v0.GetAPIVersionResponse(*response)), nil).Times(times)

}

func mockReadAPIResource(client *apidefinitions.Mock, clientV0 *v0.Mock, file string, times int) {
	mockGetAPIVersion(clientV0, file, times)
	mockGetEndpointWithActivationStatus(client, nil, nil, times)
	mockListEndpointVersions(client, times)
}

func mockReadActiveAPIResource(client *apidefinitions.Mock, clientV0 *v0.Mock, file string, times int) {
	mockGetAPIVersion(clientV0, file, times)
	mockGetEndpointWithActivationStatus(client, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive), times)
	mockListEndpointVersions(client, times)
}

func mockGetEndpointVersion(client *apidefinitions.Mock, locked bool) {
	client.On("GetEndpointVersion", mock.Anything, apidefinitions.GetEndpointVersionRequest{APIEndpointID: 1, VersionNumber: 1}).
		Return(&apidefinitions.GetEndpointVersionResponse{
			Locked: locked,
		}, nil).
		Once()
}

func mockCloneEndpointVersion(client *apidefinitions.Mock) {
	client.On("CloneEndpointVersion", mock.Anything, apidefinitions.CloneEndpointVersionRequest{APIEndpointID: 1, VersionNumber: 1}).
		Return(&apidefinitions.CloneEndpointVersionResponse{
			VersionNumber: int64(2),
		}, nil).
		Once()
}

func mockHideEndpoint(client *apidefinitions.Mock) *mock.Call {
	return client.On("HideEndpoint", mock.Anything, apidefinitions.HideEndpointRequest{APIEndpointID: 1}).
		Return(&apidefinitions.HideEndpointResponse{}, nil).
		Once()
}

func mockListEndpointVersions(client *apidefinitions.Mock, times int) *mock.Call {
	return client.On("ListEndpointVersions", mock.Anything, mock.Anything).
		Return(&apidefinitions.ListEndpointVersionsResponse{
			TotalSize: 2,
			APIVersions: []apidefinitions.APIVersion{
				{
					VersionNumber: 1,
				},
			},
		}, nil).
		Times(times)
}

func readJSON(file string) string {
	data, _ := os.ReadFile("testdata/endpoint/" + file)
	return string(data)
}

func toState(file string) *v0.API {
	data, _ := os.ReadFile("testdata/endpoint/" + file)
	response, _ := deserializeAPI(string(data))
	return response
}

func toStateJSON(file string) string {
	data := toState(file)
	json, _ := serializeIndent(data.RegisterAPIRequest.APIAttributes)
	return *json
}

func endpointResourceConfig(file string) string {
	return providerConfig + fmt.Sprintf(`
resource "akamai_apidefinitions_api" "e1" {
  api = file("testdata/endpoint/%v")
  contract_id = "Contract-1"
  group_id = 12345
}
`, file)
}

func endpointResourceConfigWithEmptyAPI() string {
	return providerConfig + `
resource "akamai_apidefinitions_api" "e1" {
  api = ""
  contract_id = "Contract-1"
  group_id = 12345
}`
}

func endpointResourceConfigWithEmptyContractID() string {
	return providerConfig + `
resource "akamai_apidefinitions_api" "e1" {
  api = "{}"
  contract_id = ""
  group_id = 12345
}`
}

func endpointResourceConfigWithoutGroupID() string {
	return providerConfig + `
resource "akamai_apidefinitions_api" "e1" {
  api = "{}"
  contract_id = "Contract-1"
}`
}

var badRequestError = v0.Error{
	Status: 400,
	Detail: "Bad Request",
}

func deserializeAPI(body string) (*v0.API, error) {
	endpoint := v0.API{}

	err := json.Unmarshal([]byte(body), &endpoint)
	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}
