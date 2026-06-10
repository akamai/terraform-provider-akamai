package gtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

const testCIDRMapName = "tfexample_cidrmap_1"

func TestResGTMCIDRMap(t *testing.T) {
	t.Parallel()
	dc := gtm.Datacenter{}

	t.Run("create CIDRMap", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockGetCIDRMap(client.GTM, testDomainName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateCIDRMap(client.GTM, testDomainName, getCIDRMap(), &gtm.CreateCIDRMapResponse{
			Resource: getCIDRMap(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetCIDRMap(client.GTM, testDomainName, getCIDRMap(), nil, testutils.FourTimes)

		mockGetDatacenter(client.GTM, testDomainName, datacenterID5400, &dc, nil, testutils.Once)

		mockGetDomainStatus(client.GTM, testDomainName, testutils.Twice)

		mockUpdateCIDRMap(client.GTM, testDomainName, &gtm.UpdateCIDRMapResponse{Status: getDefaultResponseStatus()}, nil)

		mockGetCIDRMap(client.GTM, testDomainName, getCIDRMapUpdated(), nil, testutils.ThreeTimes)

		mockDeleteCIDRMap(client.GTM, testDomainName)

		resourceName := "akamai_gtm_cidrmap.tfexample_cidrmap_1"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "tfexample_cidrmap_1"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/update_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "tfexample_cidrmap_1"),
					),
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})

	t.Run("update CIDRMap failed", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockGetCIDRMap(client.GTM, testDomainName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateCIDRMap(client.GTM, testDomainName, getCIDRMap(), &gtm.CreateCIDRMapResponse{
			Resource: getCIDRMap(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetCIDRMap(client.GTM, testDomainName, getCIDRMap(), nil, testutils.FourTimes)

		mockGetDatacenter(client.GTM, testDomainName, datacenterID5400, &dc, nil, testutils.Once)

		mockGetDomainStatus(client.GTM, testDomainName, testutils.Once)

		mockUpdateCIDRMap(client.GTM, testDomainName, nil, &gtm.Error{
			Type:       "internal_error",
			Title:      "Internal Server Error",
			Detail:     "Error updating resource",
			StatusCode: http.StatusInternalServerError,
		})

		mockGetCIDRMap(client.GTM, testDomainName, getCIDRMap(), nil, testutils.Once)

		mockDeleteCIDRMap(client.GTM, testDomainName)

		resourceName := "akamai_gtm_cidrmap.tfexample_cidrmap_1"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "tfexample_cidrmap_1"),
					),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/update_basic.tf"),
					ExpectError: regexp.MustCompile("API error"),
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})

	t.Run("update CIDRMap domain name - delete and create new CIDRMap", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockGetCIDRMap(client.GTM, testDomainName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockGetDatacenter(client.GTM, testDomainName, datacenterID5400, &dc, nil, testutils.Once)

		mockCreateCIDRMap(client.GTM, testDomainName, getCIDRMap(), &gtm.CreateCIDRMapResponse{
			Resource: getCIDRMap(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetCIDRMap(client.GTM, testDomainName, getCIDRMap(), nil, testutils.FourTimes)

		mockDeleteCIDRMap(client.GTM, testDomainName)

		domainName := "gtm_terra_testdomain_updated.akadns.net"

		mockGetCIDRMap(client.GTM, domainName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockGetDatacenter(client.GTM, domainName, datacenterID5400, &dc, nil, testutils.Once)

		mockCreateCIDRMap(client.GTM, domainName, getCIDRMap(), &gtm.CreateCIDRMapResponse{
			Resource: getCIDRMap(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetCIDRMap(client.GTM, domainName, getCIDRMap(), nil, testutils.ThreeTimes)

		mockDeleteCIDRMap(client.GTM, domainName)

		resourceName := "akamai_gtm_cidrmap.tfexample_cidrmap_1"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "tfexample_cidrmap_1"),
						resource.TestCheckResourceAttr(resourceName, "domain", "gtm_terra_testdomain.akadns.net"),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/domain_update/updated_domain_name.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "tfexample_cidrmap_1"),
						resource.TestCheckResourceAttr(resourceName, "domain", "gtm_terra_testdomain_updated.akadns.net"),
					),
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})

	t.Run("create cidrmap, remove outside of terraform, expect non-empty plan", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockGetCIDRMap(client.GTM, testDomainName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateCIDRMap(client.GTM, testDomainName, getCIDRMap(), &gtm.CreateCIDRMapResponse{
			Resource: getCIDRMap(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetCIDRMap(client.GTM, testDomainName, getCIDRMap(), nil, testutils.Twice)

		mockGetDatacenter(client.GTM, testDomainName, datacenterID5400, &dc, nil, testutils.Once)

		// Mock that the CIDRMap was deleted outside terraform
		mockGetCIDRMap(client.GTM, testDomainName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		// For terraform test framework, we need to mock GetCIDRMap as it would actually exist before deletion
		mockGetCIDRMap(client.GTM, testDomainName, getCIDRMap(), nil, testutils.Once)

		mockDeleteCIDRMap(client.GTM, testDomainName)

		resourceName := "akamai_gtm_cidrmap.tfexample_cidrmap_1"

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "tfexample_cidrmap_1"),
					),
				},
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
					ExpectNonEmptyPlan: true,
					PlanOnly:           true,
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})

	t.Run("create cidrmap failed", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockGetCIDRMap(client.GTM, testDomainName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateCIDRMap(client.GTM, testDomainName, getCIDRMap(), nil, &gtm.Error{StatusCode: http.StatusBadRequest})

		mockGetDatacenter(client.GTM, testDomainName, datacenterID5400, &dc, nil, testutils.Once)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
					ExpectError: regexp.MustCompile("cidrMap create error"),
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})

	t.Run("create cidrmap failed - cidrmap already exists", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockGetCIDRMap(client.GTM, testDomainName, getCIDRMap(), nil, testutils.Once)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
					ExpectError: regexp.MustCompile("cidrMap already exists error"),
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})

	t.Run("create cidrmap denied", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		mockGetCIDRMap(client.GTM, testDomainName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateCIDRMap(client.GTM, testDomainName, getCIDRMap(), &gtm.CreateCIDRMapResponse{
			Resource: getCIDRMap(),
			Status:   getDeniedResponseStatus(),
		}, nil)

		mockGetDatacenter(client.GTM, testDomainName, datacenterID5400, &dc, nil, testutils.Once)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
			Steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
					ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
				},
			},
		})

		client.GTM.AssertExpectations(t)
	})
}

func TestGTMCIDRMapOrder(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"reordered blocks - no diff": {
			pathForUpdate: "testdata/TestResGtmCidrmap/order/blocks/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered assignments - no diff": {
			pathForUpdate: "testdata/TestResGtmCidrmap/order/assignments/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered assignments and blocks - no diff": {
			pathForUpdate: "testdata/TestResGtmCidrmap/order/reorder_assignments_and_blocks.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"change to `name` attribute with different order of assignments and blocks - diff only for `name`": {
			pathForUpdate: "testdata/TestResGtmCidrmap/order/update_name.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"change to `domain` attribute with different order of assignments and blocks - diff only for `domain`": {
			pathForUpdate: "testdata/TestResGtmCidrmap/order/update_domain.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"change to `wait_on_complete` attribute with different order of assignments and blocks - diff only for `wait_on_complete`": {
			pathForUpdate: "testdata/TestResGtmCidrmap/order/update_wait_on_complete.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"reordered and updated blocks - diff only for updated block": {
			pathForUpdate: "testdata/TestResGtmCidrmap/order/blocks/reorder_and_update.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"reordered assignments and updated block - messy diff": {
			pathForUpdate: "testdata/TestResGtmCidrmap/order/assignments/reorder_and_update_block.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"reordered assignments and updated nickname - messy diff": {
			pathForUpdate: "testdata/TestResGtmCidrmap/order/assignments/reorder_and_update_nickname.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := getCIDRMapMocks()
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
				IsUnitTest:               true,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/order/create.tf"),
					},
					{
						Config:             testutils.LoadFixtureString(t, test.pathForUpdate),
						PlanOnly:           test.planOnly,
						ExpectNonEmptyPlan: test.nonEmptyPlan,
					},
				},
			})
			client.GTM.AssertExpectations(t)
		})
	}
}

