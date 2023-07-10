package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataPropertyIncludeParents(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		init               func(*papi.Mock)
		expectedAttributes map[string]string
		expectError        *regexp.Regexp
	}{
		"happy path": {
			givenTF: "valid.tf",
			init: func(m *papi.Mock) {
				m.On("ListIncludeParents", mock.Anything, papi.ListIncludeParentsRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					IncludeID:  "inc_1",
				}).Return(&papi.ListIncludeParentsResponse{
					Properties: papi.ParentPropertyItems{
						Items: []papi.ParentProperty{
							{
								PropertyID:        "prop_1",
								PropertyName:      "prop_name",
								StagingVersion:    tools.IntPtr(3),
								ProductionVersion: tools.IntPtr(2),
							},
							{
								PropertyID:        "prop_2",
								PropertyName:      "some_other_prop_name",
								StagingVersion:    tools.IntPtr(5),
								ProductionVersion: nil,
							},
							{
								PropertyID:        "prop_3",
								PropertyName:      "third_prop_name",
								StagingVersion:    tools.IntPtr(5),
								ProductionVersion: tools.IntPtr(5),
							},
						},
					},
				}, nil).Times(5)
				// run ListReferencedIncludes for each IncludeParent with different and not empty StagingVersion and ProductionVersion
				m.On("ListReferencedIncludes", mock.Anything, papi.ListReferencedIncludesRequest{
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
					PropertyVersion: 3,
					PropertyID:      "prop_1",
				}).Return(&papi.ListReferencedIncludesResponse{
					Includes: papi.IncludeItems{
						Items: []papi.Include{
							{
								AccountID:         "test_account",
								AssetID:           "test_asset",
								ContractID:        "ctr_1",
								GroupID:           "grp_1",
								IncludeID:         "inc_1",
								IncludeName:       "test_include",
								IncludeType:       papi.IncludeTypeMicroServices,
								LatestVersion:     1,
								ProductionVersion: tools.IntPtr(2),
								StagingVersion:    tools.IntPtr(3),
							},
						},
					},
				}, nil).Times(5)
				m.On("ListReferencedIncludes", mock.Anything, papi.ListReferencedIncludesRequest{
					ContractID:      "ctr_1",
					GroupID:         "grp_1",
					PropertyVersion: 2,
					PropertyID:      "prop_1",
				}).Return(&papi.ListReferencedIncludesResponse{
					Includes: papi.IncludeItems{
						Items: []papi.Include{
							{
								AccountID:         "test_account",
								AssetID:           "test_asset",
								ContractID:        "ctr_1",
								GroupID:           "grp_1",
								IncludeID:         "inc_2",
								IncludeName:       "test_include_2",
								IncludeType:       papi.IncludeTypeMicroServices,
								LatestVersion:     1,
								ProductionVersion: tools.IntPtr(2),
								StagingVersion:    tools.IntPtr(3),
							},
						},
					},
				}, nil).Times(5)
			},
			expectedAttributes: map[string]string{
				"parents.#": "3",

				"parents.0.id":                                    "prop_1",
				"parents.0.name":                                  "prop_name",
				"parents.0.staging_version":                       "3",
				"parents.0.production_version":                    "2",
				"parents.0.is_include_used_in_staging_version":    "true",
				"parents.0.is_include_used_in_production_version": "false",

				"parents.1.id":                                    "prop_2",
				"parents.1.name":                                  "some_other_prop_name",
				"parents.1.staging_version":                       "5",
				"parents.1.production_version":                    "",
				"parents.1.is_include_used_in_staging_version":    "true",
				"parents.1.is_include_used_in_production_version": "false",

				"parents.2.id":                                    "prop_3",
				"parents.2.name":                                  "third_prop_name",
				"parents.2.staging_version":                       "5",
				"parents.2.production_version":                    "5",
				"parents.2.is_include_used_in_staging_version":    "true",
				"parents.2.is_include_used_in_production_version": "true",
			},
			expectError: nil,
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *papi.Mock) {
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
			client := &papi.Mock{}
			if test.init != nil {
				test.init(client)
			}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_property_include_parents.parents", k, v))
			}
			useClient(client, nil, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV5ProviderFactories: testAccProviders,
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataPropertyIncludeParents/%s", test.givenTF)),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
