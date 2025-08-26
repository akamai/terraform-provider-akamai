package apidefinitions

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions"
	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestApiDefinitionConfigDataSource(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		APIID   int64
		Version int64
		init    func(*apidefinitions.Mock, *v0.Mock)
		steps   []resource.TestStep
		error   *regexp.Regexp
	}{
		"200 - By ID - API Exists": {
			init: func(m *apidefinitions.Mock, v0 *v0.Mock) {
				mockListEndpointVersions(m, 3)
				mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive), 3)
				mockGetAPIVersion(v0, "with-resources-response.json", 3)
			},
			steps: []resource.TestStep{
				{
					Config: configWithID(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "id", "1"),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "api", toStateJSON("with-resources-response.json")),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "contract_id", "1"),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "group_id", "1"),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "staging_version", "1"),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "production_version", "1"),
					),
				},
			},
		},
		"403 - By ID - API Don't Exists": {
			init: func(m *apidefinitions.Mock, _ *v0.Mock) {
				mockListEndpointVersionsError(m, 403)
			},
			steps: []resource.TestStep{
				{
					Config:      configWithID(),
					ExpectError: regexp.MustCompile("unable to find API with ID 1"),
				},
			},
		},
		"200 - By Name - API Exists": {
			init: func(m *apidefinitions.Mock, v0 *v0.Mock) {
				mockListEndpoints(m, 1, "Pet Store", 3)
				mockListEndpointVersions(m, 3)
				mockGetEndpointWithActivationStatus(m, ptr.To(int64(1)), ptr.To(apidefinitions.ActivationStatusActive), 3)
				mockGetAPIVersion(v0, "with-resources-response.json", 3)
			},
			steps: []resource.TestStep{
				{
					Config: configWithName(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "id", "1"),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "api", toStateJSON("with-resources-response.json")),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "contract_id", "1"),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "group_id", "1"),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "latest_version", "1"),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "staging_version", "1"),
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_api.test", "production_version", "1"),
					),
				},
			},
		},
		"200 - By Name - API Don't Exists": {
			init: func(m *apidefinitions.Mock, _ *v0.Mock) {
				mockListEndpoints(m, 1, "API Don't Exists", 1)
			},
			steps: []resource.TestStep{
				{
					Config:      configWithName(),
					ExpectError: regexp.MustCompile("unable to find API with Name Pet Store"),
				},
			},
		},
		"500 - error": {
			init: func(m *apidefinitions.Mock, _ *v0.Mock) {
				mockListEndpointVersionsError(m, 500)
			},
			steps: []resource.TestStep{
				{
					Config:      configWithID(),
					ExpectError: regexp.MustCompile("Error retrieving API"),
				},
			},
		},
		"check schema - id and name are not present": {
			steps: []resource.TestStep{
				{
					Config:      emptyConfig(),
					ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[id,name]`),
				},
			},
		},
		"check schema - id and name are present ": {
			steps: []resource.TestStep{
				{
					Config:      configWithIDAndName(),
					ExpectError: regexp.MustCompile(`Exactly one of these attributes must be configured: \[id,name]`),
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
			clientV0.AssertExpectations(t)
		})
	}
}

func configWithID() string {
	return providerConfig + `
data "akamai_apidefinitions_api" "test" {
	 id = 1
}`
}

func configWithName() string {
	return providerConfig + `
data "akamai_apidefinitions_api" "test" {
	 name = "Pet Store"
}`
}

func emptyConfig() string {
	return providerConfig + `
data "akamai_apidefinitions_api" "test" {
}`
}

func configWithIDAndName() string {
	return providerConfig + `
data "akamai_apidefinitions_api" "test" {
	 id = 1
	 name = "Pet Store"
}`
}

func mockListEndpoints(client *apidefinitions.Mock, ID int64, name string, times int) *mock.Call {
	return client.On("ListEndpoints", mock.Anything, mock.Anything).
		Return(&apidefinitions.ListEndpointsResponse{
			TotalSize: 2,
			APIEndpoints: []apidefinitions.Endpoint{
				{
					APIEndpointID:   ID,
					APIEndpointName: name,
				},
			},
		}, nil).Times(times)
}

func mockListEndpointVersionsError(client *apidefinitions.Mock, code int64) {
	client.On("ListEndpointVersions", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, &apidefinitions.Error{
			Status: code,
			Detail: "Server error:",
		}).Once()
}