func TestResGTMCIDRMapImport(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		domainName  string
		mapName     string
		init        func(*gtm.Mock)
		expectError *regexp.Regexp
		stateCheck  resource.ImportStateCheckFunc
	}{
		"happy path - import": {
			domainName: testDomainName,
			mapName:    testCIDRMapName,
			init: func(m *gtm.Mock) {
				// Read
				mockGetCIDRMap(m, testDomainName, getImportedCIDRMap(), nil, testutils.Twice)
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
			mapName:     testCIDRMapName,
			expectError: regexp.MustCompile(`Error: invalid resource ID: :tfexample_cidrmap_1`),
		},
		"expect error - no map name, invalid import ID": {
			domainName:  testDomainName,
			mapName:     "",
			expectError: regexp.MustCompile(`Error: invalid resource ID: gtm_terra_testdomain.akadns.net:`),
		},
		"expect error - read": {
			domainName: testDomainName,
			mapName:    testCIDRMapName,
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetCIDRMap(m, testDomainName, nil, fmt.Errorf("get failed"), testutils.Once)
			},
			expectError: regexp.MustCompile(`get failed`),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()
			if tc.init != nil {
				tc.init(client.GTM)
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, newSubproviderWithConfig(testSubproviderConfig())),
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
			client.GTM.AssertExpectations(t)
		})
	}
}

