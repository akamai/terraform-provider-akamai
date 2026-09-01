package gtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
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
	t.Parallel()

	t.Run("create datacenter", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockCreateDatacenter(client.GTM, testDomainName, &gtm.CreateDatacenterResponse{
			Resource: getTestDatacenterResp(),
			Status:   getPendingResponseStatus(),
		}, nil)

		mockGetDatacenter(client.GTM, testDomainName, datacenterID3132, getTestDatacenterResp(), nil, testutils.FourTimes)

		mockGetDomainStatus(client.GTM, testDomainName, testutils.Twice)

		mockUpdateDatacenter(client.GTM, testDomainName, &gtm.UpdateDatacenterResponse{Status: getDefaultResponseStatus()}, nil)

		mockGetDatacenter(client.GTM, testDomainName, datacenterID3132, getTestDatacenterUpdate(), nil, testutils.ThreeTimes)

		mockDeleteDatacenter(client.GTM, testDomainName)

		resourceName := "akamai_gtm_datacenter.tfexample_dc_1"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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

		client.GTM.AssertExpectations(t)
	})

	t.Run("update datacenter failed", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockCreateDatacenter(client.GTM, testDomainName, &gtm.CreateDatacenterResponse{
			Resource: getTestDatacenterResp(),
			Status:   getPendingResponseStatus(),
		}, nil)

		mockGetDatacenter(client.GTM, testDomainName, datacenterID3132, getTestDatacenterResp(), nil, testutils.FourTimes)

		mockGetDomainStatus(client.GTM, testDomainName, testutils.Once)

		mockUpdateDatacenter(client.GTM, testDomainName, nil, &gtm.Error{
			Type:       "internal_error",
			Title:      "Internal Server Error",
			Detail:     "Error updating datacenter",
			StatusCode: http.StatusInternalServerError,
		})

		mockGetDatacenter(client.GTM, testDomainName, datacenterID3132, getTestDatacenterResp(), nil, testutils.Once)

		mockDeleteDatacenter(client.GTM, testDomainName)

		resourceName := "akamai_gtm_datacenter.tfexample_dc_1"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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

		client.GTM.AssertExpectations(t)
	})

	t.Run("update datacenter domain name - delete and create new datacenter", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockCreateDatacenter(client.GTM, testDomainName, &gtm.CreateDatacenterResponse{
			Resource: getTestDatacenterResp(),
			Status:   getPendingResponseStatus(),
		}, nil)

		mockGetDatacenter(client.GTM, testDomainName, datacenterID3132, getTestDatacenterResp(), nil, testutils.FourTimes)

		mockDeleteDatacenter(client.GTM, testDomainName)

		domainName := "gtm_terra_testdomain_updated.akadns.net"

		mockCreateDatacenter(client.GTM, domainName, &gtm.CreateDatacenterResponse{
			Resource: getTestDatacenterResp(),
			Status:   getPendingResponseStatus(),
		}, nil)

		mockGetDatacenter(client.GTM, domainName, datacenterID3132, getTestDatacenterResp(), nil, testutils.ThreeTimes)

		mockDeleteDatacenter(client.GTM, domainName)

		resourceName := "akamai_gtm_datacenter.tfexample_dc_1"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "nickname", "tfexample_dc_1"),
						resource.TestCheckResourceAttr(resourceName, "continent", "EU"),
						resource.TestCheckResourceAttr(resourceName, "domain", "gtm_terra_testdomain.akadns.net"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/domain_update/updated_domain_name.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "nickname", "tfexample_dc_1"),
						resource.TestCheckResourceAttr(resourceName, "continent", "EU"),
						resource.TestCheckResourceAttr(resourceName, "domain", "gtm_terra_testdomain_updated.akadns.net"),
					),
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})

	t.Run("create datacenter, remove outside of terraform, expect non-empty plan", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockCreateDatacenter(client.GTM, testDomainName, &gtm.CreateDatacenterResponse{
			Resource: getTestDatacenterResp(),
			Status:   getPendingResponseStatus(),
		}, nil)

		mockGetDatacenter(client.GTM, testDomainName, datacenterID3132, getTestDatacenterResp(), nil, testutils.Twice)

		// Mock that the datacenter was deleted outside terraform
		mockGetDatacenter(client.GTM, testDomainName, datacenterID3132, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		// For terraform test framework, we need to mock GetDatacenter as it would actually exist before deletion
		mockGetDatacenter(client.GTM, testDomainName, datacenterID3132, getTestDatacenterResp(), nil, testutils.Once)

		mockDeleteDatacenter(client.GTM, testDomainName)

		resourceName := "akamai_gtm_datacenter.tfexample_dc_1"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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

		client.GTM.AssertExpectations(t)
	})

	t.Run("create datacenter failed", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockCreateDatacenter(client.GTM, testDomainName, nil, &gtm.Error{StatusCode: http.StatusBadRequest})

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
					ExpectError: regexp.MustCompile("Datacenter create error"),
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})

	t.Run("create datacenter denied", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockCreateDatacenter(client.GTM, testDomainName, &gtm.CreateDatacenterResponse{
			Resource: getTestDatacenterResp(),
			Status:   getDeniedResponseStatus(),
		}, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
					ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})
}

func TestResGTMDatacenterImport(t *testing.T) {
	t.Parallel()
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
				mockGetDatacenter(m, testDomainName, datacenterID3132, getImportedDatacenter(), nil, testutils.Twice)
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
			expectError:  regexp.MustCompile(`Error: invalid datacenter resource ID`),
		},
		"expect error - wrong datacenterID, invalid import ID": {
			domainName:   "",
			datacenterID: "wrong id",
			expectError:  regexp.MustCompile(`Error: invalid datacenter resource ID`),
		},
		"expect error - read": {
			domainName:   testDomainName,
			datacenterID: "3132",
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetDatacenter(m, testDomainName, datacenterID3132, nil, fmt.Errorf("get failed"), testutils.Once)
			},
			expectError: regexp.MustCompile(`get failed`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.GTM)
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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
			client.GTM.AssertExpectations(t)
		})
	}
}

func mockGetDatacenter(m *gtm.Mock, domainName string, datacenterID int, resp *gtm.Datacenter, err error, times int) *mock.Call {
	return m.On("GetDatacenter", testutils.MockContext, gtm.GetDatacenterRequest{
		DatacenterID: datacenterID,
		DomainName:   domainName,
	}).Return(resp, err).Times(times)
}

func mockUpdateDatacenter(client *gtm.Mock, domainName string, resp *gtm.UpdateDatacenterResponse, err error) *mock.Call {
	return client.On("UpdateDatacenter",
		testutils.MockContext,
		gtm.UpdateDatacenterRequest{
			Datacenter: getTestDatacenterUpdate(),
			DomainName: domainName,
		},
	).Return(resp, err).Once()
}

func mockCreateDatacenter(client *gtm.Mock, domainName string, resp *gtm.CreateDatacenterResponse, err error) *mock.Call {
	return client.On("CreateDatacenter",
		testutils.MockContext,
		gtm.CreateDatacenterRequest{
			Datacenter: getTestDatacenter(),
			DomainName: domainName,
		},
	).Return(resp, err).Once()
}

func mockDeleteDatacenter(client *gtm.Mock, domainName string) *mock.Call {
	return client.On("DeleteDatacenter",
		testutils.MockContext,
		gtm.DeleteDatacenterRequest{
			DatacenterID: datacenterID3132,
			DomainName:   domainName,
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
