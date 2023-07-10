package gtm

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

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
				ProviderFactories: testAccProviders,
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
				ProviderFactories: testAccProviders,
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
				ProviderFactories: testAccProviders,
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
				ProviderFactories: testAccProviders,
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

func TestGTMAsMapOrder(t *testing.T) {
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

// getASMapMocks mocks creation and deletion of a resource
func getASMapMocks() *gtm.Mock {
	dc := gtm.Datacenter{
		DatacenterId: asmap.DefaultDatacenter.DatacenterId,
		Nickname:     asmap.DefaultDatacenter.Nickname,
	}

	client := &gtm.Mock{}

	mockGetAsMap := client.On("GetAsMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(nil, &gtm.Error{StatusCode: http.StatusNotFound})

	resp := gtm.AsMapResponse{}
	resp.Resource = &asMapDiffOrder
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
		mockGetAsMap.ReturnArguments = mock.Arguments{resp.Resource, nil}
	})

	client.On("GetDomainStatus",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	client.On("DeleteAsMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.AsMap"),
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	return client
}

var (
	// asMapDiffOrder represents AsMap structure with values used in tests of the order of assignments and as_numbers
	asMapDiffOrder = gtm.AsMap{
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
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterId: 3133,
					Nickname:     "tfexample_dc_3",
				},
				AsNumbers: []int64{1111, 2222, 3333, 4444, 5555},
			},
		},
	}

	asmap = gtm.AsMap{
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
)
