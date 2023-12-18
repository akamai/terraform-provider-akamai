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

func TestDataGtmDomain(t *testing.T) {
	tests := map[string]struct {
		givenTF                   string
		init                      func(*gtm.Mock)
		expectedAttributes        map[string]string
		expectedMissingAttributes []string
		expectError               *regexp.Regexp
	}{
		"success - response is ok": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("GetDomain", mock.Anything, "test.cli.devexp-terraform.akadns.net").Return(&gtm.Domain{
					Name:                         "test.cli.devexp-terraform.akadns.net",
					CnameCoalescingEnabled:       false,
					DefaultErrorPenalty:          75,
					DefaultHealthMax:             0,
					DefaultHealthMultiplier:      0,
					DefaultHealthThreshold:       0,
					DefaultMaxUnreachablePenalty: 0,
					DefaultTimeoutPenalty:        25,
					DefaultUnreachableThreshold:  0,
					EmailNotificationList:        []string{"ckulinsk@akamai.com"},
					EndUserMappingEnabled:        false,
					LastModified:                 "2023-01-25T10:21:45.000+00:00",
					MaxTTL:                       172800,
					Status: &gtm.ResponseStatus{
						ChangeId:              "ca7e5b1d-1303-42d3-b6c0-8cb62ae849d4",
						Message:               "ERROR: zone is child of existing GTM domain devexp-terraform.akadns.net, which is not allowed",
						PassingValidation:     false,
						PropagationStatus:     "DENIED",
						PropagationStatusDate: "2023-01-25T10:21:00.000+00:00",
					},
					Resources: []*gtm.Resource{{
						AggregationType: "latest",
						Description:     "terraform test resource",
						Type:            "XML load object via HTTP",
						Name:            "test resource",
						UpperBound:      100,
					},
					},
					AsMaps: []*gtm.AsMap{{
						DefaultDatacenter: &gtm.DatacenterBase{
							DatacenterId: 3133,
							Nickname:     "Default (all others)",
						},
						Assignments: []*gtm.AsAssignment{{
							DatacenterBase: gtm.DatacenterBase{
								Nickname:     "New Zone 1",
								DatacenterId: 3133,
							},
							AsNumbers: []int64{
								12222,
								17334,
								16702,
							},
						}},
						Name: "New Map 1",
						Links: []*gtm.Link{{
							Rel:  "self",
							Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/as-maps/DevExpAutomatedTest_6Qil38",
						}},
					},
					},
					CidrMaps: []*gtm.CidrMap{{
						DefaultDatacenter: &gtm.DatacenterBase{
							DatacenterId: 3133,
							Nickname:     "All Other CIDR Blocks",
						},
						Assignments: []*gtm.CidrAssignment{{
							DatacenterBase: gtm.DatacenterBase{
								Nickname:     "New Zone 1",
								DatacenterId: 3133,
							},
							Blocks: []string{
								"1.2.3.4/22",
							},
						}},
						Name: "New Map 1",
						Links: []*gtm.Link{{
							Rel:  "self",
							Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/cidr-maps/New%20Map%201",
						}},
					},
					},
					GeographicMaps: []*gtm.GeoMap{{
						DefaultDatacenter: &gtm.DatacenterBase{
							DatacenterId: 3131,
							Nickname:     "terraform_datacenter_test",
						},
						Assignments: []*gtm.GeoAssignment{{
							DatacenterBase: gtm.DatacenterBase{
								Nickname:     "terraform_datacenter_test_1",
								DatacenterId: 3133,
							},
							Countries: []string{
								"GB",
							},
						}},
						Name: "tfexample_geo_2",
						Links: []*gtm.Link{{
							Rel:  "self",
							Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/geographic-maps/tfexample_geo_2",
						}},
					},
					},
					Links: []*gtm.Link{{
						Rel:  "properties",
						Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/properties",
					}, {
						Rel:  "resources",
						Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/resources"},
					},
					Properties: []*gtm.Property{{
						BalanceByDownloadScore: false,
						DynamicTTL:             60,
						GhostDemandReporting:   false,
						HandoutMode:            "Normal",
						LastModified:           "2023-01-25T09:58:09.000+00:00",
						Name:                   "property",
						Links: []*gtm.Link{{
							Href: "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/properties/property",
							Rel:  "self",
						}},
						LivenessTests: []*gtm.LivenessTest{{
							AnswersRequired:               false,
							DisableNonstandardPortWarning: false,
							HttpError3xx:                  true,
							TestObjectProtocol:            "HTTP",
						}},
						TrafficTargets: []*gtm.TrafficTarget{{
							DatacenterId: 3131,
							Enabled:      true,
							Servers: []string{
								"1.2.3.4",
								"2.3.4.5",
							},
							Weight: 1,
						}},
					}},
				}, nil)
			},
			expectedAttributes: map[string]string{
				"name":                                                           "test.cli.devexp-terraform.akadns.net",
				"cname_coalescing_enabled":                                       "false",
				"default_timeout_penalty":                                        "25",
				"default_error_penalty":                                          "75",
				"email_notification_list.0":                                      "ckulinsk@akamai.com",
				"status.change_id":                                               "ca7e5b1d-1303-42d3-b6c0-8cb62ae849d4",
				"status.message":                                                 "ERROR: zone is child of existing GTM domain devexp-terraform.akadns.net, which is not allowed",
				"status.passing_validation":                                      "false",
				"status.propagation_status":                                      "DENIED",
				"status.propagation_status_date":                                 "2023-01-25T10:21:00.000+00:00",
				"end_user_mapping_enabled":                                       "false",
				"last_modified":                                                  "2023-01-25T10:21:45.000+00:00",
				"max_ttl":                                                        "172800",
				"as_maps.0.name":                                                 "New Map 1",
				"as_maps.0.default_datacenter.datacenter_id":                     "3133",
				"as_maps.0.default_datacenter.nickname":                          "Default (all others)",
				"as_maps.0.assignments.0.nickname":                               "New Zone 1",
				"as_maps.0.assignments.0.datacenter_id":                          "3133",
				"as_maps.0.assignments.0.as_numbers.0":                           "12222",
				"as_maps.0.assignments.0.as_numbers.1":                           "16702",
				"as_maps.0.assignments.0.as_numbers.2":                           "17334",
				"as_maps.0.links.0.href":                                         "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/as-maps/DevExpAutomatedTest_6Qil38",
				"as_maps.0.links.0.rel":                                          "self",
				"cidr_maps.0.name":                                               "New Map 1",
				"cidr_maps.0.default_datacenter.datacenter_id":                   "3133",
				"cidr_maps.0.default_datacenter.nickname":                        "All Other CIDR Blocks",
				"cidr_maps.0.assignments.0.nickname":                             "New Zone 1",
				"cidr_maps.0.assignments.0.datacenter_id":                        "3133",
				"cidr_maps.0.assignments.0.blocks.0":                             "1.2.3.4/22",
				"cidr_maps.0.links.0.href":                                       "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/cidr-maps/New%20Map%201",
				"cidr_maps.0.links.0.rel":                                        "self",
				"geographic_maps.0.name":                                         "tfexample_geo_2",
				"geographic_maps.0.default_datacenter.datacenter_id":             "3131",
				"geographic_maps.0.default_datacenter.nickname":                  "terraform_datacenter_test",
				"geographic_maps.0.assignments.0.nickname":                       "terraform_datacenter_test_1",
				"geographic_maps.0.assignments.0.datacenter_id":                  "3133",
				"geographic_maps.0.assignments.0.countries.0":                    "GB",
				"geographic_maps.0.links.0.href":                                 "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/geographic-maps/tfexample_geo_2",
				"geographic_maps.0.links.0.rel":                                  "self",
				"resources.0.aggregation_type":                                   "latest",
				"resources.0.description":                                        "terraform test resource",
				"resources.0.name":                                               "test resource",
				"resources.0.type":                                               "XML load object via HTTP",
				"resources.0.upper_bound":                                        "100",
				"properties.0.balance_by_download_score":                         "false",
				"properties.0.dynamic_ttl":                                       "60",
				"properties.0.ghost_demand_reporting":                            "false",
				"properties.0.handout_mode":                                      "Normal",
				"properties.0.last_modified":                                     "2023-01-25T09:58:09.000+00:00",
				"properties.0.name":                                              "property",
				"properties.0.links.0.href":                                      "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/properties/property",
				"properties.0.links.0.rel":                                       "self",
				"properties.0.liveness_tests.0.answers_required":                 "false",
				"properties.0.liveness_tests.0.disable_nonstandard_port_warning": "false",
				"properties.0.liveness_tests.0.http_error3xx":                    "true",
				"properties.0.liveness_tests.0.test_object_protocol":             "HTTP",
				"properties.0.traffic_targets.0.datacenter_id":                   "3131",
				"properties.0.traffic_targets.0.enabled":                         "true",
				"properties.0.traffic_targets.0.servers.0":                       "1.2.3.4",
				"properties.0.traffic_targets.0.servers.1":                       "2.3.4.5",
				"properties.0.traffic_targets.0.weight":                          "1",
				"links.0.href":                                                   "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/properties",
				"links.0.rel":                                                    "properties",
				"links.1.href":                                                   "https://akaa-ouijhfns55qwgfuc-knsod5nrjl2w2gmt.luna-dev.akamaiapis.net/config-gtm/v1/domains/test.cli.devexp-terraform.akadns.net/resources",
				"links.1.rel":                                                    "resources",
			},
			expectedMissingAttributes: nil,
			expectError:               nil,
		},
		"missing required argument name": {
			givenTF:     "missing_domain_name.tf",
			expectError: regexp.MustCompile(`The argument "name" is required, but no definition was found`),
		},
		"error response from api": {
			givenTF: "valid.tf",
			init: func(m *gtm.Mock) {
				m.On("GetDomain", mock.Anything, "test.cli.devexp-terraform.akadns.net").Return(nil, fmt.Errorf("oops"))
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
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_domain.domain", k, v))
			}
			for _, v := range test.expectedMissingAttributes {
				checkFuncs = append(checkFuncs, resource.TestCheckNoResourceAttr("data.akamai_gtm_domain.domain", v))
			}
			useClient(client, func() {
				resource.Test(t, resource.TestCase{
					IsUnitTest:               true,
					ProtoV5ProviderFactories: testAccProvidersProtoV5,
					Steps: []resource.TestStep{{
						Config:      testutils.LoadFixtureString(t, fmt.Sprintf("testdata/TestDataGtmDomain/%s", test.givenTF)),
						Check:       resource.ComposeAggregateTestCheckFunc(checkFuncs...),
						ExpectError: test.expectError,
					}},
				})
			})
			client.AssertExpectations(t)
		})
	}
}
