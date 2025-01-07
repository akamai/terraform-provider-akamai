package gtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataGTMCIDRmap(t *testing.T) {
	tests := map[string]struct {
		givenTF            string
		init               func(mock *gtm.Mock)
		expectedAttributes map[string]string
		expectError        *regexp.Regexp
	}{
		"happy path": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("GetCIDRMap", testutils.MockContext, gtm.GetCIDRMapRequest{
					MapName:    "mapTest",
					DomainName: "test.cidrmap.domain.net",
				}).Return(&gtm.GetCIDRMapResponse{
					Name: "TestName",
					DefaultDatacenter: &gtm.DatacenterBase{
						Nickname:     "TestNickname",
						DatacenterID: 1,
					},
					Assignments: []gtm.CIDRAssignment{{
						DatacenterBase: gtm.DatacenterBase{
							Nickname:     "TestNicknameAssignments",
							DatacenterID: 1,
						},
						Blocks: []string{
							"test1",
							"test2",
						},
					}},
					Links: []gtm.Link{{
						Rel:  "TestRel",
						Href: "TestHref",
					}},
				}, nil)

			},
			expectedAttributes: map[string]string{
				"domain":                           "test.cidrmap.domain.net",
				"map_name":                         "TestName",
				"default_datacenter.datacenter_id": "1",
				"default_datacenter.nickname":      "TestNickname",
				"assignments.0.datacenter_id":      "1",
				"assignments.0.nickname":           "TestNicknameAssignments",
				"assignments.0.blocks.0":           "test1",
				"assignments.0.blocks.1":           "test2",
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
				m.On("GetCIDRMap", testutils.MockContext, gtm.GetCIDRMapRequest{
					MapName:    "mapTest",
					DomainName: "test.cidrmap.domain.net",
				}).Return(
					nil, fmt.Errorf("error"))
			},
			expectError: regexp.MustCompile("error"),
		},
		"no assignments": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("GetCIDRMap", testutils.MockContext, gtm.GetCIDRMapRequest{
					MapName:    "mapTest",
					DomainName: "test.cidrmap.domain.net",
				}).Return(&gtm.GetCIDRMapResponse{
					Name: "TestName",
					DefaultDatacenter: &gtm.DatacenterBase{
						Nickname:     "TestNickname",
						DatacenterID: 1,
					},
					Assignments: []gtm.CIDRAssignment{},
					Links: []gtm.Link{{
						Rel:  "TestRel",
						Href: "TestHref",
					}},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"domain":                           "test.cidrmap.domain.net",
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
			client := &gtm.Mock{}
			if test.init != nil {
				test.init(client)
			}
			var checkFuncs []resource.TestCheckFunc
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_cidrmap.gtm_cidrmap", k, v))
			}

			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataGtmCidrmap/%s", test.givenTF)),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})

			client.AssertExpectations(t)
		})
	}
}
