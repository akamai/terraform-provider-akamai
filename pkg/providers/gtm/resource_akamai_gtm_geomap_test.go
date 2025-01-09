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

func TestResGTMGeoMap(t *testing.T) {
	dc := gtm.Datacenter{
		DatacenterID: geomap.DefaultDatacenter.DatacenterID,
		Nickname:     geomap.DefaultDatacenter.Nickname,
	}

	t.Run("create geomap", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetGeoMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Twice()

		resp := geomap
		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.CreateGeoMapRequest"),
		).Return(&gtm.CreateGeoMapResponse{
			Resource: geoMapCreate.Resource,
			Status:   geoMapCreate.Status,
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		})

		client.On("GetGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetGeoMapRequest"),
		).Return(&resp, nil).Times(3)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&dc, nil)

		client.On("GetDomainStatus",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDomainStatusRequest"),
		).Return(getDomainStatusResponseStatus, nil)

		client.On("UpdateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.UpdateGeoMapRequest"),
		).Return(updateGeoMapResponseStatus, nil)

		client.On("GetGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetGeoMapRequest"),
		).Return(&geomapUpdate, nil).Times(3)

		client.On("DeleteGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.DeleteGeoMapRequest"),
		).Return(deleteGeoMapResponseStatus, nil)

		dataSourceName := "akamai_gtm_geomap.tfexample_geomap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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

	t.Run("create GEO map, remove outside of terraform, expect non-empty plan", func(t *testing.T) {
		client := &gtm.Mock{}

		getCall := client.On("GetGeoMap",
			mock.Anything,
			mock.AnythingOfType("gtm.GetGeoMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		resp := geomap
		client.On("CreateGeoMap",
			mock.Anything,
			mock.AnythingOfType("gtm.CreateGeoMapRequest"),
		).Return(&gtm.CreateGeoMapResponse{
			Resource: geoMapCreate.Resource,
			Status:   geoMapCreate.Status,
		}, nil).Run(func(args mock.Arguments) {
			getCall.ReturnArguments = mock.Arguments{&resp, nil}
		}).Once()

		client.On("GetGeoMap",
			mock.Anything,
			mock.AnythingOfType("gtm.GetGeoMapRequest"),
		).Return(&resp, nil).Twice()

		client.On("GetDatacenter",
			mock.Anything,
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&dc, nil).Once()

		// Mock that the GEOMap was deleted outside terraform
		client.On("GetGeoMap",
			mock.Anything,
			mock.AnythingOfType("gtm.GetGeoMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		// For terraform test framework, we need to mock GetGEOMap as it would actually exist before deletion
		client.On("GetGeoMap",
			mock.Anything,
			mock.AnythingOfType("gtm.GetGeoMapRequest"),
		).Return(&geomapUpdate, nil).Once()

		client.On("DeleteGeoMap",
			mock.Anything,
			mock.AnythingOfType("gtm.DeleteGeoMapRequest"),
		).Return(deleteGeoMapResponseStatus, nil).Once()

		dataSourceName := "akamai_gtm_geomap.tfexample_geomap_1"

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResGtmGeomap/create_basic.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(dataSourceName, "name", "tfexample_geomap_1"),
						),
					},
					{
						Config:             testutils.LoadFixtureString(t, "testdata/TestResGtmGeomap/create_basic.tf"),
						ExpectNonEmptyPlan: true,
						PlanOnly:           true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create geomap failed", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("GetGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetGeoMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.CreateGeoMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusBadRequest,
		})

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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

	t.Run("create geomap failed - geomap already exists", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("GetGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetGeoMapRequest"),
		).Return(&geomap, nil).Once()

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResGtmGeomap/create_basic.tf"),
						ExpectError: regexp.MustCompile("geoMap already exists error"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("create geomap denied", func(t *testing.T) {
		client := &gtm.Mock{}

		client.On("GetGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetGeoMapRequest"),
		).Return(nil, &gtm.Error{
			StatusCode: http.StatusNotFound,
		}).Once()

		dr := gtm.CreateGeoMapResponse{}
		dr.Resource = geoMapCreate.Resource
		dr.Status = &deniedResponseStatus
		client.On("CreateGeoMap",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.CreateGeoMapRequest"),
		).Return(&dr, nil)

		client.On("GetDatacenter",
			mock.Anything, // ctx is irrelevant for this test
			mock.AnythingOfType("gtm.GetDatacenterRequest"),
		).Return(&dc, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
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

func TestResGTMGeoMapImport(t *testing.T) {
	tests := map[string]struct {
		domainName  string
		mapName     string
		init        func(*gtm.Mock)
		expectError *regexp.Regexp
		stateCheck  resource.ImportStateCheckFunc
	}{
		"happy path - import": {
			domainName: "test_domain",
			mapName:    "tfexample_geomap_1",
			init: func(m *gtm.Mock) {
				// Read
				importedGeomap := gtm.GetGeoMapResponse(*getImportedGeoMap())
				mockGetGeoMap(m, &importedGeomap, nil).Times(2)
			},
			stateCheck: test.NewImportChecker().
				CheckEqual("domain", "test_domain").
				CheckEqual("name", "tfexample_geomap_1").
				CheckEqual("default_datacenter.0.datacenter_id", "5400").
				CheckEqual("default_datacenter.0.nickname", "default datacenter").
				CheckEqual("assignment.0.datacenter_id", "3131").
				CheckEqual("assignment.0.nickname", "tfexample_dc_1").
				CheckEqual("assignment.0.countries.0", "GB").
				CheckEqual("wait_on_complete", "true").Build(),
		},
		"expect error - no domain name, invalid import ID": {
			domainName:  "",
			mapName:     "tfexample_geomap_1",
			expectError: regexp.MustCompile(`Error: invalid resource ID: :tfexample_geomap_1`),
		},
		"expect error - no map name, invalid import ID": {
			domainName:  "test_domain",
			mapName:     "",
			expectError: regexp.MustCompile(`Error: invalid resource ID: test_domain:`),
		},
		"expect error - read": {
			domainName: "test_domain",
			mapName:    "tfexample_geomap_1",
			init: func(m *gtm.Mock) {
				// Read - error
				mockGetGeoMap(m, nil, fmt.Errorf("get failed")).Once()
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
							ResourceName:     "akamai_gtm_geomap.test",
							Config:           testutils.LoadFixtureString(t, "testdata/TestResGtmGeomap/import_basic.tf"),
							ExpectError:      tc.expectError,
						},
					},
				})
			})
			client.AssertExpectations(t)
		})
	}
}

