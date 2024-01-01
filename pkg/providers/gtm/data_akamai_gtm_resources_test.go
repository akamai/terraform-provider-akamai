package gtm

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"regexp"
	"testing"
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
				m.On("ListResources", mock.Anything, "test.domain.net").Return([]*gtm.Resource{
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
						Links: []*gtm.Link{{
							Rel: "self",
							Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/" +
								"test.cli.domain.net/resources/resource1",
						},
						},
						ResourceInstances: []*gtm.ResourceInstance{{
							DatacenterId:         3131,
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
						Links: []*gtm.Link{{
							Rel: "self1",
							Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/" +
								"test.cli.domain.net/resources/resource2",
						},
						},
						ResourceInstances: []*gtm.ResourceInstance{{
							DatacenterId:         3132,
							UseDefaultLoadObject: false,
							LoadObject: gtm.LoadObject{
								LoadObject:     "/test2",
								LoadObjectPort: 80,
								LoadServers:    []string{"2.3.4.5"},
							},
						},
						},
					},
				}, nil)
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
				m.On("ListResources", mock.Anything, "test.domain.net").Return(
					nil, fmt.Errorf("oops"))
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
			for k, v := range test.expectedAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_resources.my_gtm_resources", k, v))
			}
			for _, v := range test.expectedMissingAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_gtm_resources.my_gtm_resources", v))
			}

			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV5ProviderFactories: testAccProvidersProtoV5,
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataGTMResources/%s", test.givenTF)),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})

			client.AssertExpectations(t)
		})
	}
}
