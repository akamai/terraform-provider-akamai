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

func TestResGtmGeomap(t *testing.T) {
	dc := gtm.Datacenter{
		DatacenterId: geo.DefaultDatacenter.DatacenterId,
		Nickname:     geo.DefaultDatacenter.Nickname,
	}

	t.Run("create geomap", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		})

		resp := gtm.GeoMapResponse{}
		resp.Resource = &geo
		resp.Status = &pendingResponseStatus
		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			mock.AnythingOfType("string"),
		).Return(&resp, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.GeoMap), nil}
		})

		client.On("NewGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&geo, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(&dc, nil)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		client.On("UpdateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.GeoMap), nil}
		})

		client.On("DeleteGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			mock.AnythingOfType("string"),
		).Return(&completeResponseStatus, nil)

		dataSourceName := "akamai_gtm_geomap.tfexample_geomap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmGeomap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_geomap_1"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmGeomap/update_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_geomap_1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create geomap failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			gtmTestDomain,
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("NewGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&geo, nil)

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
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmGeomap/create_basic.tf"),
						ExpectError: regexp.MustCompile("geoMap Create failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create geomap denied", func(t *testing.T) {
		client := &gtm.Mock{}

		dr := gtm.GeoMapResponse{}
		dr.Resource = &geo
		dr.Status = &deniedResponseStatus
		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("*gtm.GeoMap"),
			gtmTestDomain,
		).Return(&dr, nil)

		client.On("NewGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("string"),
		).Return(&geo, nil)

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
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmGeomap/create_basic.tf"),
						ExpectError: regexp.MustCompile("Request could not be completed. Invalid credentials."),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})
}

func TestGTMGeoMapOrder(t *testing.T) {
	tests := map[string]struct {
		client        *gtm.Mock
		pathForCreate string
		pathForUpdate string
		nonEmptyPlan  bool
		planOnly      bool
	}{
		"reorder countries - no diff": {
			client:        getGeoMapMocks(),
			pathForCreate: "testdata/TestResGtmGeomap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmGeomap/order/countries/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"assignments different order - no diff": {
			client:        getGeoMapMocks(),
			pathForCreate: "testdata/TestResGtmGeomap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmGeomap/order/assignments/reorder.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"assignments and countries different order - no diff": {
			client:        getGeoMapMocks(),
			pathForCreate: "testdata/TestResGtmGeomap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmGeomap/order/reorder_assignments_and_countries.tf",
			nonEmptyPlan:  false,
			planOnly:      true,
		},
		"assignments and countries different order with updated `name` - diff only for `name`": {
			client:        getGeoMapMocks(),
			pathForCreate: "testdata/TestResGtmGeomap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmGeomap/order/update_name.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"assignments and countries different order with updated `domain` - diff only for `domain`": {
			client:        getGeoMapMocks(),
			pathForCreate: "testdata/TestResGtmGeomap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmGeomap/order/update_domain.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"assignments and countries different order with updated `wait_on_complete` - diff only for `wait_on_complete`": {
			client:        getGeoMapMocks(),
			pathForCreate: "testdata/TestResGtmGeomap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmGeomap/order/update_wait_on_complete.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reordered assignments and updated countries - messy diff": {
			client:        getGeoMapMocks(),
			pathForCreate: "testdata/TestResGtmGeomap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmGeomap/order/assignments/reorder_and_update_countries.tf",
			nonEmptyPlan:  true, // change to false to see diff
			planOnly:      true,
		},
		"reordered assignments and updated nickname - messy diff": {
			client:        getGeoMapMocks(),
			pathForCreate: "testdata/TestResGtmGeomap/order/create.tf",
			pathForUpdate: "testdata/TestResGtmGeomap/order/assignments/reorder_and_update_nickname.tf",
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

// getGeoMapMocks mock creation and deletion calls for gtm_geomap resource
func getGeoMapMocks() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetGeoMap := client.On("GetGeoMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(nil, &gtm.Error{
		StatusCode: http.StatusNotFound,
	})

	resp := gtm.GeoMapResponse{}
	resp.Resource = &geoDiffOrder
	resp.Status = &pendingResponseStatus
	client.On("CreateGeoMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.GeoMap"),
		mock.AnythingOfType("string"),
	).Return(&resp, nil).Run(func(args mock.Arguments) {
		mockGetGeoMap.ReturnArguments = mock.Arguments{args.Get(1).(*gtm.GeoMap), nil}
	})

	client.On("NewGeoMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
	).Return(&geo, nil)

	client.On("GetDatacenter",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),
	).Return(&dc, nil)

	client.On("GetDomainStatus",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	client.On("DeleteGeoMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("*gtm.GeoMap"),
		mock.AnythingOfType("string"),
	).Return(&completeResponseStatus, nil)

	return client
}

var (
	// geoDiffOrder is gtm.GeoMap structure used in testing of the assignments order
	geoDiffOrder = gtm.GeoMap{
		Name: "tfexample_geomap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterId: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []*gtm.GeoAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterId: 3131,
					Nickname:     "tfexample_dc_1",
				},
				Countries: []string{"GB", "PL", "US", "FR"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterId: 3132,
					Nickname:     "tfexample_dc_2",
				},
				Countries: []string{"GB", "AU"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterId: 3133,
					Nickname:     "tfexample_dc_3",
				},
				Countries: []string{"GB", "BG", "CN", "MC", "TR"},
			},
		},
	}

	geo = gtm.GeoMap{
		Name: "tfexample_geomap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterId: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []*gtm.GeoAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterId: 3131,
					Nickname:     "tfexample_dc_1",
				},
				Countries: []string{"GB"},
			},
		},
	}
)
