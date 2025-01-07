package gtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

const testDomainName = "gtm_terra_testdomain.akadns.net"

func TestResGTMResource(t *testing.T) {

	t.Run("create resource", func(t *testing.T) {
		client := &gtm.Mock{}

		// Create
		mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}).Once()

		resp := rsrc
		mockCreateResource(client, getDefaultResource(), &gtm.CreateResourceResponse{
			Resource: rsrcCreate.Resource,
			Status:   rsrcCreate.Status,
		}, nil)

		// Read after create + refresh
		mockGetResource(client, &resp, nil).Times(3)

		// Update
		updatedResource := gtm.GetResourceResponse(*getUpdatedResource())
		mockGetResource(client, &updatedResource, nil).Once()

		mockUpdateResource(client)

		mockGetDomainStatus(client, 1)

		// Read after create + refresh
		mockGetResource(client, &updatedResource, nil).Times(3)

		mockDeleteResource(client)
		mockGetDomainStatus(client, 1)

		dataSourceName := "akamai_gtm_resource.tfexample_resource_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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

	t.Run("create resource, remove outside of terraform, expect non-empty plan", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}).Once()

		resp := rsrc
		mockCreateResource(client, getDefaultResource(), &gtm.CreateResourceResponse{
			Resource: rsrcCreate.Resource,
			Status:   rsrcCreate.Status,
		}, nil)

		mockGetResource(client, &resp, nil).Twice()

		// Mock that the resource was deleted outside terraform
		mockGetResource(client, nil, gtm.ErrNotFound).Once()

		// For terraform test framework, we need to mock GetResource as it would actually exist before deletion
		mockGetResource(client, &resp, nil).Once()

		mockDeleteResource(client)

		dataSourceName := "akamai_gtm_resource.tfexample_resource_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmResource/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_resource_1"),
							resource.TestCheckResourceAttr(dataSourceName, "aggregation_type", "latest"),
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

		mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}).Once()

		mockCreateResource(client, getDefaultResource(), nil, &gtm.Error{StatusCode: http.StatusBadRequest})

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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

	t.Run("create resource failed - resource already exists", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetResource(client, &rsrc, nil).Once()

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

		mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}).Once()

		mockCreateResource(client, getDefaultResource(), &gtm.CreateResourceResponse{
			Resource: rsrcCreate.Resource,
			Status:   &deniedResponseStatus,
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

func mockUpdateResource(client *gtm.Mock) *mock.Call {
	return client.On("UpdateResource",
		testutils.MockContext,
		gtm.UpdateResourceRequest{
			Resource:   getUpdatedResource(),
			DomainName: testDomainName,
		},
	).Return(updateResourceResponseStatus, nil).Once()
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

func mockGetResource(client *gtm.Mock, resp *gtm.GetResourceResponse, err error) *mock.Call {
	return client.On("GetResource",
		testutils.MockContext,
		gtm.GetResourceRequest{
			ResourceName: "tfexample_resource_1",
			DomainName:   testDomainName,
		},
	).Return(resp, err)
}

func getDefaultResource() *gtm.Resource {
	return &gtm.Resource{
		Type: "XML load object via HTTP",
		ResourceInstances: []gtm.ResourceInstance{
			{
				DatacenterID:         3131,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test1",
					LoadServers:    []string{"1.2.3.4"},
					LoadObjectPort: 80,
				},
			},
		},
		AggregationType: "latest",
		Name:            "tfexample_resource_1",
	}
}

func getUpdatedResource() *gtm.Resource {
	return &gtm.Resource{
		Type: "XML load object via HTTP",
		ResourceInstances: []gtm.ResourceInstance{
			{
				DatacenterID:         3132,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test2",
					LoadServers:    []string{"1.2.3.5"},
					LoadObjectPort: 80,
				},
			},
		},
		AggregationType: "latest",
		Name:            "tfexample_resource_1",
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
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reorder and change in `load_servers` - diff only for `load_servers`": {
			pathForUpdate: "testdata/TestResGtmResource/order/load_servers/reorder_and_update.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reorder resource_instance and change in `load_servers` - messy diff": {
			pathForUpdate: "testdata/TestResGtmResource/order/resource_instance/reorder_and_update_load_servers.tf",
			nonEmptyPlan:  true, // change to false to see diff
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
			domainName:   "gtm_terra_testdomain.akadns.net",
			resourceName: "tfexample_resource_1",
			init: func(m *gtm.Mock) {
				// Read
				importedResource := gtm.GetResourceResponse(*getImportedResource())
				mockGetResource(m, &importedResource, nil).Times(2)
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
			resourceName: "tfexample_resource_1",
			expectError:  regexp.MustCompile(`Error: invalid resource ID: :tfexample_resource_1`),
		},
		"expect error - no map name, invalid import ID": {
			domainName:   "gtm_terra_testdomain.akadns.net",
			resourceName: "",
			expectError:  regexp.MustCompile(`Error: invalid resource ID: gtm_terra_testdomain.akadns.net:`),
		},
		"expect error - read": {
			domainName:   "gtm_terra_testdomain.akadns.net",
			resourceName: "tfexample_resource_1",
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetResource(m, nil, fmt.Errorf("get failed")).Times(1)
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

	mockGetResource(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}).Once()

	resourceToCreate := gtm.Resource{
		Type: "XML load object via HTTP",
		ResourceInstances: []gtm.ResourceInstance{
			{
				DatacenterID:         3131,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test1",
					LoadServers:    []string{"1.2.3.5", "1.2.3.4", "1.2.3.6"},
					LoadObjectPort: 80,
				},
			},
			{
				DatacenterID:         3132,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test2",
					LoadServers:    []string{"1.2.3.7", "1.2.3.10", "1.2.3.9", "1.2.3.8"},
					LoadObjectPort: 80,
				},
			},
		},
		AggregationType: "latest",
		Name:            "tfexample_resource_1",
	}
	mockCreateResource(client, &resourceToCreate, &gtm.CreateResourceResponse{
		Resource: &resourceToCreate,
		Status:   testStatus,
	}, nil)

	resp := resourceForOrderTests
	mockGetResource(client, &resp, nil).Times(4)

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
				DatacenterID:         3131,
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
		Name:                        "tfexample_resource_1",
		MaxUMultiplicativeIncrement: 10,
		DecayRate:                   1,
	}
}

