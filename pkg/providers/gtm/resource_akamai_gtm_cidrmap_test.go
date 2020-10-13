package gtm

import (
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"net/http"
	"regexp"
	"testing"
)

var cidr = gtm.CidrMap{
	Name: "tfexample_cidrmap_1",
	DefaultDatacenter: &gtm.DatacenterBase{
		DatacenterId: 5400,
		Nickname:     "default datacenter",
	},
	Assignments: []*gtm.CidrAssignment{
		{
			DatacenterBase: gtm.DatacenterBase{
				DatacenterId: 3131,
				Nickname:     "tfexample_dc_1",
			},
			Blocks: []string{"1.2.3.9/24"},
		},
	},
}

func TestResGtmCidrMap(t *testing.T) {

	t.Run("create cidrmap", func(t *testing.T) {
		client := &mockgtm{}

		getCall := client.On("GetCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			cidr.Name,
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.CidrMapResponse{}
		resp.Resource = &cidr
		resp.Status = &pendingResponseStatus
		client.On("CreateCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			gtmTestDomain,
		).Return(nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			gtmTestDomain,
		).Return(&completeResponseStatus, nil)

		client.On("UpdateCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			gtmTestDomain,
		).Return(&completeResponseStatus, nil)

		client.On("DeleteCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_cidrmap.tfexample_cidrmap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmCidrMap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_cidrmap_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmCidrMap/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_cidrmap_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create cidrmap failed", func(t *testing.T) {
		client := &mockgtm{}

		client.On("CreateCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
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
						Config:      loadFixtureString("testdata/TestResGtmCidrMap/create_basic.tf"),
						ExpectError: regexp.MustCompile("geoMap Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create cidrmap denied", func(t *testing.T) {
		client := &mockgtm{}

		dr := gtm.CidrMapResponse{}
		dr.Resource = &cidr
		dr.Status = &deniedResponseStatus
		client.On("CreateCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			gtmTestDomain,
		).Return(&dr, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmCidrMap/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
