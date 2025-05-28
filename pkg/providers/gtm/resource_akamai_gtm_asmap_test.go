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

const testASMapName = "tfexample_as_1"

func TestResGTMASMap(t *testing.T) {
	t.Run("create asmap", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetASMap(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockGetDatacenter(client, datacenterID5400, getDefaultDatacenter(), nil, testutils.Once)

		mockCreateASMap(client, getASMapForTestsForCreate(), &gtm.CreateASMapResponse{
			Resource: getASMapForTestsForCreateResponse(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetASMap(client, getASMapForTestsForCreateResponse(), nil, testutils.FourTimes)

		mockGetDomainStatus(client, testutils.Twice)

		mockUpdateASMap(client, &gtm.UpdateASMapResponse{Status: getDefaultResponseStatus()}, nil)

		mockGetASMap(client, getASMapUpdateResponse(), nil, testutils.ThreeTimes)

		mockDeleteASMap(client)

		resourceName := "akamai_gtm_asmap.tfexample_as_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "tfexample_as_1"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "tfexample_as_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("update asmap failed", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetASMap(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockGetDatacenter(client, datacenterID5400, getDefaultDatacenter(), nil, testutils.Once)

		mockCreateASMap(client, getASMapForTestsForCreate(), &gtm.CreateASMapResponse{
			Resource: getASMapForTestsForCreateResponse(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetASMap(client, getASMapForTestsForCreateResponse(), nil, testutils.FourTimes)

		mockGetDomainStatus(client, testutils.Once)

		mockUpdateASMap(client, nil, &gtm.Error{
			Type:       "internal_error",
			Title:      "Internal Server Error",
			Detail:     "Error updating asmap",
			StatusCode: http.StatusInternalServerError,
		})

		mockGetASMap(client, getASMapForTestsForCreateResponse(), nil, testutils.Once)

		mockDeleteASMap(client)

		resourceName := "akamai_gtm_asmap.tfexample_as_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "tfexample_as_1"),
						),
					},
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/update_basic.tf"),
						ExpectError: regexp.MustCompile("API error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create asmap, remove outside of terraform, expect non-empty plan", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetASMap(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockGetDatacenter(client, datacenterID5400, getDefaultDatacenter(), nil, testutils.Once)

		mockCreateASMap(client, getASMapForTestsForCreate(), &gtm.CreateASMapResponse{
			Resource: getASMapForTestsForCreateResponse(),
			Status:   getDefaultResponseStatus(),
		}, nil)

		mockGetASMap(client, getASMapForTestsForCreateResponse(), nil, testutils.Twice)

		// Mock that the ASMap was deleted outside terraform
		mockGetASMap(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		// For terraform test framework, we need to mock GetASMap as it would actually exist before deletion
		mockGetASMap(client, getASMapForTestsForCreateResponse(), nil, testutils.Once)

		mockDeleteASMap(client)

		resourceName := "akamai_gtm_asmap.tfexample_as_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceName, "name", "tfexample_as_1"),
						),
					},
					{
						Config:             testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/create_basic.tf"),
						ExpectNonEmptyPlan: true,
						PlanOnly:           true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create asmap failed", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetASMap(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateASMap(client, getASMapForTestsForCreate(), nil, &gtm.Error{StatusCode: http.StatusBadRequest})

		mockGetDatacenter(client, datacenterID5400, getDefaultDatacenter(), nil, testutils.Once)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("asMap create error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create asmap failed - asmap already exists", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetASMap(client, getASMapForTestsForCreateResponse(), nil, testutils.Once)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("asMap already exists error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create asmap denied", func(t *testing.T) {
		client := &gtm.Mock{}

		mockGetASMap(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

		mockCreateASMap(client, getASMapForTestsForCreate(), &gtm.CreateASMapResponse{
			Resource: getASMapForTestsForCreateResponse(),
			Status:   getDeniedResponseStatus(),
		}, nil)

		mockGetDatacenter(client, datacenterID5400, getTestDatacenterResp(), nil, testutils.Once)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestGTMASMapOrder(t *testing.T) {
	tests := map[string]struct {
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"reorder as_numbers - no diff": {
			pathForUpdate: "testdata/TestResGtmAsmap/order/as_numbers/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reorder assignments - no diff": {
			pathForUpdate: "testdata/TestResGtmAsmap/order/assignments/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reorder assignments and as_numbers - no diff": {
			pathForUpdate: "testdata/TestResGtmAsmap/order/reorder_assignments_as_numbers.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"change name attribute - diff only for name": {
			pathForUpdate: "testdata/TestResGtmAsmap/order/update_name.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"change wait_on_complete attribute - diff only for wait_on_complete": {
			pathForUpdate: "testdata/TestResGtmAsmap/order/update_wait_on_complete.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"change domain attribute - diff only for domain": {
			pathForUpdate: "testdata/TestResGtmAsmap/order/update_domain.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"reorder assignments and change in as_numbers - messy diff": {
			pathForUpdate: "testdata/TestResGtmAsmap/order/assignments/reorder_and_update_as_numbers.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"reorder and update as_numbers - diff only for update": {
			pathForUpdate: "testdata/TestResGtmAsmap/order/as_numbers/reorder_and_update.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
		"reorder and update nickname - messy diff": {
			pathForUpdate: "testdata/TestResGtmAsmap/order/assignments/reorder_and_update_nickname.tf",
			nonEmptyPlan:  true,
			planOnly:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := getASMapMocks()
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps: []resource.TestStep{
						{
							Config: testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/order/create.tf"),
						},
						{
							Config:             testutils.LoadFixtureString(t, test.pathForUpdate),
							PlanOnly:           test.planOnly,
							ExpectNonEmptyPlan: test.nonEmptyPlan,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

func TestResGTMASMapImport(t *testing.T) {
	tests := map[string]struct {
		domainName  string
		mapName     string
		init        func(*gtm.Mock)
		expectError *regexp.Regexp
		stateCheck  resource.ImportStateCheckFunc
	}{
		"happy path - import": {
			domainName: testDomainName,
			mapName:    testASMapName,
			init: func(m *gtm.Mock) {
				// Read
				mockGetASMap(m, getImportedASMap(), nil, testutils.Twice)
			},
			stateCheck: test.NewImportChecker().
				CheckEqual("domain", "gtm_terra_testdomain.akadns.net").
				CheckEqual("name", "tfexample_as_1").
				CheckEqual("default_datacenter.0.datacenter_id", "5400").
				CheckEqual("default_datacenter.0.nickname", "default datacenter").
				CheckEqual("assignment.0.datacenter_id", "3131").
				CheckEqual("assignment.0.nickname", "tfexample_dc_1").
				CheckEqual("assignment.0.as_numbers.0", "12222").
				CheckEqual("assignment.0.as_numbers.1", "16702").
				CheckEqual("assignment.0.as_numbers.2", "17334").
				CheckEqual("wait_on_complete", "true").Build(),
		},
		"expect error - no domain name, invalid import ID": {
			domainName:  "",
			mapName:     testASMapName,
			expectError: regexp.MustCompile(`Error: invalid resource ID: :tfexample_as_1`),
		},
		"expect error - no map name, invalid import ID": {
			domainName:  testDomainName,
			mapName:     "",
			expectError: regexp.MustCompile(`Error: invalid resource ID: gtm_terra_testdomain.akadns.net:`),
		},
		"expect error - read": {
			domainName: testDomainName,
			mapName:    testASMapName,
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetASMap(m, nil, fmt.Errorf("get failed"), testutils.Once)
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
							ResourceName:     "akamai_gtm_asmap.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/import_basic.tf"),
							ExpectError:      tc.expectError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

// getASMapMocks mocks creation and deletion of a resource
func getASMapMocks() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetASMap(client, nil, &gtm.Error{StatusCode: http.StatusNotFound}, testutils.Once)

	mockGetDatacenter(client, datacenterID5400, getTestDatacenterResp(), nil, testutils.Once)

	mockCreateASMap(client, getASMapForOrder(), &gtm.CreateASMapResponse{
		Resource: getASMapForOrderResp(),
		Status:   getDefaultResponseStatus(),
	}, nil)

	mockGetASMap(client, getASMapForOrderResp(), nil, testutils.FourTimes)

	mockGetDomainStatus(client, testutils.Twice)

	mockDeleteASMap(client)

	return client
}

func mockGetASMap(m *gtm.Mock, asMap *gtm.ASMap, err error, times int) *mock.Call {
	var resp *gtm.GetASMapResponse
	if asMap != nil {
		r := gtm.GetASMapResponse(*asMap)
		resp = &r
	}
	return m.On("GetASMap", testutils.MockContext, gtm.GetASMapRequest{
		ASMapName:  testASMapName,
		DomainName: testDomainName,
	}).Return(resp, err).Times(times)
}

func mockUpdateASMap(client *gtm.Mock, resp *gtm.UpdateASMapResponse, err error) *mock.Call {
	return client.On("UpdateASMap",
		testutils.MockContext,
		gtm.UpdateASMapRequest{
			ASMap:      getASMapUpdate(),
			DomainName: testDomainName,
		},
	).Return(resp, err).Once()
}

func mockCreateASMap(client *gtm.Mock, asMap *gtm.ASMap, resp *gtm.CreateASMapResponse, err error) *mock.Call {
	return client.On("CreateASMap",
		testutils.MockContext,
		gtm.CreateASMapRequest{
			ASMap:      asMap,
			DomainName: testDomainName,
		},
	).Return(resp, err).Once()
}

func mockDeleteASMap(client *gtm.Mock) *mock.Call {
	return client.On("DeleteASMap",
		testutils.MockContext,
		gtm.DeleteASMapRequest{
			ASMapName:  testASMapName,
			DomainName: testDomainName,
		},
	).Return(&gtm.DeleteASMapResponse{Status: getDefaultResponseStatus()}, nil).Once()
}

func getImportedASMap() *gtm.ASMap {
	return &gtm.ASMap{
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: datacenterID5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.ASAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3131,
					Nickname:     "tfexample_dc_1",
				},
				ASNumbers: []int64{12222, 16702, 17334},
			},
		},
		Name: testASMapName,
	}
}

func getASMapForTestsForCreate() *gtm.ASMap {
	return &gtm.ASMap{
		Name: testASMapName,
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: datacenterID5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.ASAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3131,
					Nickname:     "tfexample_dc_1",
				},
				ASNumbers: []int64{16702, 12222, 17334},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3132,
					Nickname:     "tfexample_dc_2",
				},
				ASNumbers: []int64{16703, 12229, 17335},
			},
		},
	}
}

func getASMapUpdate() *gtm.ASMap {
	asMap := getASMapForTestsForCreate()
	asMap.Assignments = []gtm.ASAssignment{
		{
			DatacenterBase: gtm.DatacenterBase{
				DatacenterID: datacenterID3132,
				Nickname:     "tfexample_dc_2",
			},
			ASNumbers: []int64{16701, 12223, 17333},
		},
		{
			DatacenterBase: gtm.DatacenterBase{
				DatacenterID: datacenterID3133,
				Nickname:     "tfexample_dc_3",
			},
			ASNumbers: []int64{12228, 17336, 16704},
		},
	}
	return asMap
}

func getASMapUpdateResponse() *gtm.ASMap {
	asMap := getASMapForTestsForCreateResponse()
	asMap.Assignments = []gtm.ASAssignment{
		{
			DatacenterBase: gtm.DatacenterBase{
				DatacenterID: datacenterID3132,
				Nickname:     "tfexample_dc_2",
			},
			ASNumbers: []int64{12223, 16701, 17333},
		},
		{
			DatacenterBase: gtm.DatacenterBase{
				DatacenterID: datacenterID3133,
				Nickname:     "tfexample_dc_3",
			},
			ASNumbers: []int64{12228, 16704, 17336},
		},
	}
	return asMap
}

func getASMapForTestsForCreateResponse() *gtm.ASMap {
	return &gtm.ASMap{
		Name: testASMapName,
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: datacenterID5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.ASAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3131,
					Nickname:     "tfexample_dc_1",
				},
				ASNumbers: []int64{12222, 16702, 17334},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3132,
					Nickname:     "tfexample_dc_2",
				},
				ASNumbers: []int64{12229, 16703, 17335},
			},
		},
	}
}

func getASMapForOrder() *gtm.ASMap {
	return &gtm.ASMap{
		Name: testASMapName,
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: datacenterID5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.ASAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3131,
					Nickname:     "tfexample_dc_1",
				},
				ASNumbers: []int64{16702, 12222, 17334},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3132,
					Nickname:     "tfexample_dc_2",
				},
				ASNumbers: []int64{16703, 12229, 17335},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3133,
					Nickname:     "tfexample_dc_3",
				},
				ASNumbers: []int64{2222, 1111, 5555, 3333, 4444},
			},
		},
	}
}

func getASMapForOrderResp() *gtm.ASMap {
	return &gtm.ASMap{
		Name: testASMapName,
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: datacenterID5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.ASAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3131,
					Nickname:     "tfexample_dc_1",
				},
				ASNumbers: []int64{12222, 16702, 17334},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3132,
					Nickname:     "tfexample_dc_2",
				},
				ASNumbers: []int64{12229, 16703, 17335},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: datacenterID3133,
					Nickname:     "tfexample_dc_3",
				},
				ASNumbers: []int64{1111, 2222, 3333, 4444, 5555},
			},
		},
	}
}

func getDefaultDatacenter() *gtm.Datacenter {
	return &gtm.Datacenter{
		DatacenterID: datacenterID5400,
		Nickname:     "default datacenter",
	}
}
