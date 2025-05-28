package gtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

const testResourceName = "tfexample_resource_1"

func TestResGTMResource(t *testing.T) {

	t.Run("create resource", func(t *testing.T) {
		client := &gtm.Mock{}

		// Create
		mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateResource(client, getDefaultResource(), &gtm.CreateResourceResponse{
			Resource: getDefaultResource(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		// Read after create + refresh
		mockGetResource(client, getDefaultResource(), nil, testutils.ThreeTimes)

		// Update
		mockGetResource(client, getUpdatedResource(), nil, testutils.Once)

		mockUpdateResource(client, &gtm.UpdateResourceResponse{Status: getDefaultResponseStatus()}, nil)

		mockGetDomainStatus(client, testutils.Once)

		// Read after create + refresh
		mockGetResource(client, getUpdatedResource(), nil, testutils.ThreeTimes)

		mockDeleteResource(client)
		mockGetDomainStatus(client, testutils.Once)

		resourceName := "akamai_gtm_resource.tfexample_resource_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmResource/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "tfexample_resource_1"),
							resource.TestCheckResourceAttr(resourceName, "aggregation_type", "latest"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmResource/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "tfexample_resource_1"),
							resource.TestCheckResourceAttr(resourceName, "aggregation_type", "latest"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("update resource failed", func(t *testing.T) {
		client := &gtm.Mock{}

		// Create
		mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateResource(client, getDefaultResource(), &gtm.CreateResourceResponse{
			Resource: getDefaultResource(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		// Read after create + refresh
		mockGetResource(client, getDefaultResource(), nil, testutils.ThreeTimes)

		// Update
		mockGetResource(client, getUpdatedResource(), nil, testutils.Once)

		mockUpdateResource(client, nil, &gtm.Error{
			Type:       "internal_error",
			Title:      "Internal Server Error",
			Detail:     "Error updating resource",
			StatusCode: http.StatusInternalServerError,
		})

		// Read after create + refresh
		mockGetResource(client, getDefaultResource(), nil, testutils.Once)

		mockDeleteResource(client)
		mockGetDomainStatus(client, testutils.Once)

		resourceName := "akamai_gtm_resource.tfexample_resource_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmResource/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "tfexample_resource_1"),
							resource.TestCheckResourceAttr(resourceName, "aggregation_type", "latest"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmResource/update_basic.tf"),
						ExpectError: regexp.MustCompile("API error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create resource, remove outside of terraform, expect non-empty plan", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateResource(client, getDefaultResource(), &gtm.CreateResourceResponse{
			Resource: getDefaultResource(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetResource(client, getDefaultResource(), nil, testutils.Twice)

		// Mock that the resource was deleted outside terraform
		mockGetResource(client, nil, gtm.ErrNotFound, testutils.Once)

		// For terraform test framework, we need to mock GetResource as it would actually exist before deletion
		mockGetResource(client, getDefaultResource(), nil, testutils.Once)

		mockDeleteResource(client)

		resourceName := "akamai_gtm_resource.tfexample_resource_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmResource/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "tfexample_resource_1"),
							resource.TestCheckResourceAttr(resourceName, "aggregation_type", "latest"),
						),
					},
					{
						Config:             testutils.LoadFixtureString(t, "testdata/TestResGtmResource/create_basic.tf"),
						ExpectNonEmptyPlan: true,
						PlanOnly:           true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create resource failed", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateResource(client, getDefaultResource(), nil, &gtm.Error{StatusCode: http.StatusBadRequest})

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmResource/create_basic.tf"),
						ExpectError: regexp.MustCompile("Resource create error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create resource failed - resource already exists", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetResource(client, getDefaultResource(), nil, testutils.Once)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmResource/create_basic.tf"),
						ExpectError: regexp.MustCompile("resource already exists error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create resource denied", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateResource(client, getDefaultResource(), &gtm.CreateResourceResponse{
			Resource: getDefaultResource(),
			Status:   getDeniedResponseStatus(),
		}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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

func mockUpdateResource(client *gtm.Mock, resp *gtm.UpdateResourceResponse, err error) *mock.Call {
	return client.On("UpdateResource",
		testutils.MockContext,
		gtm.UpdateResourceRequest{
			Resource:   getUpdatedResource(),
			DomainName: testDomainName,
		},
	).Return(resp, err).Once()
}

func mockCreateResource(client *gtm.Mock, resource *gtm.Resource, resp *gtm.CreateResourceResponse, err error) *mock.Call {
	return client.On("CreateResource",
		testutils.MockContext,
		gtm.CreateResourceRequest{
			Resource:   resource,
			DomainName: testDomainName,
		},
	).Return(resp, err).Once()
}

func mockGetResource(client *gtm.Mock, resource *gtm.Resource, err error, times int) *mock.Call {
	var resp *gtm.GetResourceResponse
	if resource != nil {
		r := gtm.GetResourceResponse(*resource)
		resp = &r
	}
	return client.On("GetResource",
		testutils.MockContext,
		gtm.GetResourceRequest{
			ResourceName: testResourceName,
			DomainName:   testDomainName,
		},
	).Return(resp, err).Times(times)
}

func getDefaultResource() *gtm.Resource {
	return &gtm.Resource{
		Type: "XML load object via HTTP",
		ResourceInstances: []gtm.ResourceInstance{
			{
				DatacenterID:         datacenterID3131,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test1",
					LoadServers:    []string{"1.2.3.4"},
					LoadObjectPort: 80,
				},
			},
		},
		AggregationType: "latest",
		Name:            testResourceName,
	}
}

func getUpdatedResource() *gtm.Resource {
	return &gtm.Resource{
		Type: "XML load object via HTTP",
		ResourceInstances: []gtm.ResourceInstance{
			{
				DatacenterID:         datacenterID3132,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test2",
					LoadServers:    []string{"1.2.3.5"},
					LoadObjectPort: 80,
				},
			},
		},
		AggregationType: "latest",
		Name:            testResourceName,
	}
}

func TestGTMResourceOrder(t *testing.T) {
	tests := map[string]struct {
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"reordered `load_servers` - no diff": {
			pathForUpdate: "testdata/TestResGtmResource/order/load_servers/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered `resource_instance` - no diff": {
			pathForUpdate: "testdata/TestResGtmResource/order/resource_instance/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered `resource_instance` and `load_servers` - no diff": {
			pathForUpdate: "testdata/TestResGtmResource/order/reorder_resource_instance_load_servers.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"change `name` attribute - diff only for `name`": {
			pathForUpdate: "testdata/TestResGtmResource/order/update_name.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"reorder and change in `load_servers` - diff only for `load_servers`": {
			pathForUpdate: "testdata/TestResGtmResource/order/load_servers/reorder_and_update.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"reorder resource_instance and change in `load_servers` - messy diff": {
			pathForUpdate: "testdata/TestResGtmResource/order/resource_instance/reorder_and_update_load_servers.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			testClient := getGTMResourceMocks()
			useClient(testClient, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, "testdata/TestResGtmResource/order/create.tf"),
						},
						{
							Config:             testutils.LoadFixtureString(t, test.pathForUpdate),
							PlanOnly:           test.planOnly,
							ExpectNonEmptyPlan: test.nonEmptyPlan,
						},
					},
				})
			})
			testClient.AssertExpectations(t)
		})
	}
}

func TestResGTMResourceImport(t *testing.T) {
	tests := map[string]struct {
		domainName   string
		resourceName string
		init         func(*gtm.Mock)
		expectError  *regexp.Regexp
		stateCheck   resource.ImportStateCheckFunc
	}{
		"happy path - import": {
			domainName:   testDomainName,
			resourceName: testResourceName,
			init: func(m *gtm.Mock) {
				// Read
				mockGetResource(m, getImportedResource(), nil, testutils.Twice)
			},
			stateCheck: test.NewImportChecker().
				CheckEqual("domain", "gtm_terra_testdomain.akadns.net").
				CheckEqual("name", "tfexample_resource_1").
				CheckEqual("type", "XML load object via HTTP").
				CheckEqual("host_header", "test host").
				CheckEqual("least_squares_decay", "1").
				CheckEqual("description", "test description").
				CheckEqual("leader_string", "test string").
				CheckEqual("constrained_property", "test property").
				CheckEqual("aggregation_type", "latest").
				CheckEqual("load_imbalance_percentage", "1").
				CheckEqual("upper_bound", "5").
				CheckEqual("max_u_multiplicative_increment", "10").
				CheckEqual("decay_rate", "1").
				CheckEqual("resource_instance.0.datacenter_id", "3131").
				CheckEqual("resource_instance.0.use_default_load_object", "false").
				CheckEqual("resource_instance.0.load_object", "/test1").
				CheckEqual("resource_instance.0.load_object_port", "80").
				CheckEqual("resource_instance.0.load_servers.0", "1.2.3.4").
				CheckEqual("resource_instance.0.load_servers.1", "1.2.3.5").
				CheckEqual("resource_instance.0.load_servers.2", "1.2.3.6").
				CheckEqual("wait_on_complete", "true").Build(),
		},
		"expect error - no domain name, invalid import ID": {
			domainName:   "",
			resourceName: testResourceName,
			expectError:  regexp.MustCompile(`Error: invalid resource ID: :tfexample_resource_1`),
		},
		"expect error - no map name, invalid import ID": {
			domainName:   testDomainName,
			resourceName: "",
			expectError:  regexp.MustCompile(`Error: invalid resource ID: gtm_terra_testdomain.akadns.net:`),
		},
		"expect error - read": {
			domainName:   testDomainName,
			resourceName: testResourceName,
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetResource(m, nil, fmt.Errorf("get failed"), testutils.Once)
			},
			expectError: regexp.MustCompile(`get failed`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &gtm.Mock{}
			if tc.init != nil {
				tc.init(client)
			}
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					Steps: []resource.TestStep{
						{
							ImportStateCheck: tc.stateCheck,
							ImportStateId:    fmt.Sprintf("%s:%s", tc.domainName, tc.resourceName),
							ImportState:      true,
							ResourceName:     "akamai_gtm_resource.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResGtmResource/import_basic.tf"),
							ExpectError:      tc.expectError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

// getGTMResourceMocks mocks creation and deletion calls for the gtm_resource
func getGTMResourceMocks() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

	mockCreateResource(client, getCreatedResource(), &gtm.CreateResourceResponse{
		Resource: getCreatedResourceResp(),
		Status:   getDefaultResponseStatus(),
	}, nil)

	mockGetResource(client, getCreatedResourceResp(), nil, testutils.FourTimes)

	mockDeleteResource(client)

	return client
}

func getImportedResource() *gtm.Resource {
	return &gtm.Resource{
		Type:                "XML load object via HTTP",
		HostHeader:          "test host",
		LeastSquaresDecay:   1,
		Description:         "test description",
		LeaderString:        "test string",
		ConstrainedProperty: "test property",
		ResourceInstances: []gtm.ResourceInstance{
			{
				DatacenterID:         datacenterID3131,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test1",
					LoadServers:    []string{"1.2.3.4", "1.2.3.5", "1.2.3.6"},
					LoadObjectPort: 80,
				},
			},
		},
		AggregationType:             "latest",
		LoadImbalancePercentage:     1,
		UpperBound:                  5,
		Name:                        testResourceName,
		MaxUMultiplicativeIncrement: 10,
		DecayRate:                   1,
	}
}

func mockDeleteResource(client *gtm.Mock) *mock.Call {
	return client.On("DeleteResource",
		testutils.MockContext,
		gtm.DeleteResourceRequest{
			ResourceName: testResourceName,
			DomainName:   testDomainName,
		},
	).Return(&gtm.DeleteResourceResponse{
		Status: getDefaultResponseStatus(),
	}, nil).Once()
}

func getCreatedResource() *gtm.Resource {
	return &gtm.Resource{
		Type: "XML load object via HTTP",
		ResourceInstances: []gtm.ResourceInstance{
			{
				DatacenterID:         datacenterID3131,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test1",
					LoadServers:    []string{"1.2.3.5", "1.2.3.4", "1.2.3.6"},
					LoadObjectPort: 80,
				},
			},
			{
				DatacenterID:         datacenterID3132,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test2",
					LoadServers:    []string{"1.2.3.7", "1.2.3.10", "1.2.3.9", "1.2.3.8"},
					LoadObjectPort: 80,
				},
			},
		},
		AggregationType: "latest",
		Name:            testResourceName,
	}
}

func getCreatedResourceResp() *gtm.Resource {
	return &gtm.Resource{
		Name:            testResourceName,
		AggregationType: "latest",
		Type:            "XML load object via HTTP",
		ResourceInstances: []gtm.ResourceInstance{
			{
				DatacenterID:         datacenterID3131,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test1",
					LoadServers:    []string{"1.2.3.4", "1.2.3.5", "1.2.3.6"},
					LoadObjectPort: 80,
				},
			},
			{
				DatacenterID:         datacenterID3132,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test2",
					LoadServers:    []string{"1.2.3.7", "1.2.3.8", "1.2.3.9", "1.2.3.10"},
					LoadObjectPort: 80,
				},
			},
		},
	}
}
