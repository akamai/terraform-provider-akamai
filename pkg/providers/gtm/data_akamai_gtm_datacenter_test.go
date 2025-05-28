package gtm

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataGTMDatacenter(t *testing.T) {
	tests := map[string]struct {
		init       func(*gtm.Mock, *gtm.Datacenter)
		mockData   *gtm.Datacenter
		configPath string
		error      *regexp.Regexp
	}{
		"happy path - all fields populated": {
			init: func(m *gtm.Mock, datacenter *gtm.Datacenter) {
				mockGetDatacenter(m, datacenter.DatacenterID, datacenter, nil, testutils.ThreeTimes)
			},
			mockData:   getTestGTMDatacenter(),
			configPath: "testdata/TestDataGTMDatacenter/default.tf",
		},
		"happy path - minimal fields": {
			init: func(m *gtm.Mock, datacenter *gtm.Datacenter) {
				mockGetDatacenter(m, datacenter.DatacenterID, datacenter, nil, testutils.ThreeTimes)
			},
			mockData:   getMinimalTestDatacenter(),
			configPath: "testdata/TestDataGTMDatacenter/default.tf",
		},
		"happy path - no load_servers in default_load_object": {
			init: func(m *gtm.Mock, datacenter *gtm.Datacenter) {
				mockGetDatacenter(m, datacenter.DatacenterID, datacenter, nil, testutils.ThreeTimes)
			},
			mockData:   getNoLoadServersDatacenter(),
			configPath: "testdata/TestDataGTMDatacenter/default.tf",
		},
		"error - GetDatacenter fail": {
			init: func(m *gtm.Mock, datacenter *gtm.Datacenter) {
				mockGetDatacenter(m, datacenter.DatacenterID, nil, fmt.Errorf("GetDatacenter error"), testutils.Once)
			},
			mockData:   getTestGTMDatacenter(),
			configPath: "testdata/TestDataGTMDatacenter/default.tf",
			error:      regexp.MustCompile("GetDatacenter error"),
		},
		"error - no domain attribute": {
			configPath: "testdata/TestDataGTMDatacenter/no_domain.tf",
			error:      regexp.MustCompile(`The argument "domain" is required, but no definition was found.`),
		},
		"error - no datacenter_id attribute": {
			configPath: "testdata/TestDataGTMDatacenter/no_datacenter_id.tf",
			error:      regexp.MustCompile(`The argument "datacenter_id" is required, but no definition was found.`),
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
							Check:       checkAttrsForGTMDatacenter(test.mockData),
							ExpectError: test.error,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func checkAttrsForGTMDatacenter(datacenter *gtm.Datacenter) resource.TestCheckFunc {
	if datacenter == nil {
		return nil
	}

	const datasourceName = "data.akamai_gtm_datacenter.test"
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(datasourceName, "nickname", datacenter.Nickname),
		resource.TestCheckResourceAttr(datasourceName, "score_penalty", strconv.Itoa(datacenter.ScorePenalty)),
		resource.TestCheckResourceAttr(datasourceName, "city", datacenter.City),
		resource.TestCheckResourceAttr(datasourceName, "state_or_province", datacenter.StateOrProvince),
		resource.TestCheckResourceAttr(datasourceName, "country", datacenter.Country),
		resource.TestCheckResourceAttr(datasourceName, "latitude", strconv.FormatFloat(datacenter.Latitude, 'f', -1, 64)),
		resource.TestCheckResourceAttr(datasourceName, "longitude", strconv.FormatFloat(datacenter.Longitude, 'f', -1, 64)),
		resource.TestCheckResourceAttr(datasourceName, "clone_of", strconv.Itoa(datacenter.CloneOf)),
		resource.TestCheckResourceAttr(datasourceName, "virtual", strconv.FormatBool(datacenter.Virtual)),
		resource.TestCheckResourceAttr(datasourceName, "continent", datacenter.Continent),
		resource.TestCheckResourceAttr(datasourceName, "servermonitor_pool", datacenter.ServermonitorPool),
		resource.TestCheckResourceAttr(datasourceName, "cloud_server_targeting", strconv.FormatBool(datacenter.CloudServerTargeting)),
		resource.TestCheckResourceAttr(datasourceName, "cloud_server_host_header_override", strconv.FormatBool(datacenter.CloudServerHostHeaderOverride)),
	}

	if datacenter.DefaultLoadObject != nil {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(datasourceName, "default_load_object.0.load_servers.#", strconv.Itoa(len(datacenter.DefaultLoadObject.LoadServers))),
			resource.TestCheckResourceAttr(datasourceName, "default_load_object.0.load_object", datacenter.DefaultLoadObject.LoadObject),
			resource.TestCheckResourceAttr(datasourceName, "default_load_object.0.load_object_port", strconv.Itoa(datacenter.DefaultLoadObject.LoadObjectPort)),
		}...)

		for i, server := range datacenter.DefaultLoadObject.LoadServers {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("default_load_object.0.load_servers.%d", i), server))
		}
	}

	if datacenter.Links != nil {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr(datasourceName, "links.#", strconv.Itoa(len(datacenter.Links))))
		for i, link := range datacenter.Links {
			checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("links.%d.rel", i), link.Rel),
				resource.TestCheckResourceAttr(datasourceName, fmt.Sprintf("links.%d.href", i), link.Href),
			}...)

		}
	}

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

func getTestGTMDatacenter() *gtm.Datacenter {
	return &gtm.Datacenter{
		DatacenterID:    1,
		Nickname:        "testNickname",
		ScorePenalty:    2,
		City:            "city",
		StateOrProvince: "state",
		Country:         "country",
		Latitude:        3.3,
		Longitude:       4.4,
		CloneOf:         5,
		Virtual:         true,
		DefaultLoadObject: &gtm.LoadObject{
			LoadObject:     "loadObject",
			LoadObjectPort: 80,
			LoadServers:    []string{"1.1.1.1", "2.2.2.2"},
		},
		Continent:                     "continent",
		ServermonitorPool:             "serverMonitorPool",
		CloudServerTargeting:          true,
		CloudServerHostHeaderOverride: true,
		Links: []gtm.Link{
			{
				Rel:  "rel1",
				Href: "href1",
			},
			{
				Rel:  "rel2",
				Href: "href2",
			},
		},
	}
}

func getMinimalTestDatacenter() *gtm.Datacenter {
	return &gtm.Datacenter{
		DatacenterID:                  1,
		Nickname:                      "testNickname",
		ScorePenalty:                  2,
		Latitude:                      3.3,
		Longitude:                     4.4,
		Virtual:                       true,
		CloudServerTargeting:          true,
		CloudServerHostHeaderOverride: true,
	}
}

func getNoLoadServersDatacenter() *gtm.Datacenter {
	return &gtm.Datacenter{
		DatacenterID: 1,
		Nickname:     "testNickname",
		ScorePenalty: 2,
		Latitude:     3.3,
		Longitude:    4.4,
		Virtual:      true,
		DefaultLoadObject: &gtm.LoadObject{
			LoadObject:     "loadObject",
			LoadObjectPort: 443,
			LoadServers:    nil,
		},
		CloudServerTargeting:          true,
		CloudServerHostHeaderOverride: true,
	}
}
