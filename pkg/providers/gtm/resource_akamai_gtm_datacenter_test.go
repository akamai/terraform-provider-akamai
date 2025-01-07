package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
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
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := dc
		client.On("CreateDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.CreateDatacenterRequest"),
		).Return(&gtm.CreateDatacenterResponse{
			Resource: &dc,
			Status:   &pendingResponseStatus,
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDomainStatusRequest"),
		).Return(getDomainStatusResponseStatus, nil)

		client.On("UpdateDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.UpdateDatacenterRequest"),
		).Return(updateDatacenterResponseStatus, nil)

		client.On("DeleteDatacenter",
			mock.Anything, // ctx is irrelevant for this test
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
			mock.Anything,
			mock.AnythingOfType("gtm.CreateDatacenterRequest"),
		).Return(&gtm.CreateDatacenterResponse{
			Resource: &dc,
			Status:   &pendingResponseStatus,
		}, nil).Once()

		client.On("GetDatacenter",
			mock.Anything,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&resp, nil).Twice()

		// Mock that the datacenter was deleted outside terraform
		client.On("GetDatacenter",
			mock.Anything,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		// For terraform test framework, we need to mock GetDatacenter as it would actually exist before deletion
		client.On("GetDatacenter",
			mock.Anything,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&resp, nil).Once()

		client.On("DeleteDatacenter",
			mock.Anything,
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
			mock.Anything, // ctx is irrelevant for this test
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
						ExpectError: regexp.MustCompile("Datacenter Create failed"),
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
			mock.Anything, // ctx is irrelevant for this test
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
