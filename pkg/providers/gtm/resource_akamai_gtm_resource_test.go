package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

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
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmResource/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_resource_1"),
							resource.TestCheckResourceAttr(dataSourceName, "aggregation_type", "latest"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmResource/update_basic.tf"),
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
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmResource/create_basic.tf"),
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
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmResource/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestGTMResourceOrder(t *testing.T) {
	tests := map[string]struct {
		client        *gtm.Mock
		pathForCreate string
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"reordered `load_servers` - no diff": {
			client:        getGTMResourceMocks(),
			pathForCreate: "testdata/TestResGtmResource/order/create.tf",
			pathForUpdate: "testdata/TestResGtmResource/order/load_servers/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered `resource_instance` - no diff": {
			client:        getGTMResourceMocks(),
			pathForCreate: "testdata/TestResGtmResource/order/create.tf",
			pathForUpdate: "testdata/TestResGtmResource/order/resource_instance/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered `resource_instance` and `load_servers` - no diff": {
			client:        getGTMResourceMocks(),
			pathForCreate: "testdata/TestResGtmResource/order/create.tf",
			pathForUpdate: "testdata/TestResGtmResource/order/reorder_resource_instance_load_servers.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"change `name` attribute - diff only for `name`": {
			client:        getGTMResourceMocks(),
			pathForCreate: "testdata/TestResGtmResource/order/create.tf",
			pathForUpdate: "testdata/TestResGtmResource/order/update_name.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reorder and change in `load_servers` - diff only for `load_servers`": {
			client:        getGTMResourceMocks(),
			pathForCreate: "testdata/TestResGtmResource/order/create.tf",
			pathForUpdate: "testdata/TestResGtmResource/order/load_servers/reorder_and_update.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reorder resource_instance and change in `load_servers` - messy diff": {
			client:        getGTMResourceMocks(),
			pathForCreate: "testdata/TestResGtmResource/order/create.tf",
			pathForUpdate: "testdata/TestResGtmResource/order/resource_instance/reorder_and_update_load_servers.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			useClient(test.client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, test.pathForCreate),
						},
						{
							Config:             testutils.LoadFixtureString(t, test.pathForUpdate),
							PlanOnly:           test.planOnly,
							ExpectNonEmptyPlan: test.nonEmptyPlan,
						},
					},
				})
			})
			test.client.AssertExpectations(t)
		})
	}
}

// getGTMResourceMocks mocks creation and deletion calls for the gtm_resource
func getGTMResourceMocks() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetResource := client.On("GetResource",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(nil, &gtm.Error{
		StatusCode: http.StatusNotFound,
	})

	resp := gtm.ResourceResponse{}
	resp.Resource = &resourceForOrderTests
	resp.Status = &pendingResponseStatus
	client.On("CreateResource",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.Resource"),
		mock.AnythingOfType("string"),
	).Return(&resp, nil).Run(func(args mock.Arguments) {
		mockGetResource.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.Resource), nil}
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

	client.On("DeleteResource",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.Resource"),
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	return client
}

var (
	// resourceForOrderTests is a gtm.Resource structure used in testing the order of resource_instance
	resourceForOrderTests = gtm.Resource{
		Name:            "tfexample_resource_1",
		AggregationType: "latest",
		Type:            "XML load object via HTTP",
		ResourceInstances: []*gtm.ResourceInstance{
			{
				DatacenterId:         3131,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test1",
					LoadServers:    []string{"1.2.3.4", "1.2.3.5", "1.2.3.6"},
					LoadObjectPort: 80,
				},
			},
			{
				DatacenterId:         3132,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test2",
					LoadServers:    []string{"1.2.3.7", "1.2.3.8", "1.2.3.9", "1.2.3.10"},
					LoadObjectPort: 80,
				},
			},
		},
	}

	rsrc = gtm.Resource{
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
)
