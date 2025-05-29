package apidefinitions

import (
	"regexp"
	"testing"

	v0 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/apidefinitions/v0"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestOpenAPIDataSource(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		configPath string
		init       func(*v0.Mock)
		steps      []resource.TestStep
		error      *regexp.Regexp
	}{
		"200 - ok": {
			init: func(m *v0.Mock) {
				mockFromOpenAPIFile(m, 3)
			},
			steps: []resource.TestStep{
				{
					Config: datasourceConfig(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_apidefinitions_openapi.o1", "api", toStateJSON("with-resources-response.json")),
					),
				},
			},
		},
		"500 - error": {
			init: func(m *v0.Mock) {
				mockFromOpenAPIFileFail(m)
			},
			steps: []resource.TestStep{
				{
					Config:      datasourceConfig(),
					ExpectError: regexp.MustCompile("Mapping OpenAPI File Failed"),
				},
			},
		},
		"check schema - missing required attributes": {
			steps: []resource.TestStep{
				{
					Config:      datasourceConfigEmpty(),
					ExpectError: regexp.MustCompile("Attribute file_path cannot be empty"),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &v0.Mock{}
			if test.init != nil {
				test.init(client)
			}
			useClient(nil, client, func() {
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

func mockFromOpenAPIFile(client *v0.Mock, times int) {
	response, _ := deserialize(readJSON("api-configuration-with-resources.json"))
	client.On("FromOpenAPIFile", mock.Anything, mock.Anything, mock.Anything).
		Return(&v0.FromOpenAPIFileResponse{
			API:      *response,
			Problems: []v0.Error{},
		}, nil).
		Times(times)
}

func mockFromOpenAPIFileFail(client *v0.Mock) {
	client.On("FromOpenAPIFile", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, &internalServerError).
		Once()
}

func datasourceConfig() string {
	return providerConfig + `
data "akamai_apidefinitions_openapi" "o1" {
  file_path = "testdata/endpoint/openapi.yaml"
}
`
}

func datasourceConfigEmpty() string {
	return providerConfig + `
data "akamai_apidefinitions_openapi" "o1" {
  file_path = ""
}
`
}

var internalServerError = v0.Error{
	Status: 500,
	Detail: "Internal server error",
}