// getCIDRMapMocks mocks creation and deletion of a resource
func getCIDRMapMocks() *edgegrid.TestClient {
	client := edgegrid.NewTestClient()

	mockGetCIDRMap(client.GTM, testDomainName, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

	mockCreateCIDRMap(client.GTM, testDomainName, getCIDRMapForOrder(), &gtm.CreateCIDRMapResponse{
		Resource: getCIDRMapForOrderResp(),
		Status:   getDefaultResponseStatus(),
	}, nil)

	mockGetCIDRMap(client.GTM, testDomainName, getCIDRMapForOrderResp(), nil, testutils.FourTimes)

	mockGetDatacenter(client.GTM, testDomainName, datacenterID5400, getTestDatacenterResp(), nil, testutils.Once)

	mockGetDomainStatus(client.GTM, testDomainName, testutils.Twice)

	mockDeleteCIDRMap(client.GTM, testDomainName)

	return client
}

func mockGetCIDRMap(m *gtm.Mock, domainName string, cidrMap *gtm.CIDRMap, err error, times int) *mock.Call {
	var resp *gtm.GetCIDRMapResponse
	if cidrMap != nil {
		r := gtm.GetCIDRMapResponse(*cidrMap)
		resp = &r
	}
	return m.On("GetCIDRMap", testutils.MockContext, gtm.GetCIDRMapRequest{
		MapName:    testCIDRMapName,
		DomainName: domainName,
	}).Return(resp, err).Times(times)
}

func mockUpdateCIDRMap(client *gtm.Mock, domainName string, resp *gtm.UpdateCIDRMapResponse, err error) *mock.Call {
	return client.On("UpdateCIDRMap",
		testutils.MockContext,
		gtm.UpdateCIDRMapRequest{
			CIDR:       getCIDRMapUpdated(),
			DomainName: domainName,
		},
	).Return(resp, err).Once()
}

func mockCreateCIDRMap(client *gtm.Mock, domainName string, cidrMap *gtm.CIDRMap, resp *gtm.CreateCIDRMapResponse, err error) *mock.Call {
	return client.On("CreateCIDRMap",
		testutils.MockContext,
		gtm.CreateCIDRMapRequest{
			CIDR:       cidrMap,
			DomainName: domainName,
		},
	).Return(resp, err).Once()
}

func mockDeleteCIDRMap(client *gtm.Mock, domainName string) *mock.Call {
	return client.On("DeleteCIDRMap",
		testutils.MockContext,
		gtm.DeleteCIDRMapRequest{
			MapName:    testCIDRMapName,
			DomainName: domainName,
		},
	).Return(&gtm.DeleteCIDRMapResponse{Status: getDefaultResponseStatus()}, nil).Once()
}

func getImportedCIDRMap() *gtm.CIDRMap {
	return &gtm.CIDRMap{
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: datacenterID5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.CIDRAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					Nickname:     "tfexample_dc_1",
					DatacenterID: datacenterID3131,
				},
				Blocks: []string{"1.2.3.9/24"},
			},
		},
		Name: testCIDRMapName,
	}
}

func getCIDRMap() *gtm.CIDRMap {
	return &gtm.CIDRMap{
		Name: testCIDRMapName,
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: datacenterID5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.CIDRAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3131,
					Nickname:     "tfexample_dc_1",
				},
				Blocks: []string{"1.2.3.9/24"},
			},
		},
	}
}

func getCIDRMapUpdated() *gtm.CIDRMap {
	return &gtm.CIDRMap{
		Name: testCIDRMapName,
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: datacenterID5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.CIDRAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3132,
					Nickname:     "tfexample_dc_2",
				},
				Blocks: []string{"1.2.3.9/16"},
			},
		},
	}
}

func getCIDRMapForOrder() *gtm.CIDRMap {
	return &gtm.CIDRMap{
		Name: testCIDRMapName,
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: datacenterID5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.CIDRAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3131,
					Nickname:     "tfexample_dc_1",
				},
				Blocks: []string{"1.2.3.4/24", "1.2.3.5/24"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3132,
					Nickname:     "tfexample_dc_2",
				},
				Blocks: []string{"1.2.3.7/24", "1.2.3.6/24", "1.2.3.8/24"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3133,
					Nickname:     "tfexample_dc_3",
				},
				Blocks: []string{"1.2.3.9/24", "1.2.3.10/24"},
			},
		},
	}
}

func getCIDRMapForOrderResp() *gtm.CIDRMap {
	return &gtm.CIDRMap{
		Name: testCIDRMapName,
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: datacenterID5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.CIDRAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3131,
					Nickname:     "tfexample_dc_1",
				},
				Blocks: []string{"1.2.3.4/24", "1.2.3.5/24"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3132,
					Nickname:     "tfexample_dc_2",
				},
				Blocks: []string{"1.2.3.6/24", "1.2.3.7/24", "1.2.3.8/24"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3133,
					Nickname:     "tfexample_dc_3",
				},
				Blocks: []string{"1.2.3.9/24", "1.2.3.10/24"},
			},
		},
	}
}