func mockDeleteResource(client *gtm.Mock) *mock.Call {
	return client.On("DeleteResource",
		testutils.MockContext,
		gtm.DeleteResourceRequest{
			ResourceName: "tfexample_resource_1",
			DomainName:   testDomainName,
		},
	).Return(deleteResourceResponseStatus, nil)
}

var (
	// resourceForOrderTests is a gtm.Resource structure used in testing the order of resource_instance
	resourceForOrderTests = gtm.GetResourceResponse{
		Name:            "tfexample_resource_1",
		AggregationType: "latest",
		Type:            "XML load object via HTTP",
		ResourceInstances: []gtm.ResourceInstance{
			{
				DatacenterID:         3131,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test1",
					LoadServers:    []string{"1.2.3.4", "1.2.3.5", "1.2.3.6"},
					LoadObjectPort: 80,
				},
			},
			{
				DatacenterID:         3132,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test2",
					LoadServers:    []string{"1.2.3.7", "1.2.3.8", "1.2.3.9", "1.2.3.10"},
					LoadObjectPort: 80,
				},
			},
		},
	}

	rsrcCreate = gtm.CreateResourceResponse{
		Resource: &gtm.Resource{
			Name:            "tfexample_resource_1",
			AggregationType: "latest",
			Type:            "XML load object via HTTP",
			ResourceInstances: []gtm.ResourceInstance{
				{
					DatacenterID:         3131,
					UseDefaultLoadObject: false,
					LoadObject: gtm.LoadObject{
						LoadObject:     "/test1",
						LoadServers:    []string{"1.2.3.4"},
						LoadObjectPort: 80,
					},
				},
			},
		},
		Status: &gtm.ResponseStatus{
			ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
			Links: []gtm.Link{
				{
					Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
					Rel:  "self",
				},
			},
			Message:               "Current configuration has been propagated to all GTM nameservers",
			PassingValidation:     true,
			PropagationStatus:     "COMPLETE",
			PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
		},
	}

	rsrc = gtm.GetResourceResponse{
		Name:            "tfexample_resource_1",
		AggregationType: "latest",
		Type:            "XML load object via HTTP",
		ResourceInstances: []gtm.ResourceInstance{
			{
				DatacenterID:         3131,
				UseDefaultLoadObject: false,
				LoadObject: gtm.LoadObject{
					LoadObject:     "/test1",
					LoadServers:    []string{"1.2.3.4"},
					LoadObjectPort: 80,
				},
			},
		},
	}

	updateResourceResponseStatus = &gtm.UpdateResourceResponse{
		Status: &gtm.ResponseStatus{
			ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
			Links: []gtm.Link{
				{
					Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
					Rel:  "self",
				},
			},
			Message:               "Current configuration has been propagated to all GTM nameservers",
			PassingValidation:     true,
			PropagationStatus:     "COMPLETE",
			PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
		},
	}
	deleteResourceResponseStatus = &gtm.DeleteResourceResponse{
		Status: &gtm.ResponseStatus{
			ChangeID: "40e36abd-bfb2-4635-9fca-62175cf17007",
			Links: []gtm.Link{
				{
					Href: "https://akab-ymtebc45gco3ypzj-apz4yxpek55y7fyv.luna.akamaiapis.net/config-gtm/v1/domains/gtmdomtest.akadns.net/status/current",
					Rel:  "self",
				},
			},
			Message:               "Current configuration has been propagated to all GTM nameservers",
			PassingValidation:     true,
			PropagationStatus:     "COMPLETE",
			PropagationStatusDate: "2019-04-25T14:54:00.000+00:00",
		},
	}
)
