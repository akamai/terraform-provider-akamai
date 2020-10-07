package property

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"regexp"
	"testing"
)

func TestResourcePropertyVariables_Create(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		expectedAttributes map[string]string
		withError          *regexp.Regexp
	}{
		"all fields set": {
			givenTF: "all_variables.tf",
			expectedAttributes: map[string]string{
				"json": compactJSON(loadFixtureBytes("testdata/TestResourcePorpertyVariables/all_variables.json")),
			},
		},
		"optional fields not set": {
			givenTF: "optional_empty.tf",
			expectedAttributes: map[string]string{
				"json": compactJSON(loadFixtureBytes("testdata/TestResourcePorpertyVariables/optional_empty.json")),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_property_variables.var", k, v))
			}
			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest: true,
					Providers:  testAccProviders,
					Steps: []resource.TestStep{
						{
							Destroy:     false,
							Config:      loadFixtureString(fmt.Sprintf("testdata/TestResourcePorpertyVariables/%s", test.givenTF)),
							Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
							ExpectError: test.withError,
						},
					},
				})
			})
		})
	}
}

func TestResourcePropertyVariables_WithImport(t *testing.T) {
	tests := map[string]struct {
		givenTF          string
		init             func(*mockpapi)
		expectAttributes map[string]string
	}{
		"valid property ID": {
			givenTF: "all_variables.tf",
			init: func(m *mockpapi) {
				m.On("SearchProperties", mock.Anything, papi.SearchRequest{
					Key:   papi.SearchKeyPropertyName,
					Value: "id",
				}).Return(nil, papi.ErrNotFound).Once()
				m.On("SearchProperties", mock.Anything, papi.SearchRequest{
					Key:   papi.SearchKeyHostname,
					Value: "id",
				}).Return(nil, papi.ErrNotFound).Once()
				m.On("SearchProperties", mock.Anything, papi.SearchRequest{
					Key:   papi.SearchKeyEdgeHostname,
					Value: "id",
				}).Return(&papi.SearchResponse{Versions: papi.SearchItems{
					Items: []papi.SearchItem{
						{
							PropertyID: "prp_1",
						},
					},
				}}, nil).Once()
				m.On("GetProperty", mock.Anything, papi.GetPropertyRequest{
					PropertyID: "prp_1",
				}).Return(&papi.GetPropertyResponse{Property: &papi.Property{
					AccountID:     "acc",
					ContractID:    "ctr",
					GroupID:       "grp",
					LatestVersion: 1,
					PropertyID:    "prp_1",
					PropertyName:  "property_name",
				}}, nil)
			},
			expectAttributes: map[string]string{
				"id": "id",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			test.init(client)
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_property_variables.var", k, v))
			}
			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest: true,
					Providers:  testAccProviders,
					Steps: []resource.TestStep{
						{
							Destroy: false,
							Config:  loadFixtureString(fmt.Sprintf("testdata/TestResourcePorpertyVariables/%s", test.givenTF)),
						},
						{
							ImportState:       true,
							ImportStateVerify: true,
							ResourceName:      "akamai_property_variables.var",
							ImportStateId:     "id",
							Check:             resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestResourcePropertyVariables_Update(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		givenTFUpdate      string
		expectedAttributes map[string]string
		withError          *regexp.Regexp
	}{
		"update variable name": {
			givenTF:       "all_variables.tf",
			givenTFUpdate: "update_variables.tf",
			expectedAttributes: map[string]string{
				"json": compactJSON(loadFixtureBytes("testdata/TestResourcePorpertyVariables/update_variables.json")),
			},
		},
		"no variables passed": {
			givenTF:       "all_variables.tf",
			givenTFUpdate: "no_variables.tf",
			expectedAttributes: map[string]string{
				"json": compactJSON(loadFixtureBytes("testdata/TestResourcePorpertyVariables/all_variables.json")),
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("akamai_property_variables.var", k, v))
			}
			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest: true,
					Providers:  testAccProviders,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(fmt.Sprintf("testdata/TestResourcePorpertyVariables/%s", test.givenTF)),
						},
						{
							Config: loadFixtureString(fmt.Sprintf("testdata/TestResourcePorpertyVariables/%s", test.givenTFUpdate)),
							Check:  resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						},
					},
				})
			})
		})
	}
}
