package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/gtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

var geo = gtm.GeoMap{
	Name: "tfexample_geomap_1",
	DefaultDatacenter: &gtm.DatacenterBase{
		DatacenterId: 5400,
		Nickname:     "default datacenter",
	},
	Assignments: []*gtm.GeoAssignment{
		{
			DatacenterBase: gtm.DatacenterBase{
				DatacenterId: 3131,
				Nickname:     "tfexample_dc_1",
			},
			Countries: []string{"GB"},
		},
	},
}

func TestResGtmGeomap(t *testing.T) {
	dc := gtm.Datacenter{
		DatacenterId: geo.DefaultDatacenter.DatacenterId,
		Nickname:     geo.DefaultDatacenter.Nickname,
	}

	t.Run("create geomap", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.GeoMapResponse{}
		resp.Resource = &geo
		resp.Status = &pendingResponseStatus
		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.GeoMap), nil}
		})

		client.On("NewGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&geo, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("UpdateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.GeoMap), nil}
		})

		client.On("DeleteGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_geomap.tfexample_geomap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmGeomap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_geomap_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmGeomap/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_geomap_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create geomap failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("NewGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&geo, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmGeomap/create_basic.tf"),
						ExpectError: regexp.MustCompile("geoMap Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create geomap denied", func(t *testing.T) {
		client := &gtm.Mock{}

		dr := gtm.GeoMapResponse{}
		dr.Resource = &geo
		dr.Status = &deniedResponseStatus
		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			gtmTestDomain,
		).Return(&dr, nil)

		client.On("NewGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&geo, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmGeomap/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
