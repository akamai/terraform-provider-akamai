package gtm

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataGTMDatacenters(t *testing.T) {
	tests := map[string]struct {
		init       func(*gtm.Mock, []gtm.Datacenter)
		mockData   []gtm.Datacenter
		configPath string
		error      *regexp.Regexp
	}{
		"happy path - three datacenters": {
			init: func(m *gtm.Mock, datacenters []gtm.Datacenter) {
				mockListDatacenters(m, datacenters, nil, testutils.ThreeTimes)
			},
			mockData:   getTestGTMDatacenters(),
			configPath: "testdata/TestDataGTMDatacenters/default.tf",
		},
		"happy path - one datacenter": {
			init: func(m *gtm.Mock, datacenters []gtm.Datacenter) {
				mockListDatacenters(m, datacenters, nil, testutils.ThreeTimes)
			},
			mockData:   getSingleTestDatacenters(),
			configPath: "testdata/TestDataGTMDatacenters/default.tf",
		},
		"happy path - no datacenters": {
			init: func(m *gtm.Mock, datacenters []gtm.Datacenter) {
				mockListDatacenters(m, datacenters, nil, testutils.ThreeTimes)
			},
			mockData:   getEmptyTestDatacenters(),
			configPath: "testdata/TestDataGTMDatacenters/default.tf",
		},
		"error - ListDatacenters fail": {
			init: func(m *gtm.Mock, _ []gtm.Datacenter) {
				mockListDatacenters(m, nil, fmt.Errorf("ListDatacenters error"), testutils.Once)
			},
			mockData:   getTestGTMDatacenters(),
			configPath: "testdata/TestDataGTMDatacenters/default.tf",
			error:      regexp.MustCompile("ListDatacenters error"),
		},
		"error - no domain attribute": {
			configPath: "testdata/TestDataGTMDatacenters/no_domain.tf",
			error:      regexp.MustCompile(`The argument "domain" is required, but no definition was found.`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &gtm.Mock{}
			if test.init != nil {
				test.init(client, test.mockData)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config:      testutils.LoadFixtureString(t, test.configPath),
							Check:       checkAttrsForGTMDatacenters(test.mockData),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkAttrsForGTMDatacenters(datacenters []gtm.Datacenter) resource.TestCheckFunc {
	const datasourceName = "data.akamai_gtm_datacenters.test"
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(datasourceName, "id", testDomainName),
		resource.TestCheckResourceAttr(datasourceName, "datacenters.#", strconv.Itoa(len(datacenters))),
	}

	for i, dc := range datacenters {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.nickname", i), dc.Nickname),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.datacenter_id", i), strconv.Itoa(dc.DatacenterID)),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.score_penalty", i), strconv.Itoa(dc.ScorePenalty)),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.city", i), dc.City),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.state_or_province", i), dc.StateOrProvince),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.country", i), dc.Country),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.latitude", i), strconv.FormatFloat(dc.Latitude, 'f', -1, 64)),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.longitude", i), strconv.FormatFloat(dc.Longitude, 'f', -1, 64)),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.clone_of", i), strconv.Itoa(dc.CloneOf)),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.virtual", i), strconv.FormatBool(dc.Virtual)),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.continent", i), dc.Continent),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.servermonitor_pool", i), dc.ServermonitorPool),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.cloud_server_targeting", i), strconv.FormatBool(dc.CloudServerTargeting)),
			resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.cloud_server_host_header_override", i), strconv.FormatBool(dc.CloudServerHostHeaderOverride)),
		}...)

		if dc.DefaultLoadObject != nil {
			checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.default_load_object.0.load_servers.#", i), strconv.Itoa(len(dc.DefaultLoadObject.LoadServers))),
				resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.default_load_object.0.load_object", i), dc.DefaultLoadObject.LoadObject),
				resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.default_load_object.0.load_object_port", i), strconv.Itoa(dc.DefaultLoadObject.LoadObjectPort)),
			}...)

			for j, server := range dc.DefaultLoadObject.LoadServers {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.default_load_object.0.load_servers.%d", i, j), server))
			}
		}

		if dc.Links != nil {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.links.#", i), strconv.Itoa(len(dc.Links))))
			for j, link := range dc.Links {
				checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
					resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.links.%d.rel", i, j), link.Rel),
					resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("datacenters.%d.links.%d.href", i, j), link.Href),
				}...)
			}
		}
	}

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func mockListDatacenters(client *gtm.Mock, dcs []gtm.Datacenter, err error, times int) *mock.Call {
	return client.On("ListDatacenters", testutils.MockContext, gtm.ListDatacentersRequest{
		DomainName: testDomainName,
	}).Return(dcs, err).Times(times)
}

func getTestGTMDatacenters() []gtm.Datacenter {
	return []gtm.Datacenter{
		*getTestDatacenter(),
		*getMinimalTestDatacenter(),
		*getNoLoadServersDatacenter(),
	}
}

func getSingleTestDatacenters() []gtm.Datacenter {
	return []gtm.Datacenter{*getTestDatacenter()}
}

func getEmptyTestDatacenters() []gtm.Datacenter {
	return []gtm.Datacenter{}
}
