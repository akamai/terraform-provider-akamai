package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/papi"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataPropertyInclude(t *testing.T) {
	tests := map[string]struct {
		givenTF                   string
		init                      func(*papi.Mock)
		expectedAttributes        map[string]string
		expectedMissingAttributes []string
		expectError               *regexp.Regexp
	}{
		"happy path - both staging and production versions are returned": {
			givenTF: "valid.tf",
			init: func(m *papi.Mock) {
				m.On("GetInclude", testutils.MockContext, papi.GetIncludeRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					IncludeID:  "inc_1",
				}).Return(&papi.GetIncludeResponse{
					Includes: papi.IncludeItems{
						Items: []papi.Include{
							{
								AssetID:           "aid_555",
								IncludeName:       "inc_name",
								IncludeType:       "MICROSERVICES",
								LatestVersion:     4,
								ProductionVersion: ptr.To(3),
								StagingVersion:    ptr.To(2),
							},
						},
					},
					Include: papi.Include{
						AssetID:           "aid_555",
						IncludeName:       "inc_name",
						IncludeType:       "MICROSERVICES",
						LatestVersion:     4,
						ProductionVersion: ptr.To(3),
						StagingVersion:    ptr.To(2),
					},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"asset_id":           "aid_555",
				"name":               "inc_name",
				"type":               "MICROSERVICES",
				"latest_version":     "4",
				"production_version": "3",
				"staging_version":    "2",
			},
			expectError: nil,
		},
		"happy path - missing production version and staging version": {
			givenTF: "valid.tf",
			init: func(m *papi.Mock) {
				m.On("GetInclude", testutils.MockContext, papi.GetIncludeRequest{
					ContractID: "ctr_1",
					GroupID:    "grp_1",
					IncludeID:  "inc_1",
				}).Return(&papi.GetIncludeResponse{
					Includes: papi.IncludeItems{
						Items: []papi.Include{
							{
								AssetID:           "aid_555",
								IncludeName:       "inc_name",
								IncludeType:       "MICROSERVICES",
								LatestVersion:     4,
								ProductionVersion: nil,
								StagingVersion:    nil,
							},
						},
					},
					Include: papi.Include{
						AssetID:           "aid_555",
						IncludeName:       "inc_name",
						IncludeType:       "MICROSERVICES",
						LatestVersion:     4,
						ProductionVersion: nil,
						StagingVersion:    nil,
					},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"name":           "inc_name",
				"type":           "MICROSERVICES",
				"latest_version": "4",
			},
			expectedMissingAttributes: []string{
				"production_version",
				"staging_version",
			},
			expectError: nil,
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *papi.Mock) {
				m.On("GetInclude", testutils.MockContext, papi.GetIncludeRequest{
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
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_property_include.include", k, v))
			}
			for _, v := range test.expectedMissingAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_property_include.include", v))
			}
			useClient(client, nil, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataPropertyInclude/%s", test.givenTF),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
