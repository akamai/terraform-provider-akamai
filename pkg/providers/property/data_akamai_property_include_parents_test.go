package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataPropertyIncludeParents(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		init               func(*mockpapi)
		expectedAttributes map[string]string
		expectError        *regexp.Regexp
	}{
		"happy path": {
			givenTF: "valid.tf",
			init: func(m *mockpapi) {
				m.On("ListIncludeParents", mock.Anything, papi.ListIncludeParentsRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					IncludeID:  "inc_1",
				}).Return(&papi.ListIncludeParentsResponse{
					Properties: papi.ParentPropertyItems{
						Items: []papi.ParentProperty{
							{
								PropertyID:                       "prop_1",
								PropertyName:                     "prop_name",
								StagingVersion:                   tools.IntPtr(3),
								ProductionVersion:                tools.IntPtr(2),
								IsIncludeUsedInProductionVersion: true,
								IsIncludeUsedInStagingVersion:    true,
							},
							{
								PropertyID:                       "prop_2",
								PropertyName:                     "some_other_prop_name",
								StagingVersion:                   tools.IntPtr(5),
								ProductionVersion:                nil,
								IsIncludeUsedInStagingVersion:    true,
								IsIncludeUsedInProductionVersion: false,
							},
						},
					},
				}, nil).Times(5)
			},
			expectedAttributes: map[string]string{
				"parents.#": "2",

				"parents.0.id":                                    "prop_1",
				"parents.0.name":                                  "prop_name",
				"parents.0.staging_version":                       "3",
				"parents.0.production_version":                    "2",
				"parents.0.is_include_used_in_staging_version":    "true",
				"parents.0.is_include_used_in_production_version": "true",

				"parents.1.id":                                    "prop_2",
				"parents.1.name":                                  "some_other_prop_name",
				"parents.1.staging_version":                       "5",
				"parents.1.production_version":                    "",
				"parents.1.is_include_used_in_staging_version":    "true",
				"parents.1.is_include_used_in_production_version": "false",
			},
			expectError: nil,
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *mockpapi) {
				m.On("ListIncludeParents", mock.Anything, papi.ListIncludeParentsRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					IncludeID:  "inc_1",
				}).Return(nil, fmt.Errorf("oops"))
			},
			expectError: regexp.MustCompile("oops"),
		},
		"missing required argument contract_id": {
			givenTF:     "missing_contract_id.tf",
			expectError: regexp.MustCompile(`The argument "contract_id" is required, but no definition was found`),
		},
		"missing required argument group_id": {
			givenTF:     "missing_group_id.tf",
			expectError: regexp.MustCompile(`The argument "group_id" is required, but no definition was found`),
		},
		"missing required argument include_id": {
			givenTF:     "missing_include_id.tf",
			expectError: regexp.MustCompile(`The argument "include_id" is required, but no definition was found`),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &mockpapi{}
			if test.init != nil {
				test.init(client)
			}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_parents.parents", k, v))
			}
			useClient(client, nil, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest: true,
					Providers:  testAccProviders,
					Steps: []resource.TestStep{{
						Config:      loadFixtureString(fmt.Sprintf("testdata/TestDataPropertyIncludeParents/%s", test.givenTF)),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