// getGeoMapMocks mock creation and deletion calls for gtm_geomap resource
func getGeoMapMocks() *gtm.Mock {
	client := &gtm.Mock{}

	mockGetGeoMap := client.On("GetGeoMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("gtm.GetGeoMapRequest"),
	).Return(nil, &gtm.Error{
		StatusCode: http.StatusNotFound,
	})

	resp := geoDiffOrder
	client.On("CreateGeoMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("gtm.CreateGeoMapRequest"),
	).Return(&gtm.CreateGeoMapResponse{
		Resource: geoMapCreateDiif.Resource,
		Status:   geoMapCreateDiif.Status,
	}, nil).Run(func(args mock.Arguments) {
		mockGetGeoMap.ReturnArguments = mock.Arguments{&resp, nil}
	})

	client.On("GetDatacenter",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("gtm.GetDatacenterRequest"),
	).Return(&dc, nil)

	client.On("GetDomainStatus",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("gtm.GetDomainStatusRequest"),
	).Return(getDomainStatusResponseStatus, nil)

	client.On("DeleteGeoMap",
		mock.Anything, // ctx is irrelevant for this test
		mock.AnythingOfType("gtm.DeleteGeoMapRequest"),
	).Return(deleteGeoMapResponseStatus, nil)

	return client
}

func mockGetGeoMap(m *gtm.Mock, resp *gtm.GetGeoMapResponse, err error) *mock.Call {
	return m.On("GetGeoMap", mock.Anything, gtm.GetGeoMapRequest{
		MapName:    "tfexample_geomap_1",
		DomainName: "test_domain",
	}).Return(resp, err)
}

func getImportedGeoMap() *gtm.GeoMap {
	return &gtm.GeoMap{
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.GeoAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3131,
					Nickname:     "tfexample_dc_1",
				},
				Countries: []string{"GB"},
			},
		},
		Name: "tfexample_geomap_1",
	}
}

var (
	// geoDiffOrder is gtm.GeoMap structure used in testing of the assignments order
	geoDiffOrder = gtm.GetGeoMapResponse{
		Name: "tfexample_geomap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.GeoAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3131,
					Nickname:     "tfexample_dc_1",
				},
				Countries: []string{"GB", "PL", "US", "FR"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3132,
					Nickname:     "tfexample_dc_2",
				},
				Countries: []string{"GB", "AU"},
			},
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3133,
					Nickname:     "tfexample_dc_3",
				},
				Countries: []string{"GB", "BG", "CN", "MC", "TR"},
			},
		},
	}

	geoMapCreateDiif = gtm.CreateGeoMapResponse{
		Resource: &gtm.GeoMap{
			Name: "tfexample_geomap_1",
			DefaultDatacenter: &gtm.DatacenterBase{
				DatacenterID: 5400,
				Nickname:     "default datacenter",
			},
			Assignments: []gtm.GeoAssignment{
				{
					DatacenterBase: gtm.DatacenterBase{
						DatacenterID: 3131,
						Nickname:     "tfexample_dc_1",
					},
					Countries: []string{"GB", "PL", "US", "FR"},
				},
				{
					DatacenterBase: gtm.DatacenterBase{
						DatacenterID: 3132,
						Nickname:     "tfexample_dc_2",
					},
					Countries: []string{"GB", "AU"},
				},
				{
					DatacenterBase: gtm.DatacenterBase{
						DatacenterID: 3133,
						Nickname:     "tfexample_dc_3",
					},
					Countries: []string{"GB", "BG", "CN", "MC", "TR"},
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

	geomap = gtm.GetGeoMapResponse{
		Name: "tfexample_geomap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.GeoAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3131,
					Nickname:     "tfexample_dc_1",
				},
				Countries: []string{"GB"},
			},
		},
	}

	geomapUpdate = gtm.GetGeoMapResponse{
		Name: "tfexample_geomap_1",
		DefaultDatacenter: &gtm.DatacenterBase{
			DatacenterID: 5400,
			Nickname:     "default datacenter",
		},
		Assignments: []gtm.GeoAssignment{
			{
				DatacenterBase: gtm.DatacenterBase{
					DatacenterID: 3132,
					Nickname:     "tfexample_dc_2",
				},
				Countries: []string{"US"},
			},
		},
	}

	geoMapCreate = gtm.CreateGeoMapResponse{
		Resource: &gtm.GeoMap{
			Name: "tfexample_geomap_1",
			DefaultDatacenter: &gtm.DatacenterBase{
				DatacenterID: 5400,
				Nickname:     "default datacenter",
			},
			Assignments: []gtm.GeoAssignment{
				{
					DatacenterBase: gtm.DatacenterBase{
						DatacenterID: 3131,
						Nickname:     "tfexample_dc_1",
					},
					Countries: []string{"GB"},
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

	updateGeoMapResponseStatus = &gtm.UpdateGeoMapResponse{
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
	deleteGeoMapResponseStatus = &gtm.DeleteGeoMapResponse{
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
