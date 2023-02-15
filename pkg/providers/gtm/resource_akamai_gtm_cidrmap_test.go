package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/gtm"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResGtmCidrmap(t *testing.T) {
	dc := gtm.Datacenter{}

	t.Run("create cidrmap", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.CidrMapResponse{}
		resp.Resource = &cidr
		resp.Status = &pendingResponseStatus
		client.On("CreateCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{resp.Resource, nil}
		})

		client.On("NewCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&cidr, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("UpdateCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("DeleteCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_cidrmap.tfexample_cidrmap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResGtmCidrmap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_cidrmap_1"),
						),
					},
					{
						Config: loadFixtureString("testdata/TestResGtmCidrmap/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_cidrmap_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create cidrmap failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("CreateCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("NewCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&cidr, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmCidrmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("cidrMap Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create cidrmap denied", func(t *testing.T) {
		client := &gtm.Mock{}

		dr := gtm.CidrMapResponse{}
		dr.Resource = &cidr
		dr.Status = &deniedResponseStatus
		client.On("CreateCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CidrMap"),
			gtmTestDomain,
		).Return(&dr, nil)

		client.On("NewCidrMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&cidr, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestResGtmCidrmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestGTMCidrMapOrder(t *testing.T) {
	tests := map[string]struct {
		client        *gtm.Mock
		pathForCreate string
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"reordered blocks - no diff": {
			client:        getCidrMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/blocks/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered assignments - no diff": {
			client:        getCidrMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/assignments/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"reordered assignments and blocks - no diff": {
			client:        getCidrMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/reorder_assignments_and_blocks.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"change to `name` attribute with different order of assignments and blocks - diff only for `name`": {
			client:        getCidrMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/update_name.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"change to `domain` attribute with different order of assignments and blocks - diff only for `domain`": {
			client:        getCidrMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/update_domain.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"change to `wait_on_complete` attribute with different order of assignments and blocks - diff only for `wait_on_complete`": {
			client:        getCidrMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/update_wait_on_complete.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reordered and updated blocks - diff only for updated block": {
			client:        getCidrMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/blocks/reorder_and_update.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reordered assignments and updated block - messy diff": {
			client:        getCidrMapMocks(),
			pathForCreate: "testdata/TestResGtmCidrmap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmCidrmap/order/assignments/reorder_and_update_block.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reordered assignments and updated nickname - messy diff": {
			client:        getCidrMapMocks(),
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
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps: []resource.TestStep{
						{
							Config: loadFixtureString(test.pathForCreate),
						},
						{
							Config:             loadFixtureString(test.pathForUpdate),
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

// getCidrMapMocks mocks creation and deletion of a resource
func getCidrMapMocks() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetCIDRMap := client.On("GetCidrMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(nil, &gtm.Error{
		StatusCode: http.StatusNotFound,
	})

	resp := gtm.CidrMapResponse{}
	resp.Resource = &cidrMapDiffOrder
	resp.Status = &pendingResponseStatus
	client.On("CreateCidrMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.CidrMap"),
		mock.AnythingOfType("string"),
	).Return(&resp, nil).Run(func(args mock.Arguments) {
		mockGetCIDRMap.ReturnArguments = mock.Arguments{resp.Resource, nil}
	})

	client.On("NewCidrMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
	).Return(&cidr, nil)

	client.On("GetDatacenter",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),
	).Return(&dc, nil)

	client.On("GetDomainStatus",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	client.On("DeleteCidrMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.CidrMap"),
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	return client
}

var (
	// cidrMapDiffOrder is a gtm.CidrMap structure used in tests of order of assignments and block in gtm_cidrmap resource
	cidrMapDiffOrder = gtm.CidrMap{
		Name: "tfexample_cidrmap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterId: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []*gtm.CidrAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterId: 3131,
					Nickname:     "tfexample_dc_1",
				},
				Blocks: []string{"1.2.3.4/24", "1.2.3.5/24"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterId: 3132,
					Nickname:     "tfexample_dc_2",
				},
				Blocks: []string{"1.2.3.6/24", "1.2.3.7/24", "1.2.3.8/24"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterId: 3133,
					Nickname:     "tfexample_dc_3",
				},
				Blocks: []string{"1.2.3.9/24", "1.2.3.10/24"},
			},
		},
	}

	cidr = gtm.CidrMap{
		Name: "tfexample_cidrmap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterId: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []*gtm.CidrAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterId: 3131,
					Nickname:     "tfexample_dc_1",
				},
				Blocks: []string{"1.2.3.9/24"},
			},
		},
	}
)
