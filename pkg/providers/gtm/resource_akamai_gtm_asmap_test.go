package gtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResGTMASMap(t *testing.T) {
	t.Run("create asmap", func(t *testing.T) {
		client := &gtm.Mock{}

		asmap, dc := getASMapTestData()

		getCall := client.On("GetASMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.ASMapResponse{}
		resp.Resource = &asmap
		resp.Status = &pendingResponseStatus

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("CreateASMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.ASMap"),
			mock.AnythingOfType("string"),
		).Return(&gtm.ASMapResponse{
			Resource: &asmap,
			Status:   &gtm.ResponseStatus{},
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{resp.Resource, nil}
		})

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("UpdateASMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.ASMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("DeleteASMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.ASMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_asmap.tfexample_as_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_as_1"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/update_basic.tf"),
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

		_, dc := getASMapTestData()

		client.On("CreateASMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.ASMap"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("asMap Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create asmap denied", func(t *testing.T) {
		client := &gtm.Mock{}

		asmap, dc := getASMapTestData()

		dr := gtm.ASMapResponse{}
		dr.Resource = &asmap
		dr.Status = &deniedResponseStatus
		client.On("CreateASMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.ASMap"),
			gtmTestDomain,
		).Return(&dr, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

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

	t.Run("import asmap", func(t *testing.T) {
		client := &gtm.Mock{}

		asmap, dc := getASMapTestData()

		client.On("GetDatacenter",
			mock.Anything,
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("CreateASMap",
			mock.Anything,
			mock.AnythingOfType("*gtm.ASMap"),
			mock.AnythingOfType("string"),
		).Return(&gtm.ASMapResponse{
			Resource: &asmap,
			Status:   &gtm.ResponseStatus{},
		}, nil)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("GetASMap",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&asmap, nil)

		client.On("DeleteASMap",
			mock.Anything,
			mock.AnythingOfType("*gtm.ASMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_asmap.tfexample_as_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/import_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_as_1"),
						),
					},
					{
						Config:            testutils.LoadFixtureString(t, "testdata/TestResGtmAsmap/create_basic.tf"),
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

func TestGTMASMapOrder(t *testing.T) {
	tests := map[string]struct {
		client        *gtm.Mock
		pathForCreate string
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"reorder as_numbers - no diff": {
			client:        getASMapMocks(),
			pathForCreate: "testdata/TestResGtmAsmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmAsmap/order/as_numbers/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reorder assignments - no diff": {
			client:        getASMapMocks(),
			pathForCreate: "testdata/TestResGtmAsmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmAsmap/order/assignments/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reorder assignments and as_numbers - no diff": {
			client:        getASMapMocks(),
			pathForCreate: "testdata/TestResGtmAsmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmAsmap/order/reorder_assignments_as_numbers.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"change name attribute - diff only for name": {
			client:        getASMapMocks(),
			pathForCreate: "testdata/TestResGtmAsmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmAsmap/order/update_name.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"change wait_on_complete attribute - diff only for wait_on_complete": {
			client:        getASMapMocks(),
			pathForCreate: "testdata/TestResGtmAsmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmAsmap/order/update_wait_on_complete.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"change domain attribute - diff only for domain": {
			client:        getASMapMocks(),
			pathForCreate: "testdata/TestResGtmAsmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmAsmap/order/update_domain.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reorder assignments and change in as_numbers - messy diff": {
			client:        getASMapMocks(),
			pathForCreate: "testdata/TestResGtmAsmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmAsmap/order/assignments/reorder_and_update_as_numbers.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reorder and update as_numbers - diff only for update": {
			client:        getASMapMocks(),
			pathForCreate: "testdata/TestResGtmAsmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmAsmap/order/as_numbers/reorder_and_update.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reorder and update nickname - messy diff": {
			client:        getASMapMocks(),
			pathForCreate: "testdata/TestResGtmAsmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmAsmap/order/assignments/reorder_and_update_nickname.tf",
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

// getASMapMocks mocks creation and deletion of a resource
func getASMapMocks() *gtm.Mock {
	asmap, dc := getASMapTestData()

	client := &gtm.Mock{}

	mockGetAsMap := client.On("GetASMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(nil, &gtm.Error{StatusCode: http.StatusNotFound})

	resp := gtm.ASMapResponse{}
	resp.Resource = &asMapDiffOrder
	resp.Status = &pendingResponseStatus

	client.On("GetDatacenter",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),
	).Return(&dc, nil)

	client.On("CreateASMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.ASMap"),
		mock.AnythingOfType("string"),
	).Return(&gtm.ASMapResponse{
		Resource: &asmap,
		Status:   &gtm.ResponseStatus{},
	}, nil).Run(func(args mock.Arguments) {
		mockGetAsMap.ReturnArguments = mock.Arguments{resp.Resource, nil}
	})

	client.On("GetDomainStatus",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	client.On("DeleteASMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.ASMap"),
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	return client
}

var (
	// asMapDiffOrder represents AsMap structure with values used in tests of the order of assignments and as_numbers
	asMapDiffOrder = gtm.ASMap{
		Name: "tfexample_as_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []*gtm.ASAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3131,
					Nickname:     "tfexample_dc_1",
				},
				ASNumbers: []int64{12222, 16702, 17334},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3132,
					Nickname:     "tfexample_dc_2",
				},
				ASNumbers: []int64{12229, 16703, 17335},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3133,
					Nickname:     "tfexample_dc_3",
				},
				ASNumbers: []int64{1111, 2222, 3333, 4444, 5555},
			},
		},
	}
)

func getASMapTestData() (gtm.ASMap, gtm.Datacenter) {
	asmap := gtm.ASMap{
		Name: "tfexample_as_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []*gtm.ASAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3131,
					Nickname:     "tfexample_dc_1",
				},
				ASNumbers: []int64{12222, 16702, 17334},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3132,
					Nickname:     "tfexample_dc_2",
				},
				ASNumbers: []int64{12229, 16703, 17335},
			},
		},
	}
	dc := gtm.Datacenter{
		DatacenterID: asmap.DefaultDatacenter.DatacenterID,
		Nickname:     asmap.DefaultDatacenter.Nickname,
	}
	return asmap, dc
}
