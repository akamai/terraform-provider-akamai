package gtm

import (
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"net/http"
	"regexp"
	"testing"
)

var asmap = gtm.AsMap{
	Name: "tfexample_as_1",
	DefaultDatacenter: &gtm.DatacenterBase{
		DatacenterId: 5400,
		Nickname:     "default datacenter",
	},
	Assignments: []*gtm.AsAssignment{
		{
			DatacenterBase: gtm.DatacenterBase{
				DatacenterId: 3131,
				Nickname:     "tfexample_dc_1",
			},
			AsNumbers: []int64{12222, 16702, 17334},
		},
		{
			DatacenterBase: gtm.DatacenterBase{
				DatacenterId: 3132,
				Nickname:     "tfexample_dc_2",
			},
			AsNumbers: []int64{12229, 16703, 17335},
		},
	},
}

func TestResGtmAsMap(t *testing.T) {

	t.Run("create asmap", func(t *testing.T) {
		client := &mockgtm{}

		getCall := client.On("GetAsMap",
			mock.Anything, // ctx is irrelevant for this test
			asmap.Name,
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.AsMapResponse{}
		resp.Resource = &asmap
		resp.Status = &pendingResponseStatus
		client.On("CreateAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.AsMap"),
			gtmTestDomain,
		).Return(nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			gtmTestDomain,
		).Return(&completeResponseStatus, nil)

		client.On("UpdateAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.AsMap"),
			gtmTestDomain,
		).Return(&completeResponseStatus, nil)

		client.On("DeleteAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.AsMap"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_asmap.tfexample_asmap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmAsMap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_asmap_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmAsMap/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_asmap_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create asmap failed", func(t *testing.T) {
		client := &mockgtm{}

		client.On("CreateAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.AsMap"),
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
						Config:      loadFixtureString("testdata/TestResGtmAsMap/create_basic.tf"),
						ExpectError: regexp.MustCompile("asMap Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create asmap denied", func(t *testing.T) {
		client := &mockgtm{}

		dr := gtm.AsMapResponse{}
		dr.Resource = &asmap
		dr.Status = &deniedResponseStatus
		client.On("CreateAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.AsMap"),
			gtmTestDomain,
		).Return(&dr, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheck(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmAsMap/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
