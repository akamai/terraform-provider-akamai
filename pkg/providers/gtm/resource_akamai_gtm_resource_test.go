package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/gtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

var rsrc = gtm.Resource{
	Name:            "tfexample_resource_1",
	AggregationType: "latest",
	Type:            "XML load object via HTTP",
	ResourceInstances: []*gtm.ResourceInstance{
		{
			DatacenterId:         3131,
			UseDefaultLoadObject: false,
			LoadObject: gtm.LoadObject{
				LoadObject:     "/test1",
				LoadServers:    []string{"1.2.3.4"},
				LoadObjectPort: 80,
			},
		},
	},
}

func TestResGtmResource(t *testing.T) {

	t.Run("create resource", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.ResourceResponse{}
		resp.Resource = &rsrc
		resp.Status = &pendingResponseStatus
		client.On("CreateResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Resource), nil}
		})

		resCall := client.On("NewResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		)

		resCall.RunFn = func(args mock.Arguments) {
			resCall.ReturnArguments = mock.Arguments{
				&gtm.Resource{
					Name: args.String(1),
				},
			}
		}

		resInstCall := client.On("NewResourceInstance",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			mock.AnythingOfType("int"),
		)

		resInstCall.RunFn = func(args mock.Arguments) {
			resInstCall.ReturnArguments = mock.Arguments{
				&gtm.ResourceInstance{
					DatacenterId: args.Int(2),
				},
			}
		}

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("UpdateResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Resource), nil}
		})

		client.On("DeleteResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_resource.tfexample_resource_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmResource/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_resource_1"),
							resource.TestCheckResourceAttr(dataSourceName, "aggregation_type", "latest"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmResource/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_resource_1"),
							resource.TestCheckResourceAttr(dataSourceName, "aggregation_type", "latest"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create resource failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("CreateResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		resCall := client.On("NewResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		)

		resCall.RunFn = func(args mock.Arguments) {
			resCall.ReturnArguments = mock.Arguments{
				&gtm.Resource{
					Name: args.String(1),
				},
			}
		}

		resInstCall := client.On("NewResourceInstance",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			mock.AnythingOfType("int"),
		)

		resInstCall.RunFn = func(args mock.Arguments) {
			resInstCall.ReturnArguments = mock.Arguments{
				&gtm.ResourceInstance{
					DatacenterId: args.Int(2),
				},
			}
		}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmResource/create_basic.tf"),
						ExpectError: regexp.MustCompile("Resource Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create resource denied", func(t *testing.T) {
		client := &gtm.Mock{}

		dr := gtm.ResourceResponse{}
		dr.Resource = &rsrc
		dr.Status = &deniedResponseStatus
		client.On("CreateResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			gtmTestDomain,
		).Return(&dr, nil)

		resCall := client.On("NewResource",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		)

		resCall.RunFn = func(args mock.Arguments) {
			resCall.ReturnArguments = mock.Arguments{
				&gtm.Resource{
					Name: args.String(1),
				},
			}
		}

		resInstCall := client.On("NewResourceInstance",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.Resource"),
			mock.AnythingOfType("int"),
		)

		resInstCall.RunFn = func(args mock.Arguments) {
			resInstCall.ReturnArguments = mock.Arguments{
				&gtm.ResourceInstance{
					DatacenterId: args.Int(2),
				},
			}
		}

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmResource/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}
