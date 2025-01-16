package gtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataGTMASMap(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		init               func(*gtm.Mock)
		expectedAttributes map[string]string
		expectError        *regexp.Regexp
	}{
		"happy path": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("GetASMap", testutils.MockContext, gtm.GetASMapRequest{
					ASMapName:  "map1",
					DomainName: "test.domain.net",
				}).Return(&gtm.GetASMapResponse{
					Name: "TestName",
					DefaultDatacenter: &gtm.DatacenterBase{
						Nickname:     "TestDefaultDatacenterNickname",
						DatacenterID: 1,
					},
					Assignments: []gtm.ASAssignment{{
						DatacenterBase: gtm.DatacenterBase{
							Nickname:     "TestAssignmentNickname",
							DatacenterID: 1,
						},
						ASNumbers: []int64{
							1,
							2,
							3,
						},
					}},
					Links: []gtm.Link{{
						Href: "href.test",
						Rel:  "TestRel",
					}},
				}, nil)

			},
			expectedAttributes: map[string]string{
				"domain":                           "test.domain.net",
				"map_name":                         "TestName",
				"default_datacenter.datacenter_id": "1",
				"default_datacenter.nickname":      "TestDefaultDatacenterNickname",
				"assignments.0.datacenter_id":      "1",
				"assignments.0.nickname":           "TestAssignmentNickname",
				"assignments.0.as_numbers.0":       "1",
				"assignments.0.as_numbers.1":       "2",
				"assignments.0.as_numbers.2":       "3",
				"links.0.rel":                      "TestRel",
				"links.0.href":                     "href.test",
			},
			expectError: nil,
		},
		"missing required argument domain": {
			givenTF:     "missing_domain.tf",
			expectError: regexp.MustCompile(`The argument "domain" is required, but no definition was found`),
		},
		"missing required argument map_name": {
			givenTF:     "missing_map_name.tf",
			expectError: regexp.MustCompile(`The argument "map_name" is required, but no definition was found`),
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("GetASMap", testutils.MockContext, gtm.GetASMapRequest{
					ASMapName:  "map1",
					DomainName: "test.domain.net",
				}).Return(
					nil, fmt.Errorf("test error"))
			},
			expectError: regexp.MustCompile("test error"),
		},
		"no assignments": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("GetASMap", testutils.MockContext, gtm.GetASMapRequest{
					ASMapName:  "map1",
					DomainName: "test.domain.net",
				}).Return(&gtm.GetASMapResponse{
					Name: "TestName",
					DefaultDatacenter: &gtm.DatacenterBase{
						Nickname:     "TestDefaultDatacenterNickname",
						DatacenterID: 1,
					},
					Assignments: []gtm.ASAssignment{},
					Links: []gtm.Link{{
						Href: "href.test",
						Rel:  "TestRel",
					},
					},
				}, nil)

			},
			expectedAttributes: map[string]string{
				"map_name":                         "TestName",
				"default_datacenter.datacenter_id": "1",
				"default_datacenter.nickname":      "TestDefaultDatacenterNickname",
				"assignments.#":                    "0",
				"links.0.rel":                      "TestRel",
				"links.0.href":                     "href.test",
			},
			expectError: nil,
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
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_asmap.my_gtm_asmap", k, v))
			}

			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataGtmAsmap/%s", test.givenTF),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})

			client.AssertExpectations(t)
		})
	}
}
