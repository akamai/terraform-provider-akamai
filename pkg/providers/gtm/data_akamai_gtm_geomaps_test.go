package gtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataGTMGeoMaps(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		init               func(mock *gtm.Mock)
		expectedAttributes map[string]string
		expectError        *regexp.Regexp
	}{
		"happy path": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				mockListGeoMaps(m, []gtm.GeoMap{
					{
						Name: "TestName1",
						DefaultDatacenter: &gtm.DatacenterBase{
							Nickname:     "TestNickname1",
							DatacenterID: 1,
						},
						Assignments: []gtm.GeoAssignment{{
							DatacenterBase: gtm.DatacenterBase{
								Nickname:     "TestNicknameAssignments1",
								DatacenterID: 2,
							},
							Countries: []string{
								"PL",
								"US",
							},
						}},
						Links: []gtm.Link{{
							Rel:  "TestRel1",
							Href: "TestHref1",
						}},
					},
					{
						Name: "TestName2",
						DefaultDatacenter: &gtm.DatacenterBase{
							Nickname:     "TestNickname2",
							DatacenterID: 3,
						},
						Assignments: []gtm.GeoAssignment{
							{
								DatacenterBase: gtm.DatacenterBase{
									Nickname:     "TestNicknameAssignments2",
									DatacenterID: 4,
								},
								Countries: []string{
									"AR",
									"CA",
									"IS",
								},
							},
							{
								DatacenterBase: gtm.DatacenterBase{
									Nickname:     "TestNicknameAssignments3",
									DatacenterID: 5,
								},
								Countries: []string{
									"IT",
								},
							}},
						Links: []gtm.Link{{
							Rel:  "TestRel2",
							Href: "TestHref2",
						}},
					},
				}, nil, testutils.ThreeTimes)
			},
			expectedAttributes: map[string]string{
				"domain":                                 "gtm_terra_testdomain.akadns.net",
				"geo_maps.0.name":                        "TestName1",
				"geo_maps.1.name":                        "TestName2",
				"geo_maps.0.default_datacenter.nickname": "TestNickname1",
				"geo_maps.0.default_datacenter.datacenter_id": "1",
				"geo_maps.1.default_datacenter.nickname":      "TestNickname2",
				"geo_maps.1.default_datacenter.datacenter_id": "3",
				"geo_maps.0.assignments.0.datacenter_id":      "2",
				"geo_maps.0.assignments.0.nickname":           "TestNicknameAssignments1",
				"geo_maps.0.assignments.0.countries.0":        "PL",
				"geo_maps.0.assignments.0.countries.1":        "US",
				"geo_maps.1.assignments.0.datacenter_id":      "4",
				"geo_maps.1.assignments.0.nickname":           "TestNicknameAssignments2",
				"geo_maps.1.assignments.0.countries.0":        "AR",
				"geo_maps.1.assignments.0.countries.1":        "CA",
				"geo_maps.1.assignments.0.countries.2":        "IS",
				"geo_maps.1.assignments.1.datacenter_id":      "5",
				"geo_maps.1.assignments.1.nickname":           "TestNicknameAssignments3",
				"geo_maps.1.assignments.1.countries.0":        "IT",
				"geo_maps.0.links.0.rel":                      "TestRel1",
				"geo_maps.0.links.0.href":                     "TestHref1",
				"geo_maps.1.links.0.rel":                      "TestRel2",
				"geo_maps.1.links.0.href":                     "TestHref2",
			},
			expectError: nil,
		},
		"missing required argument domain": {
			givenTF:     "missing_domain.tf",
			expectError: regexp.MustCompile(`The argument "domain" is required, but no definition was found.`),
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				mockListGeoMaps(m, nil, fmt.Errorf("API error"), testutils.Once)
			},
			expectError: regexp.MustCompile("API error"),
		},
		"no assignments": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				mockListGeoMaps(m, []gtm.GeoMap{{
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
				}}, nil, testutils.ThreeTimes)
			},
			expectedAttributes: map[string]string{
				"domain":          "gtm_terra_testdomain.akadns.net",
				"geo_maps.0.name": "TestName",
				"geo_maps.0.default_datacenter.datacenter_id": "1",
				"geo_maps.0.default_datacenter.nickname":      "TestNickname",
				"geo_maps.0.assignments.#":                    "0",
				"geo_maps.0.links.0.rel":                      "TestRel",
				"geo_maps.0.links.0.href":                     "TestHref",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &gtm.Mock{}
			if test.init != nil {
				test.init(client)
			}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_geomaps.testmaps", k, v))
			}

			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataGtmGeomaps/%s", test.givenTF),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})

			client.AssertExpectations(t)
		})
	}
}

func mockListGeoMaps(client *gtm.Mock, resp []gtm.GeoMap, err error, times int) *mock.Call {
	return client.On("ListGeoMaps",
		testutils.MockContext, gtm.ListGeoMapsRequest{DomainName: testDomainName},
	).Return(resp, err).Times(times)
}
