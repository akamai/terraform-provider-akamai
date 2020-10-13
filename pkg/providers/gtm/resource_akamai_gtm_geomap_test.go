package gtm

import (
	"net/http"
	"regexp"
	"testing"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
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

	t.Run("create geomap", func(t *testing.T) {
		client := &mockgtm{}

		getCall := client.On("GetGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			geo.Name,
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.GeoMapResponse{}
		resp.Resource = &geo
		resp.Status = &pendingResponseStatus
		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			gtmTestDomain,
		).Return(nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			gtmTestDomain,
		).Return(&completeResponseStatus, nil)

		client.On("UpdateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			gtmTestDomain,
		).Return(&completeResponseStatus, nil)

		client.On("DeleteGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_geomap.tfexample_geomap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
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
		client := &mockgtm{}

		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
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
		client := &mockgtm{}

		dr := gtm.GeoMapResponse{}
		dr.Resource = &geo
		dr.Status = &deniedResponseStatus
		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			gtmTestDomain,
		).Return(&dr, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
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
