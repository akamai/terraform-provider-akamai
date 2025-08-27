package gtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

const (
	datacenterID5400 = 5400
	datacenterID3131 = 3131
	datacenterID3132 = 3132
	datacenterID3133 = 3133
)

func TestResGTMDatacenter(t *testing.T) {

	t.Run("create datacenter", func(t *testing.T) {
		client := &gtm.Mock{}

		mockCreateDatacenter(client, &gtm.CreateDatacenterResponse{
			Resource: getTestDatacenterResp(),
			Status:   getPendingResponseStatus(),
		}, nil)

		mockGetDatacenter(client, datacenterID3132, getTestDatacenterResp(), nil, testutils.FourTimes)

		mockGetDomainStatus(client, testutils.Twice)

		mockUpdateDatacenter(client, &gtm.UpdateDatacenterResponse{Status: getDefaultResponseStatus()}, nil)

		mockGetDatacenter(client, datacenterID3132, getTestDatacenterUpdate(), nil, testutils.ThreeTimes)

		mockDeleteDatacenter(client)

		resourceName := "akamai_gtm_datacenter.tfexample_dc_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "nickname", "tfexample_dc_1"),
							resource.TestCheckResourceAttr(resourceName, "continent", "EU"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "nickname", "tfexample_dc_1"),
							resource.TestCheckResourceAttr(resourceName, "continent", "NA"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("update datacenter failed", func(t *testing.T) {
		client := &gtm.Mock{}

		mockCreateDatacenter(client, &gtm.CreateDatacenterResponse{
			Resource: getTestDatacenterResp(),
			Status:   getPendingResponseStatus(),
		}, nil)

		mockGetDatacenter(client, datacenterID3132, getTestDatacenterResp(), nil, testutils.FourTimes)

		mockGetDomainStatus(client, testutils.Once)

		mockUpdateDatacenter(client, nil, &gtm.Error{
			Type:       "internal_error",
			Title:      "Internal Server Error",
			Detail:     "Error updating datacenter",
			StatusCode: http.StatusInternalServerError,
		})

		mockGetDatacenter(client, datacenterID3132, getTestDatacenterResp(), nil, testutils.Once)

		mockDeleteDatacenter(client)

		resourceName := "akamai_gtm_datacenter.tfexample_dc_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "nickname", "tfexample_dc_1"),
							resource.TestCheckResourceAttr(resourceName, "continent", "EU"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/update_basic.tf"),
						ExpectError: regexp.MustCompile("API error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create datacenter, remove outside of terraform, expect non-empty plan", func(t *testing.T) {
		client := &gtm.Mock{}

		mockCreateDatacenter(client, &gtm.CreateDatacenterResponse{
			Resource: getTestDatacenterResp(),
			Status:   getPendingResponseStatus(),
		}, nil)

		mockGetDatacenter(client, datacenterID3132, getTestDatacenterResp(), nil, testutils.Twice)

		// Mock that the datacenter was deleted outside terraform
		mockGetDatacenter(client, datacenterID3132, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		// For terraform test framework, we need to mock GetDatacenter as it would actually exist before deletion
		mockGetDatacenter(client, datacenterID3132, getTestDatacenterResp(), nil, testutils.Once)

		mockDeleteDatacenter(client)

		resourceName := "akamai_gtm_datacenter.tfexample_dc_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "nickname", "tfexample_dc_1"),
							resource.TestCheckResourceAttr(resourceName, "continent", "EU"),
						),
					},
					{
						Config:             testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
						ExpectNonEmptyPlan: true,
						PlanOnly:           true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create datacenter failed", func(t *testing.T) {
		client := &gtm.Mock{}

		mockCreateDatacenter(client, nil, &gtm.Error{StatusCode: http.StatusBadRequest})

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
						ExpectError: regexp.MustCompile("Datacenter create error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create datacenter denied", func(t *testing.T) {
		client := &gtm.Mock{}

		mockCreateDatacenter(client, &gtm.CreateDatacenterResponse{
			Resource: getTestDatacenterResp(),
			Status:   getDeniedResponseStatus(),
		}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestResGTMDatacenterImport(t *testing.T) {
	tests := map[string]struct {
		domainName   string
		datacenterID string
		init         func(*gtm.Mock)
		expectError  *regexp.Regexp
		stateCheck   resource.ImportStateCheckFunc
	}{
		"happy path - import": {
			domainName:   testDomainName,
			datacenterID: "3132",
			init: func(m *gtm.Mock) {
				// Read
				mockGetDatacenter(m, datacenterID3132, getImportedDatacenter(), nil, testutils.Twice)
			},
			stateCheck: test.NewImportChecker().
				CheckEqual("domain", "gtm_terra_testdomain.akadns.net").
				CheckEqual("datacenter_id", "3132").
				CheckEqual("nickname", "testNickname").
				CheckEqual("city", "city").
				CheckEqual("wait_on_complete", "true").
				CheckEqual("clone_of", "5").
				CheckEqual("cloud_server_host_header_override", "true").
				CheckEqual("cloud_server_targeting", "true").
				CheckEqual("default_load_object.0.load_object", "loadObject").
				CheckEqual("default_load_object.0.load_object_port", "80").
				CheckEqual("default_load_object.0.load_servers.0", "1.1.1.1").
				CheckEqual("default_load_object.0.load_servers.1", "2.2.2.2").
				CheckEqual("continent", "continent").
				CheckEqual("country", "country").
				CheckEqual("latitude", "3.3").
				CheckEqual("longitude", "4.4").
				CheckEqual("score_penalty", "2").
				CheckEqual("servermonitor_pool", "serverMonitorPool").
				CheckEqual("servermonitor_liveness_count", "1").
				CheckEqual("servermonitor_load_count", "123").
				CheckEqual("ping_interval", "1").
				CheckEqual("ping_packet_size", "123").
				CheckEqual("state_or_province", "state").
				CheckEqual("virtual", "true").Build(),
		},
		"expect error - no domain name, invalid import ID": {
			domainName:   "",
			datacenterID: "3132",
			expectError:  regexp.MustCompile(`Error: Invalid Datacenter resource ID`),
		},
		"expect error - wrong datacenterID, invalid import ID": {
			domainName:   "",
			datacenterID: "wrong id",
			expectError:  regexp.MustCompile(`Error: Invalid Datacenter resource ID`),
		},
		"expect error - read": {
			domainName:   testDomainName,
			datacenterID: "3132",
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetDatacenter(m, datacenterID3132, nil, fmt.Errorf("get failed"), testutils.Once)
			},
			expectError: regexp.MustCompile(`get failed`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &gtm.Mock{}
			if tc.init != nil {
				tc.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							ImportStateCheck: tc.stateCheck,
							ImportStateId:    fmt.Sprintf("%s:%s", tc.domainName, tc.datacenterID),
							ImportState:      true,
							ResourceName:     "akamai_gtm_datacenter.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/import_basic.tf"),
							ExpectError:      tc.expectError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func mockGetDatacenter(m *gtm.Mock, datacenterID int, resp *gtm.Datacenter, err error, times int) *mock.Call {
	return m.On("GetDatacenter", testutils.MockContext, gtm.GetDatacenterRequest{
		DatacenterID: datacenterID,
		DomainName:   testDomainName,
	}).Return(resp, err).Times(times)
}

func mockUpdateDatacenter(client *gtm.Mock, resp *gtm.UpdateDatacenterResponse, err error) *mock.Call {
	return client.On("UpdateDatacenter",
		testutils.MockContext,
		gtm.UpdateDatacenterRequest{
			Datacenter: getTestDatacenterUpdate(),
			DomainName: testDomainName,
		},
	).Return(resp, err).Once()
}

func mockCreateDatacenter(client *gtm.Mock, resp *gtm.CreateDatacenterResponse, err error) *mock.Call {
	return client.On("CreateDatacenter",
		testutils.MockContext,
		gtm.CreateDatacenterRequest{
			Datacenter: getTestDatacenter(),
			DomainName: testDomainName,
		},
	).Return(resp, err).Once()
}

func mockDeleteDatacenter(client *gtm.Mock) *mock.Call {
	return client.On("DeleteDatacenter",
		testutils.MockContext,
		gtm.DeleteDatacenterRequest{
			DatacenterID: datacenterID3132,
			DomainName:   testDomainName,
		},
	).Return(&gtm.DeleteDatacenterResponse{
		Status: getDefaultResponseStatus(),
	}, nil).Once()
}

func getImportedDatacenter() *gtm.Datacenter {
	return &gtm.Datacenter{
		DatacenterID:    datacenterID3132,
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
		ServermonitorLivenessCount:    1,
		ServermonitorLoadCount:        123,
		CloudServerTargeting:          true,
		CloudServerHostHeaderOverride: true,
		PingPacketSize:                123,
		PingInterval:                  1,
	}
}

func getTestDatacenter() *gtm.Datacenter {
	return &gtm.Datacenter{
		City:                 "Snæfellsjökull",
		CloudServerTargeting: false,
		Continent:            "EU",
		Country:              "IS",
		DefaultLoadObject: &gtm.LoadObject{
			LoadObject:     "/test",
			LoadObjectPort: 80,
			LoadServers:    []string{"1.2.3.4", "1.2.3.9"},
		},
		Latitude:        64.808,
		Longitude:       -23.776,
		Nickname:        "tfexample_dc_1",
		StateOrProvince: "",
	}
}

func getTestDatacenterResp() *gtm.Datacenter {
	return &gtm.Datacenter{
		City:                 "Snæfellsjökull",
		CloudServerTargeting: false,
		Continent:            "EU",
		Country:              "IS",
		DatacenterID:         datacenterID3132,
		DefaultLoadObject: &gtm.LoadObject{
			LoadObject:     "/test",
			LoadObjectPort: 80,
			LoadServers:    []string{"1.2.3.4", "1.2.3.9"},
		},
		Latitude: 64.808,
		Links: []gtm.Link{
			{
				Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3132",
				Rel:  "self",
			},
		},
		Longitude:       -23.776,
		Nickname:        "tfexample_dc_1",
		StateOrProvince: "",
		Virtual:         true,
	}
}

func getTestDatacenterUpdate() *gtm.Datacenter {
	return &gtm.Datacenter{
		CloudServerTargeting: false,
		Continent:            "NA",
		DatacenterID:         datacenterID3132,
		DefaultLoadObject: &gtm.LoadObject{
			LoadObject:     "/test",
			LoadObjectPort: 80,
			LoadServers:    []string{"1.2.3.5", "1.2.3.6"},
		},
		Links: []gtm.Link{
			{
				Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/datacenters/3132",
				Rel:  "self",
			},
		},
		Nickname:        "tfexample_dc_1",
		StateOrProvince: "",
		Virtual:         true,
	}
}
