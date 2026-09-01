package gtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataGTMGeoMap(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		givenTF            string
		init               func(mock *gtm.Mock)
		expectedAttributes map[string]string
		expectError        *regexp.Regexp
	}{
		"happy path": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				mockGetGeoMap(m, testDomainName, &gtm.GeoMap{
					Name: "TestName",
					DefaultDatacenter: &gtm.DatacenterBase{
						Nickname:     "TestNickname",
						DatacenterID: 1,
					},
					Assignments: []gtm.GeoAssignment{{
						DatacenterBase: gtm.DatacenterBase{
							Nickname:     "TestNicknameAssignments",
							DatacenterID: 2,
						},
						Countries: []string{
							"PL",
							"US",
						},
					}},
					Links: []gtm.Link{{
						Rel:  "TestRel",
						Href: "TestHref",
					}},
				}, nil, testutils.ThreeTimes)
			},
			expectedAttributes: map[string]string{
				"domain":                           "gtm_terra_testdomain.akadns.net",
				"map_name":                         "TestName",
				"default_datacenter.datacenter_id": "1",
				"default_datacenter.nickname":      "TestNickname",
				"assignments.0.datacenter_id":      "2",
				"assignments.0.nickname":           "TestNicknameAssignments",
				"assignments.0.countries.0":        "PL",
				"assignments.0.countries.1":        "US",
				"links.0.rel":                      "TestRel",
				"links.0.href":                     "TestHref",
			},
			expectError: nil,
		},
		"missing required argument domain": {
			givenTF:     "missing_domain.tf",
			expectError: regexp.MustCompile(`The argument "domain" is required, but no definition was found.`),
		},
		"missing required argument map_name": {
			givenTF:     "missing_map_name.tf",
			expectError: regexp.MustCompile(`The argument "map_name" is required, but no definition was found.`),
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				mockGetGeoMap(m, testDomainName, nil, fmt.Errorf("API error"), testutils.Once)
			},
			expectError: regexp.MustCompile("API error"),
		},
		"no assignments": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				mockGetGeoMap(m, testDomainName, &gtm.GeoMap{
					Name: "TestName",
					DefaultDatacenter: &gtm.DatacenterBase{
						Nickname:     "TestNickname",
						DatacenterID: 1,
					},
					Assignments: []gtm.GeoAssignment{},
					Links: []gtm.Link{{
						Rel:  "TestRel",
						Href: "TestHref",
					}},
				}, nil, testutils.ThreeTimes)
			},
			expectedAttributes: map[string]string{
				"domain":                           "gtm_terra_testdomain.akadns.net",
				"map_name":                         "TestName",
				"default_datacenter.datacenter_id": "1",
				"default_datacenter.nickname":      "TestNickname",
				"assignments.#":                    "0",
				"links.0.rel":                      "TestRel",
				"links.0.href":                     "TestHref",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if test.init != nil {
				test.init(client.GTM)
			}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_geomap.testmap", k, v))
			}

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataGtmGeomap/%s", test.givenTF),
					Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
					ExpectError: test.expectError,
				}},
			})

			client.GTM.AssertExpectations(t)
		})
	}
}
