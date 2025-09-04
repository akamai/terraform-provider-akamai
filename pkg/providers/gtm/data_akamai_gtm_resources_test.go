package gtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataGTMResources(t *testing.T) {
	tests := map[string]struct {
		givenTF                   string
		init                      func(*gtm.Mock)
		expectedAttributes        map[string]string
		expectedMissingAttributes []string
		expectError               *regexp.Regexp
	}{
		"happy path - GTM data resources should be returned": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				mockListResources(m, []gtm.Resource{
					{
						Type:                        "XML load object via HTTP",
						LeastSquaresDecay:           0,
						Description:                 "terraform test resource1",
						AggregationType:             "latest1",
						LoadImbalancePercentage:     0,
						UpperBound:                  100,
						Name:                        "property1",
						MaxUMultiplicativeIncrement: 0,
						DecayRate:                   0,
						Links: []gtm.Link{{
							Rel: "self",
							Href: "https://test.domain1.net/config-gtm/v1/domains/" +
								"test.cli.domain.net/resources/resource1",
						},
						},
						ResourceInstances: []gtm.ResourceInstance{{
							DatacenterID:         datacenterID3131,
							UseDefaultLoadObject: false,
							LoadObject: gtm.LoadObject{
								LoadObject:     "/test1",
								LoadObjectPort: 80,
								LoadServers:    []string{"2.3.4.5"},
							},
						},
						},
					},
					{
						Type:                        "XML load object via HTTP",
						LeastSquaresDecay:           0,
						Description:                 "terraform test resource2",
						AggregationType:             "latest2",
						LoadImbalancePercentage:     0,
						UpperBound:                  100,
						Name:                        "property2",
						MaxUMultiplicativeIncrement: 0,
						DecayRate:                   0,
						Links: []gtm.Link{{
							Rel: "self1",
							Href: "https://test.domain1.net/config-gtm/v1/domains/" +
								"test.cli.domain.net/resources/resource2",
						},
						},
						ResourceInstances: []gtm.ResourceInstance{{
							DatacenterID:         datacenterID3132,
							UseDefaultLoadObject: false,
							LoadObject: gtm.LoadObject{
								LoadObject:     "/test2",
								LoadObjectPort: 80,
								LoadServers:    []string{"2.3.4.5"},
							},
						},
						},
					},
				}, nil, testutils.ThreeTimes)
			},
			expectedAttributes: map[string]string{
				"resources.0.aggregation_type":                   "latest1",
				"resources.0.description":                        "terraform test resource1",
				"resources.0.resource_instances.0.datacenter_id": "3131",
				"resources.1.aggregation_type":                   "latest2",
				"resources.1.description":                        "terraform test resource2",
				"resources.1.resource_instances.0.datacenter_id": "3132",
			},
		},
		"missing required argument domain": {
			givenTF:     "missing_domain.tf",
			expectError: regexp.MustCompile(`The argument "domain" is required, but no definition was found`),
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				mockListResources(m, nil, fmt.Errorf("oops"), testutils.Once)
			},
			expectError: regexp.MustCompile("oops"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &gtm.Mock{}
			if test.init != nil {
				test.init(client)
			}
			var checkFuncs []resource.TestCheckFunc
			const datasourceName = "data.akamai_gtm_resources.my_gtm_resources"
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(datasourceName, k, v))
			}
			for _, v := range test.expectedMissingAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr(datasourceName, v))
			}

			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureStringf(t, "testdata/TestDataGTMResources/%s", test.givenTF),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})

			client.AssertExpectations(t)
		})
	}
}

func mockListResources(client *gtm.Mock, resp []gtm.Resource, err error, times int) *mock.Call {
	return client.On("ListResources", testutils.MockContext,
		gtm.ListResourcesRequest{
			DomainName: testDomainName,
		}).Return(resp, err).Times(times)
}
