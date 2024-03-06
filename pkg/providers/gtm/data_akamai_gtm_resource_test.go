package gtm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataGTMResource(t *testing.T) {
	tests := map[string]struct {
		givenTF                   string
		init                      func(*gtm.Mock)
		expectedAttributes        map[string]string
		expectedMissingAttributes []string
		expectError               *regexp.Regexp
	}{
		"happy path - GTM data resource should be returned": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("GetResource", mock.Anything, "resource1", "test.domain.net").Return(&gtm.Resource{
					Type:                        "XML load object via HTTP",
					LeastSquaresDecay:           0,
					Description:                 "terraform test resource",
					AggregationType:             "latest",
					LoadImbalancePercentage:     0,
					UpperBound:                  100,
					Name:                        "property",
					MaxUMultiplicativeIncrement: 0,
					DecayRate:                   0,
					Links: []*gtm.Link{{
						Rel: "self",
						Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/" +
							"test.cli.domain.net/resources/resource1",
					},
					},
					ResourceInstances: []*gtm.ResourceInstance{{
						DatacenterID:         3131,
						UseDefaultLoadObject: false,
						LoadObject: gtm.LoadObject{
							LoadObject:     "/test2",
							LoadObjectPort: 80,
							LoadServers:    []string{"2.3.4.5"},
						},
					}},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"aggregation_type":                   "latest",
				"description":                        "terraform test resource",
				"type":                               "XML load object via HTTP",
				"upper_bound":                        "100",
				"decay_rate":                         "0",
				"links.#":                            "1",
				"resource_instances.#":               "1",
				"resource_instances.0.datacenter_id": "3131",
			},
		},
		"missing required argument resource_name": {
			givenTF:     "missing_resource_name.tf",
			expectError: regexp.MustCompile(`The argument "resource_name" is required, but no definition was found`),
		},
		"missing required argument domain": {
			givenTF:     "missing_domain.tf",
			expectError: regexp.MustCompile(`The argument "domain" is required, but no definition was found`),
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("GetResource", mock.Anything, "resource1", "test.domain.net").Return(
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
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_resource.my_gtm_resource", k, v))
			}
			for _, v := range test.expectedMissingAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_gtm_resource.my_gtm_resource", v))
			}

			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataGTMResource/%s", test.givenTF)),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})

			client.AssertExpectations(t)
		})
	}
}
