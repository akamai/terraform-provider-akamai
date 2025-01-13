package gtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

var dc = gtm.Datacenter{
	City:                 "Snæfellsjökull",
	CloudServerTargeting: false,
	Continent:            "EU",
	Country:              "IS",
	DatacenterID:         3132,
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

func TestResGTMDatacenter(t *testing.T) {

	t.Run("create datacenter", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := dc
		client.On("CreateDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.CreateDatacenterRequest"),
		).Return(&gtm.CreateDatacenterResponse{
			Resource: &dc,
			Status:   &pendingResponseStatus,
		}, nil).Run(func(_ mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		})

		client.On("GetDomainStatus",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetDomainStatusRequest"),
		).Return(getDomainStatusResponseStatus, nil)

		client.On("UpdateDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.UpdateDatacenterRequest"),
		).Return(updateDatacenterResponseStatus, nil)

		client.On("DeleteDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.DeleteDatacenterRequest"),
		).Return(deleteDatacenterResponseStatus, nil)

		dataSourceName := "akamai_gtm_datacenter.tfexample_dc_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "nickname", "tfexample_dc_1"),
							resource.TestCheckResourceAttr(dataSourceName, "continent", "EU"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "nickname", "tfexample_dc_1"),
							resource.TestCheckResourceAttr(dataSourceName, "continent", "NA"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create datacenter, remove outside of terraform, expect non-empty plan", func(t *testing.T) {
		client := &gtm.Mock{}

		resp := dc
		client.On("CreateDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.CreateDatacenterRequest"),
		).Return(&gtm.CreateDatacenterResponse{
			Resource: &dc,
			Status:   &pendingResponseStatus,
		}, nil).Once()

		client.On("GetDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&resp, nil).Twice()

		// Mock that the datacenter was deleted outside terraform
		client.On("GetDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		// For terraform test framework, we need to mock GetDatacenter as it would actually exist before deletion
		client.On("GetDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&resp, nil).Once()

		client.On("DeleteDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.DeleteDatacenterRequest"),
		).Return(deleteDatacenterResponseStatus, nil).Once()

		dataSourceName := "akamai_gtm_datacenter.tfexample_dc_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmDatacenter/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "nickname", "tfexample_dc_1"),
							resource.TestCheckResourceAttr(dataSourceName, "continent", "EU"),
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

		client.On("CreateDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.CreateDatacenterRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

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

		dr := gtm.CreateDatacenterResponse{}
		dr.Resource = &dc
		dr.Status = &deniedResponseStatus
		client.On("CreateDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.CreateDatacenterRequest"),
		).Return(&dr, nil)

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
			domainName:   "gtm_terra_testdomain.akadns.net",
			datacenterID: "3132",
			init: func(m *gtm.Mock) {
				// Read
				mockGetDatacenterImport(m, getImportedDatacenter(), nil).Times(2)
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
			domainName:   "gtm_terra_testdomain.akadns.net",
			datacenterID: "3132",
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetDatacenterImport(m, nil, fmt.Errorf("get failed")).Once()
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

func mockGetDatacenterImport(m *gtm.Mock, resp *gtm.Datacenter, err error) *mock.Call {
	return m.On("GetDatacenter", testutils.MockContext, gtm.GetDatacenterRequest{
		DatacenterID: 3132,
		DomainName:   "gtm_terra_testdomain.akadns.net",
	}).Return(resp, err)
}

func getImportedDatacenter() *gtm.Datacenter {
	return &gtm.Datacenter{
		DatacenterID:    3132,
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

var (
	updateDatacenterResponseStatus = &gtm.UpdateDatacenterResponse{
		Status: &gtm.ResponseStatus{
			ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
			Links: []gtm.Link{
				{
					Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
					Rel:  "self",
				},
			},
			Message:               "Current configuration has been propagated to all GTM nameservers",
			PassingValidation:     true,
			PropagationStatus:     "COMPLETE",
			PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
		},
	}

	deleteDatacenterResponseStatus = &gtm.DeleteDatacenterResponse{
		Status: &gtm.ResponseStatus{
			ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
			Links: []gtm.Link{
				{
					Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
					Rel:  "self",
				},
			},
			Message:               "Current configuration has been propagated to all GTM nameservers",
			PassingValidation:     true,
			PropagationStatus:     "COMPLETE",
			PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
		},
	}
)
