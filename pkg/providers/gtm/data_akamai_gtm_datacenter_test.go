package gtm

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDataGTMDatacenter(t *testing.T) {
	tests := map[string]struct {
		init       func(*testing.T, *gtm.Mock, testDataForGTMDatacenter)
		mockData   testDataForGTMDatacenter
		configPath string
		error      *regexp.Regexp
	}{
		"happy path - all fields populated": {
			init: func(t *testing.T, m *gtm.Mock, data testDataForGTMDatacenter) {
				mockGetDatacenter(t, m, data, 5)
			},
			mockData:   testGTMDatacenter,
			configPath: "testdata/TestDataGTMDatacenter/default.tf",
		},
		"happy path - minimal fields": {
			init: func(t *testing.T, m *gtm.Mock, data testDataForGTMDatacenter) {
				mockGetDatacenter(t, m, data, 5)
			},
			mockData:   testGTMDatacenterMinimal,
			configPath: "testdata/TestDataGTMDatacenter/default.tf",
		},
		"happy path - no load_servers in default_load_object": {
			init: func(t *testing.T, m *gtm.Mock, data testDataForGTMDatacenter) {
				mockGetDatacenter(t, m, data, 5)
			},
			mockData:   testGTMDatacenterNoLoadServers,
			configPath: "testdata/TestDataGTMDatacenter/default.tf",
		},
		"error - GetDatacenter fail": {
			init: func(t *testing.T, m *gtm.Mock, data testDataForGTMDatacenter) {
				m.On("GetDatacenter", mock.Anything, data.datacenterID, data.domain).Return(
					nil, fmt.Errorf("GetDatacenter error")).Once()
			},
			mockData:   testGTMDatacenter,
			configPath: "testdata/TestDataGTMDatacenter/default.tf",
			error:      regexp.MustCompile("GetDatacenter error"),
		},
		"error - no domain attribute": {
			init:       func(t *testing.T, _ *gtm.Mock, _ testDataForGTMDatacenter) {},
			configPath: "testdata/TestDataGTMDatacenter/no_domain.tf",
			error:      regexp.MustCompile(`The argument "domain" is required, but no definition was found.`),
		},
		"error - no datacenter_id attribute": {
			init:       func(t *testing.T, _ *gtm.Mock, _ testDataForGTMDatacenter) {},
			configPath: "testdata/TestDataGTMDatacenter/no_datacenter_id.tf",
			error:      regexp.MustCompile(`The argument "datacenter_id" is required, but no definition was found.`),
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

func checkAttrsForGTMDatacenter(data testDataForGTMDatacenter) resource.TestCheckFunc {
	var checkFuncs []resource.TestCheckFunc

	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "nickname", data.nickname))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "score_penalty", strconv.Itoa(data.scorePenalty)))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "city", data.city))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "state_or_province", data.stateOrProvince))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "country", data.country))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "latitude", strconv.FormatFloat(data.latitude, 'f', -1, 64)))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "longitude", strconv.FormatFloat(data.longitude, 'f', -1, 64)))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "clone_of", strconv.Itoa(data.cloneOf)))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "virtual", strconv.FormatBool(data.virtual)))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "clone_of", strconv.Itoa(data.cloneOf)))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "continent", data.continent))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "servermonitor_pool", data.serverMonitorPool))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "cloud_server_targeting", strconv.FormatBool(data.cloudServerTargeting)))
	checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "cloud_server_host_header_override", strconv.FormatBool(data.cloudServerHostHeaderOverride)))

	if data.defaultLoadObject != nil {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "default_load_object.0.load_servers.#", strconv.Itoa(len(data.defaultLoadObject.LoadServers))))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "default_load_object.0.load_object", data.defaultLoadObject.LoadObject))
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "default_load_object.0.load_object_port", strconv.Itoa(data.defaultLoadObject.LoadObjectPort)))
		for i, server := range data.defaultLoadObject.LoadServers {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", fmt.Sprintf("default_load_object.0.load_servers.%d", i), server))
		}
	}

	if data.links != nil {
		checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", "links.#", strconv.Itoa(len(data.links))))
		for i, link := range data.links {
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", fmt.Sprintf("links.%d.rel", i), link.Rel))
			checkFuncs = append(checkFuncs, resource.TestCheckResourceAttr("data.akamai_gtm_datacenter.test", fmt.Sprintf("links.%d.href", i), link.Href))
		}
	}

	return resource.ComposeAggregateTestCheckFunc(checkFuncs...)
}

type testDataForGTMDatacenter struct {
	domain                        string
	datacenterID                  int
	nickname                      string
	scorePenalty                  int
	city                          string
	stateOrProvince               string
	country                       string
	latitude                      float64
	longitude                     float64
	cloneOf                       int
	virtual                       bool
	defaultLoadObject             *gtm.LoadObject
	continent                     string
	serverMonitorPool             string
	cloudServerTargeting          bool
	cloudServerHostHeaderOverride bool
	links                         []*gtm.Link
}

var (
	// mockGetDatacenter mocks GetDatacenter call with provided data
	mockGetDatacenter = func(t *testing.T, client *gtm.Mock, data testDataForGTMDatacenter, timesToRun int) {
		dc := gtm.Datacenter{
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
		client.On("GetDatacenter", mock.Anything, data.datacenterID, data.domain).Return(&dc, nil).Times(timesToRun)
	}

	testGTMDatacenter = testDataForGTMDatacenter{
		domain:          "test.domain.com",
		datacenterID:    1,
		nickname:        "testNickname",
		scorePenalty:    2,
		city:            "city",
		stateOrProvince: "state",
		country:         "country",
		latitude:        3.3,
		longitude:       4.4,
		cloneOf:         5,
		virtual:         true,
		defaultLoadObject: &gtm.LoadObject{
			LoadObject:     "loadObject",
			LoadObjectPort: 80,
			LoadServers:    []string{"1.1.1.1", "2.2.2.2"},
		},
		continent:                     "continent",
		serverMonitorPool:             "serverMonitorPool",
		cloudServerTargeting:          true,
		cloudServerHostHeaderOverride: true,
		links: []*gtm.Link{
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

	testGTMDatacenterMinimal = testDataForGTMDatacenter{
		domain:                        "test.domain.com",
		datacenterID:                  1,
		nickname:                      "testNickname",
		scorePenalty:                  2,
		latitude:                      3.3,
		longitude:                     4.4,
		virtual:                       true,
		cloudServerTargeting:          true,
		cloudServerHostHeaderOverride: true,
	}

	testGTMDatacenterNoLoadServers = testDataForGTMDatacenter{
		domain:       "test.domain.com",
		datacenterID: 1,
		nickname:     "testNickname",
		scorePenalty: 2,
		latitude:     3.3,
		longitude:    4.4,
		virtual:      true,
		defaultLoadObject: &gtm.LoadObject{
			LoadObject:     "loadObject",
			LoadObjectPort: 443,
			LoadServers:    nil,
		},
		cloudServerTargeting:          true,
		cloudServerHostHeaderOverride: true,
	}
)
