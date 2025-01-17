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

func TestResGTMCIDRMap(t *testing.T) {
	dc := gtm.Datacenter{}

	t.Run("create CIDRMap", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetCIDRMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Twice()

		resp := cidr
		client.On("CreateCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.CreateCIDRMapRequest"),
		).Return(&gtm.CreateCIDRMapResponse{
			Resource: cidrCreate.Resource,
			Status:   cidrCreate.Status,
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		})

		client.On("GetCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetCIDRMapRequest"),
		).Return(&resp, nil).Times(3)

		client.On("GetDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&dc, nil)

		client.On("GetDomainStatus",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetDomainStatusRequest"),
		).Return(getDomainStatusResponseStatus, nil)

		client.On("UpdateCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.UpdateCIDRMapRequest"),
		).Return(updateCIDRMapResponseStatus, nil)

		client.On("GetCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetCIDRMapRequest"),
		).Return(&cidrUpdate, nil).Times(3)

		client.On("DeleteCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.DeleteCIDRMapRequest"),
		).Return(deleteCIDRMapResponseStatus, nil)

		dataSourceName := "akamai_gtm_cidrmap.tfexample_cidrmap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_cidrmap_1"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_cidrmap_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create cidrmap, remove outside of terraform, expect non-empty plan", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetCIDRMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		resp := cidr
		client.On("CreateCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.CreateCIDRMapRequest"),
		).Return(&gtm.CreateCIDRMapResponse{
			Resource: cidrCreate.Resource,
			Status:   cidrCreate.Status,
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		}).Once()

		client.On("GetCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetCIDRMapRequest"),
		).Return(&resp, nil).Twice()

		client.On("GetDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&dc, nil).Once()

		// Mock that the CIDRMap was deleted outside terraform
		client.On("GetCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetCIDRMapRequest"),
		).Return(nil, gtm.ErrNotFound).Once()

		// For terraform test framework, we need to mock GetCIDRMap as it would actually exist before deletion
		client.On("GetCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetCIDRMapRequest"),
		).Return(&resp, nil).Once()

		client.On("DeleteCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.DeleteCIDRMapRequest"),
		).Return(deleteCIDRMapResponseStatus, nil).Once()

		dataSourceName := "akamai_gtm_cidrmap.tfexample_cidrmap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_cidrmap_1"),
						),
					},
					{
						Config:             testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
						ExpectNonEmptyPlan: true,
						PlanOnly:           true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create cidrmap failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("GetCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetCIDRMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		client.On("CreateCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.CreateCIDRMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("GetDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("cidrMap create error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create cidrmap failed - cidrmap already exists", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("GetCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetCIDRMapRequest"),
		).Return(&cidr, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("cidrMap already exists error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create cidrmap denied", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("GetCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetCIDRMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		dr := gtm.CreateCIDRMapResponse{}
		dr.Resource = cidrCreate.Resource
		dr.Status = &deniedResponseStatus
		client.On("CreateCIDRMap",
			testutils.MockContext,
			mock.AnythingOfType("gtm.CreateCIDRMapRequest"),
		).Return(&dr, nil)

		client.On("GetDatacenter",
			testutils.MockContext,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestGTMCIDRMapOrder(t *testing.T) {
	tests := map[string]struct {
		client        *gtm.Mock
		pathForCreate string
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"reordered blocks - no diff": {
			client:        getCIDRMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/blocks/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered assignments - no diff": {
			client:        getCIDRMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/assignments/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered assignments and blocks - no diff": {
			client:        getCIDRMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/reorder_assignments_and_blocks.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"change to `name` attribute with different order of assignments and blocks - diff only for `name`": {
			client:        getCIDRMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/update_name.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"change to `domain` attribute with different order of assignments and blocks - diff only for `domain`": {
			client:        getCIDRMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/update_domain.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"change to `wait_on_complete` attribute with different order of assignments and blocks - diff only for `wait_on_complete`": {
			client:        getCIDRMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/update_wait_on_complete.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reordered and updated blocks - diff only for updated block": {
			client:        getCIDRMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/blocks/reorder_and_update.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reordered assignments and updated block - messy diff": {
			client:        getCIDRMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/assignments/reorder_and_update_block.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reordered assignments and updated nickname - messy diff": {
			client:        getCIDRMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/assignments/reorder_and_update_nickname.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			useClient(test.client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
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

func TestResGTMCIDRMapImport(t *testing.T) {
	tests := map[string]struct {
		domainName  string
		mapName     string
		init        func(*gtm.Mock)
		expectError *regexp.Regexp
		stateCheck  resource.ImportStateCheckFunc
	}{
		"happy path - import": {
			domainName: "gtm_terra_testdomain.akadns.net",
			mapName:    "tfexample_cidrmap_1",
			init: func(m *gtm.Mock) {
				// Read
				importedCidrMap := gtm.GetCIDRMapResponse(*getImportedCIDRMap())
				mockGetCIDRMap(m, &importedCidrMap, nil).Times(2)
			},
			stateCheck: test.NewImportChecker().
				CheckEqual("domain", "gtm_terra_testdomain.akadns.net").
				CheckEqual("name", "tfexample_cidrmap_1").
				CheckEqual("default_datacenter.0.datacenter_id", "5400").
				CheckEqual("default_datacenter.0.nickname", "default datacenter").
				CheckEqual("assignment.0.datacenter_id", "3131").
				CheckEqual("assignment.0.nickname", "tfexample_dc_1").
				CheckEqual("assignment.0.blocks.0", "1.2.3.9/24").
				CheckEqual("wait_on_complete", "true").Build(),
		},
		"expect error - no domain name, invalid import ID": {
			domainName:  "",
			mapName:     "tfexample_cidrmap_1",
			expectError: regexp.MustCompile(`Error: invalid resource ID: :tfexample_cidrmap_1`),
		},
		"expect error - no map name, invalid import ID": {
			domainName:  "gtm_terra_testdomain.akadns.net",
			mapName:     "",
			expectError: regexp.MustCompile(`Error: invalid resource ID: gtm_terra_testdomain.akadns.net:`),
		},
		"expect error - read": {
			domainName: "gtm_terra_testdomain.akadns.net",
			mapName:    "tfexample_cidrmap_1",
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetCIDRMap(m, nil, fmt.Errorf("get failed")).Once()
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
							ImportStateId:    fmt.Sprintf("%s:%s", tc.domainName, tc.mapName),
							ImportState:      true,
							ResourceName:     "akamai_gtm_cidrmap.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/import_basic.tf"),
							ExpectError:      tc.expectError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

// getCIDRMapMocks mocks creation and deletion of a resource
func getCIDRMapMocks() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetCIDRMap := client.On("GetCIDRMap",
		testutils.MockContext,
		mock.AnythingOfType("gtm.GetCIDRMapRequest"),
	).Return(nil, &gtm.Error{
		StatusCode: http.StatusNotFound,
	})

	resp := cidrMapDiffOrder
	client.On("CreateCIDRMap",
		testutils.MockContext,
		mock.AnythingOfType("gtm.CreateCIDRMapRequest"),
	).Return(&gtm.CreateCIDRMapResponse{
		Resource: cidrCreate.Resource,
		Status:   cidrCreate.Status,
	}, nil).Run(func(args mock.Arguments) {
		mockGetCIDRMap.ReturnArguments = mock.Arguments{&resp, nil}
	})

	client.On("GetDatacenter",
		testutils.MockContext,
		mock.AnythingOfType("gtm.GetDatacenterRequest"),
	).Return(&dc, nil)

	client.On("GetDomainStatus",
		testutils.MockContext,
		mock.AnythingOfType("gtm.GetDomainStatusRequest"),
	).Return(getDomainStatusResponseStatus, nil)

	client.On("DeleteCIDRMap",
		testutils.MockContext,
		mock.AnythingOfType("gtm.DeleteCIDRMapRequest"),
	).Return(deleteCIDRMapResponseStatus, nil)

	return client
}

func mockGetCIDRMap(m *gtm.Mock, resp *gtm.GetCIDRMapResponse, err error) *mock.Call {
	return m.On("GetCIDRMap", testutils.MockContext, gtm.GetCIDRMapRequest{
		MapName:    "tfexample_cidrmap_1",
		DomainName: "gtm_terra_testdomain.akadns.net",
	}).Return(resp, err)
}

func getImportedCIDRMap() *gtm.CIDRMap {
	return &gtm.CIDRMap{
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.CIDRAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					Nickname:     "tfexample_dc_1",
					DatacenterID: 3131,
				},
				Blocks: []string{"1.2.3.9/24"},
			},
		},
		Name: "tfexample_cidrmap_1",
	}
}

var (
	// cidrMapDiffOrder is a gtm.CidrMap structure used in tests of order of assignments and block in gtm_cidrmap resource
	cidrMapDiffOrder = gtm.GetCIDRMapResponse{
		Name: "tfexample_cidrmap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.CIDRAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3131,
					Nickname:     "tfexample_dc_1",
				},
				Blocks: []string{"1.2.3.4/24", "1.2.3.5/24"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3132,
					Nickname:     "tfexample_dc_2",
				},
				Blocks: []string{"1.2.3.6/24", "1.2.3.7/24", "1.2.3.8/24"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3133,
					Nickname:     "tfexample_dc_3",
				},
				Blocks: []string{"1.2.3.9/24", "1.2.3.10/24"},
			},
		},
	}

	cidrCreate = gtm.CreateCIDRMapResponse{
		Resource: &gtm.CIDRMap{
			Name: "tfexample_cidrmap_1",
			DefaultDatacenter: &gtm.DatacenterBase{
				DatacenterID: 5400,
				Nickname:     "default datacenter",
			},
			Assignments: []gtm.CIDRAssignment{
				{
					DatacenterBase: gtm.DatacenterBase{
						DatacenterID: 3131,
						Nickname:     "tfexample_dc_1",
					},
					Blocks: []string{"1.2.3.9/24"},
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

	cidr = gtm.GetCIDRMapResponse{
		Name: "tfexample_cidrmap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.CIDRAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3131,
					Nickname:     "tfexample_dc_1",
				},
				Blocks: []string{"1.2.3.9/24"},
			},
		},
	}

	cidrUpdate = gtm.GetCIDRMapResponse{
		Name: "tfexample_cidrmap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.CIDRAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3132,
					Nickname:     "tfexample_dc_2",
				},
				Blocks: []string{"1.2.3.9/16"},
			},
		},
	}

	updateCIDRMapResponseStatus = &gtm.UpdateCIDRMapResponse{
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

	deleteCIDRMapResponseStatus = &gtm.DeleteCIDRMapResponse{
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
