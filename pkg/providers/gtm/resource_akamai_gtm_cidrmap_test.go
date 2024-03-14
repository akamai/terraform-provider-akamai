package gtm

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/gtm"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResGTMCIDRMap(t *testing.T) {
	dc := gtm.Datacenter{}

	t.Run("create CIDRMap", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetCIDRMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.CIDRMapResponse{}
		resp.Resource = &cidr
		resp.Status = &pendingResponseStatus
		client.On("CreateCIDRMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CIDRMap"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{resp.Resource, nil}
		})

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("UpdateCIDRMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CIDRMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("DeleteCIDRMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CIDRMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

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

	t.Run("create cidrmap failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("CreateCIDRMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CIDRMap"),
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
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmCidrmap/create_basic.tf"),
						ExpectError: regexp.MustCompile("cidrMap Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create cidrmap denied", func(t *testing.T) {
		client := &gtm.Mock{}

		dr := gtm.CIDRMapResponse{}
		dr.Resource = &cidr
		dr.Status = &deniedResponseStatus
		client.On("CreateCIDRMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.CIDRMap"),
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

// getCIDRMapMocks mocks creation and deletion of a resource
func getCIDRMapMocks() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetCIDRMap := client.On("GetCIDRMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(nil, &gtm.Error{
		StatusCode: http.StatusNotFound,
	})

	resp := gtm.CIDRMapResponse{}
	resp.Resource = &cidrMapDiffOrder
	resp.Status = &pendingResponseStatus
	client.On("CreateCIDRMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.CIDRMap"),
		mock.AnythingOfType("string"),
	).Return(&resp, nil).Run(func(args mock.Arguments) {
		mockGetCIDRMap.ReturnArguments = mock.Arguments{resp.Resource, nil}
	})

	client.On("GetDatacenter",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),
	).Return(&dc, nil)

	client.On("GetDomainStatus",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	client.On("DeleteCIDRMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.CIDRMap"),
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	return client
}

var (
	// cidrMapDiffOrder is a gtm.CidrMap structure used in tests of order of assignments and block in gtm_cidrmap resource
	cidrMapDiffOrder = gtm.CIDRMap{
		Name: "tfexample_cidrmap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []*gtm.CIDRAssignment{
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

	cidr = gtm.CIDRMap{
		Name: "tfexample_cidrmap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []*gtm.CIDRAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3131,
					Nickname:     "tfexample_dc_1",
				},
				Blocks: []string{"1.2.3.9/24"},
			},
		},
	}
)
