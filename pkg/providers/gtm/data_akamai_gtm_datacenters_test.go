package gtm

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataGTMDatacenters(t *testing.T) {
	tests := map[string]struct {
		init       func(*testing.T, *gtm.Mock, testDataForGTMDatacenters)
		mockData   testDataForGTMDatacenters
		configPath string
		error      *regexp.Regexp
	}{
		"happy path - three datacenters": {
			init: func(t *testing.T, m *gtm.Mock, data testDataForGTMDatacenters) {
				mockListDatacenters(t, m, data, 5)
			},
			mockData:   testGTMDatacenters,
			configPath: "testdata/TestDataGTMDatacenters/default.tf",
		},
		"happy path - one datacenter": {
			init: func(t *testing.T, m *gtm.Mock, data testDataForGTMDatacenters) {
				mockListDatacenters(t, m, data, 5)
			},
			mockData:   testGTMSingleDatacenter,
			configPath: "testdata/TestDataGTMDatacenters/default.tf",
		},
		"happy path - no datacenters": {
			init: func(t *testing.T, m *gtm.Mock, data testDataForGTMDatacenters) {
				mockListDatacenters(t, m, data, 5)
			},
			mockData:   testGTMEmptyDatacenters,
			configPath: "testdata/TestDataGTMDatacenters/default.tf",
		},
		"error - ListDatacenters fail": {
			init: func(t *testing.T, m *gtm.Mock, data testDataForGTMDatacenters) {
				m.On("ListDatacenters", mock.Anything, data.domain).Return(
					nil, fmt.Errorf("ListDatacenters error")).Once()
			},
			mockData:   testGTMDatacenters,
			configPath: "testdata/TestDataGTMDatacenters/default.tf",
			error:      regexp.MustCompile("ListDatacenters error"),
		},
		"error - no domain attribute": {
			init:       func(t *testing.T, _ *gtm.Mock, _ testDataForGTMDatacenters) {},
			configPath: "testdata/TestDataGTMDatacenters/no_domain.tf",
			error:      regexp.MustCompile(`The argument "domain" is required, but no definition was found.`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &gtm.Mock{}
			test.init(t, client, test.mockData)
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

func checkAttrsForGTMDatacenters(data testDataForGTMDatacenters) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", "id", data.domain))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", "datacenters.#", strconv.Itoa(len(data.datacenters))))

	for i, dc := range data.datacenters {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.nickname", i), dc.nickname))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.datacenter_id", i), strconv.Itoa(dc.datacenterID)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.score_penalty", i), strconv.Itoa(dc.scorePenalty)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.city", i), dc.city))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.state_or_province", i), dc.stateOrProvince))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.country", i), dc.country))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.latitude", i), strconv.FormatFloat(dc.latitude, 'f', -1, 64)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.longitude", i), strconv.FormatFloat(dc.longitude, 'f', -1, 64)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.clone_of", i), strconv.Itoa(dc.cloneOf)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.virtual", i), strconv.FormatBool(dc.virtual)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.clone_of", i), strconv.Itoa(dc.cloneOf)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.continent", i), dc.continent))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.servermonitor_pool", i), dc.serverMonitorPool))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.cloud_server_targeting", i), strconv.FormatBool(dc.cloudServerTargeting)))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.cloud_server_host_header_override", i), strconv.FormatBool(dc.cloudServerHostHeaderOverride)))

		if dc.defaultLoadObject != nil {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.default_load_object.0.load_servers.#", i), strconv.Itoa(len(dc.defaultLoadObject.LoadServers))))
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.default_load_object.0.load_object", i), dc.defaultLoadObject.LoadObject))
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.default_load_object.0.load_object_port", i), strconv.Itoa(dc.defaultLoadObject.LoadObjectPort)))
			for j, server := range dc.defaultLoadObject.LoadServers {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.default_load_object.0.load_servers.%d", i, j), server))
			}
		}

		if dc.links != nil {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.links.#", i), strconv.Itoa(len(dc.links))))
			for j, link := range dc.links {
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.links.%d.rel", i, j), link.Rel))
				checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenters.test", fmt.Sprintf("datacenters.%d.links.%d.href", i, j), link.Href))
			}
		}
	}

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

type testDataForGTMDatacenters struct {
	domain      string
	datacenters []testDataForGTMDatacenter
}

var (
	testGTMDatacenters = testDataForGTMDatacenters{
		domain: "test.domain.com",
		datacenters: []testDataForGTMDatacenter{
			testGTMDatacenter,
			testGTMDatacenterMinimal,
			testGTMDatacenterNoLoadServers,
		},
	}

	testGTMSingleDatacenter = testDataForGTMDatacenters{
		domain:      "test.domain.com",
		datacenters: []testDataForGTMDatacenter{testGTMDatacenter},
	}

	testGTMEmptyDatacenters = testDataForGTMDatacenters{
		domain:      "test.domain.com",
		datacenters: []testDataForGTMDatacenter{},
	}

	// mockListDatacenters mocks ListDatacenters call with provided data
	mockListDatacenters = func(t *testing.T, client *gtm.Mock, data testDataForGTMDatacenters, timesToRun int) {
		var dcs []*gtm.Datacenter

		for _, data := range data.datacenters {
			dc := &gtm.Datacenter{
				City:                          data.city,
				CloneOf:                       data.cloneOf,
				CloudServerHostHeaderOverride: data.cloudServerHostHeaderOverride,
				CloudServerTargeting:          data.cloudServerTargeting,
				Continent:                     data.continent,
				Country:                       data.country,
				DefaultLoadObject:             data.defaultLoadObject,
				Latitude:                      data.latitude,
				Links:                         data.links,
				Longitude:                     data.longitude,
				Nickname:                      data.nickname,
				DatacenterID:                  data.datacenterID,
				ScorePenalty:                  data.scorePenalty,
				ServermonitorPool:             data.serverMonitorPool,
				StateOrProvince:               data.stateOrProvince,
				Virtual:                       data.virtual,
			}
			dcs = append(dcs, dc)
		}

		client.On("ListDatacenters", mock.Anything, data.domain).Return(dcs, nil).Times(timesToRun)
	}
)
