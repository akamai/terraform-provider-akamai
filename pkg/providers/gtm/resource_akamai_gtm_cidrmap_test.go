package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/gtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
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

func TestResGtmCidrmap(t *testing.T) {
	dc := gtm.Datacenter{}

	t.Run("create cidrmap", func(t *testing.T) {
		client := &mockgtm{}

		getCall := client.On("GetCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.CidrMapResponse{}
		resp.Resource = &cidr
		resp.Status = &pendingResponseStatus
		client.On("CreateCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{resp.Resource, nil}
		})

		client.On("NewCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&cidr, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("UpdateCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("DeleteCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_cidrmap.tfexample_cidrmap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmCidrmap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_cidrmap_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmCidrmap/update_basic.tf"),
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

		client.On("NewCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&cidr, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmCidrmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("cidrMap Create failed"),
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

		client.On("NewCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&cidr, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmCidrmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
