package gtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/gtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
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

func TestResGtmAsmap(t *testing.T) {
	dc := gtm.Datacenter{
		DatacenterId: asmap.DefaultDatacenter.DatacenterId,
		Nickname:     asmap.DefaultDatacenter.Nickname,
	}

	t.Run("create asmap", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.AsMapResponse{}
		resp.Resource = &asmap
		resp.Status = &pendingResponseStatus

		client.On("NewAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&asmap, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("CreateAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.AsMap"),
			mock.AnythingOfType("string"),
		).Return(&gtm.AsMapResponse{
			Resource: &asmap,
			Status:   &gtm.ResponseStatus{},
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{resp.Resource, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("UpdateAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.AsMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("DeleteAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.AsMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_asmap.tfexample_as_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmAsmap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_as_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmAsmap/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_as_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create asmap failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("CreateAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.AsMap"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("NewAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&asmap, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmAsmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("asMap Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create asmap denied", func(t *testing.T) {
		client := &gtm.Mock{}

		dr := gtm.AsMapResponse{}
		dr.Resource = &asmap
		dr.Status = &deniedResponseStatus
		client.On("CreateAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.AsMap"),
			gtmTestDomain,
		).Return(&dr, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("NewAsMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&asmap, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmAsmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("import asmap", func(t *testing.T) {
		client := &gtm.Mock{}

		resp := gtm.AsMapResponse{}
		resp.Resource = &asmap
		resp.Status = &pendingResponseStatus

		client.On("NewAsMap",
			mock.Anything,
			mock.AnythingOfType("string"),
		).Return(&asmap, nil)

		client.On("GetDatacenter",
			mock.Anything,
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("CreateAsMap",
			mock.Anything,
			mock.AnythingOfType("*gtm.AsMap"),
			mock.AnythingOfType("string"),
		).Return(&gtm.AsMapResponse{
			Resource: &asmap,
			Status:   &gtm.ResponseStatus{},
		}, nil)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("GetAsMap",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&asmap, nil)

		client.On("DeleteAsMap",
			mock.Anything,
			mock.AnythingOfType("*gtm.AsMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_asmap.tfexample_as_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmAsmap/import_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_as_1"),
						),
					},
					{
						Config:            loadFixtureString("testdata/TestResGtmAsmap/create_basic.tf"),
						ImportState:       true,
						ImportStateVerify: true,
						ImportStateId:     fmt.Sprintf("%s:%s", gtmTestDomain, "tfexample_as_1"),
						ResourceName:      dataSourceName,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
