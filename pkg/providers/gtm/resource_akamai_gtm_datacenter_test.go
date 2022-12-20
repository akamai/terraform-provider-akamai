package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/gtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

var dc = gtm.Datacenter{
	City:                 "Snæfellsjökull",
	CloudServerTargeting: false,
	Continent:            "EU",
	Country:              "IS",
	DatacenterId:         3132,
	DefaultLoadObject: &gtm.LoadObject{
		LoadObject:     "/test",
		LoadObjectPort: 80,
		LoadServers:    make([]string, 0),
	},
	Latitude: 64.808,
	Links: []*gtm.Link{
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

func TestResGtmDatacenter(t *testing.T) {

	t.Run("create datacenter", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.DatacenterResponse{}
		resp.Resource = &dc
		resp.Status = &pendingResponseStatus
		client.On("CreateDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Datacenter"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Datacenter), nil}
		})

		client.On("NewDatacenter",
			mock.Anything, // ctx is irrelevant for this test
		).Return(&gtm.Datacenter{})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("UpdateDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Datacenter"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Datacenter), nil}
		})

		client.On("DeleteDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Datacenter"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_datacenter.tfexample_dc_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmDatacenter/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "nickname", "tfexample_dc_1"),
							resource.TestCheckResourceAttr(dataSourceName, "continent", "EU"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmDatacenter/update_basic.tf"),
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

	t.Run("create datacenter failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("CreateDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Datacenter"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("NewDatacenter",
			mock.Anything, // ctx is irrelevant for this test
		).Return(&dc)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmDatacenter/create_basic.tf"),
						ExpectError: regexp.MustCompile("Datacenter Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create datacenter denied", func(t *testing.T) {
		client := &gtm.Mock{}

		dr := gtm.DatacenterResponse{}
		dr.Resource = &dc
		dr.Status = &deniedResponseStatus
		client.On("CreateDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Datacenter"),
			gtmTestDomain,
		).Return(&dr, nil)

		client.On("NewDatacenter",
			mock.Anything, // ctx is irrelevant for this test
		).Return(&dc)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmDatacenter/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
