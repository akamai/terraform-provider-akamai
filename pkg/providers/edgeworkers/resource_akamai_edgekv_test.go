package edgeworkers

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/edgeworkers"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestResourceEdgeKV(t *testing.T) {
	initWindow = time.Duration(1) * time.Millisecond

	basicData := edgeKVmockData{
		network:   "staging",
		name:      "DevExpTest",
		retention: ptr.To(86401),
		groupID:   ptr.To(1234),
	}

	tests := map[string]struct {
		init  func(*edgeworkers.Mock)
		steps []resource.TestStep
		data  edgeKVmockData
	}{
		"basic": {
			init: func(m *edgeworkers.Mock) {
				// create
				mockEdgeKVCreate(m, basicData)
				// read
				mockEdgeKVRead(m, basicData).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", "DevExpTest:staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", "DevExpTest"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", 1234)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 86401)),
					),
				},
			},
		},
		"basic - retention 0": {
			init: func(m *edgeworkers.Mock) {
				data := basicData
				data.retention = ptr.To(0)
				// create
				mockEdgeKVCreate(m, data)
				// read
				mockEdgeKVRead(m, data).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic_retention_0.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", "DevExpTest:staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", "DevExpTest"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", 1234)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 0)),
					),
				},
			},
		},
		"error in namespace initialization": {
			init: func(m *edgeworkers.Mock) {
				// create
				mockGetEdgeKVInitializationStatus(m, "UNINITIALIZED")
				// expect an error on InitializeEdgeKV
				m.On("InitializeEdgeKV", testutils.MockContext).Return(nil, fmt.Errorf("error on initialization edgeKV")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					ExpectError: regexp.MustCompile("error on initialization edgeKV"),
				},
			},
		},
		"namespace status PENDING - only waiting for INITIALIZED status": {
			init: func(m *edgeworkers.Mock) {
				// create
				// initial check
				mockGetEdgeKVInitializationStatus(m, "PENDING")
				// wait function - 1st call
				mockGetEdgeKVInitializationStatus(m, "PENDING")
				// namespace initialized, exit waiting function
				mockGetEdgeKVInitializationStatus(m, "INITIALIZED")
				// proceed with create flow
				mockCreateEdgeKVNamespace(m, basicData)
				// read
				mockEdgeKVRead(m, basicData).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", "DevExpTest:staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", "DevExpTest"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", 1234)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 86401)),
					),
				},
			},
		},
		"namespace already INITIALIZED": {
			init: func(m *edgeworkers.Mock) {
				// create
				// Namespace status is already initialized, we skip waiting and proceed normally
				mockGetEdgeKVInitializationStatus(m, "INITIALIZED")
				mockCreateEdgeKVNamespace(m, basicData)
				// read
				mockEdgeKVRead(m, basicData).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", "DevExpTest:staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", "DevExpTest"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", 1234)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 86401)),
					),
				},
			},
		},
		"error creating namespace": {
			init: func(m *edgeworkers.Mock) {
				// create
				mockGetEdgeKVInitializationStatus(m, "UNINITIALIZED")
				mockInitializeEdgeKV(m)
				mockGetEdgeKVInitializationStatus(m, "INITIALIZED")
				// expect error on creating EdgeKV namespace
				m.On("CreateEdgeKVNamespace", testutils.MockContext, edgeworkers.CreateEdgeKVNamespaceRequest{
					Network: basicData.network,
					Namespace: edgeworkers.Namespace{
						Name:        basicData.name,
						GeoLocation: basicData.geoLocation,
						Retention:   basicData.retention,
						GroupID:     basicData.groupID,
					},
				}).Return(nil, fmt.Errorf("error creating edgeKV namespace")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					ExpectError: regexp.MustCompile("error creating edgeKV namespace"),
				},
			},
		},
		"basic error read": {
			init: func(m *edgeworkers.Mock) {
				// create
				mockEdgeKVCreate(m, basicData)
				// read - expect error on GetEdgeKVNamespace
				m.On("GetEdgeKVNamespace", testutils.MockContext, edgeworkers.GetEdgeKVNamespaceRequest{
					Network: basicData.network,
					Name:    basicData.name,
				}).Return(nil, fmt.Errorf("error reading edgekv namespace")).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					ExpectError: regexp.MustCompile("error reading edgekv namespace"),
				},
			},
		},
		"basic no diff no update": {
			init: func(m *edgeworkers.Mock) {
				// create
				mockEdgeKVCreate(m, basicData)
				// read
				mockEdgeKVRead(m, basicData).Times(2)
				// read before update
				mockEdgeKVRead(m, basicData).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", "DevExpTest:staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", "DevExpTest"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", 1234)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 86401)),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", "DevExpTest:staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", "DevExpTest"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", 1234)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 86401)),
					),
				},
			},
		},
		"ignore diff on group_id": {
			init: func(m *edgeworkers.Mock) {
				// create
				mockEdgeKVCreate(m, basicData)
				// read
				mockEdgeKVRead(m, basicData).Times(2)
				// read before update
				mockEdgeKVRead(m, basicData).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", "DevExpTest:staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", "DevExpTest"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", 1234)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 86401)),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/update_diff_group_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", "DevExpTest:staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", "DevExpTest"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", 1234)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 86401)),
					),
				},
			},
		},
		"basic diff retention update": {
			init: func(m *edgeworkers.Mock) {
				data := basicData
				// create
				mockEdgeKVCreate(m, data)
				// read
				mockEdgeKVRead(m, data).Times(2)
				// read before update
				mockEdgeKVRead(m, data)

				// update retention value
				data.retention = ptr.To(88401)
				mockUpdateEdgeKVNamespace(m, data)

				// read
				mockEdgeKVRead(m, data).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", "DevExpTest:staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", "DevExpTest"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", 1234)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 86401)),
					),
				},
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/update_retention.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", "DevExpTest:staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", "DevExpTest"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", "staging"),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", 1234)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 88401)),
					),
				},
			},
		},
		"test import": {
			init: func(m *edgeworkers.Mock) {
				// 1st step: create
				mockEdgeKVCreate(m, basicData)
				// 1st step: read
				mockEdgeKVRead(m, basicData).Times(2)
				// 2nd step: import (read)
				mockEdgeKVRead(m, basicData)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
				},
				{
					ImportState:       true,
					ImportStateId:     "DevExpTest:staging",
					ResourceName:      "akamai_edgekv.test",
					ImportStateVerify: true,
				},
			},
		},
		"test import - invalid ID": {
			init: func(m *edgeworkers.Mock) {
				// create
				mockEdgeKVCreate(m, basicData)
				// read
				mockEdgeKVRead(m, basicData).Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
				},
				{
					ImportState:       true,
					ImportStateId:     "DevExpTest",
					ResourceName:      "akamai_edgekv.test",
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile("invalid EdgeKV identifier: DevExpTest"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &edgeworkers.Mock{}
			test.init(client)
			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})
			client.AssertExpectations(t)
		})
	}
}

type (
	edgeKVmockData struct {
		network     edgeworkers.NamespaceNetwork
		name        string
		geoLocation string
		retention   *int
		groupID     *int
		items       []mockItem
	}

	mockItem struct {
		key   string
		value string
		group string
	}
)

func mockEdgeKVCreate(m *edgeworkers.Mock, data edgeKVmockData) {
	mockGetEdgeKVInitializationStatus(m, "UNINITIALIZED")
	mockInitializeEdgeKV(m)
	mockGetEdgeKVInitializationStatus(m, "INITIALIZED")
	mockCreateEdgeKVNamespace(m, data)
	mockUpsertItems(m, data)
}

func mockEdgeKVRead(m *edgeworkers.Mock, data edgeKVmockData) *mock.Call {
	return m.On("GetEdgeKVNamespace", testutils.MockContext, edgeworkers.GetEdgeKVNamespaceRequest{
		Network: data.network,
		Name:    data.name,
	}).Return(&edgeworkers.Namespace{
		Name:      data.name,
		Retention: data.retention,
		GroupID:   data.groupID,
	}, nil).Once()
}

func mockGetEdgeKVInitializationStatus(m *edgeworkers.Mock, status string) {
	m.On("GetEdgeKVInitializationStatus", testutils.MockContext).Return(&edgeworkers.EdgeKVInitializationStatus{
		AccountStatus:    status,
		CPCode:           "123456",
		ProductionStatus: status,
		StagingStatus:    status,
	}, nil).Once()
}

func mockInitializeEdgeKV(m *edgeworkers.Mock) {
	m.On("InitializeEdgeKV", testutils.MockContext).Return(&edgeworkers.EdgeKVInitializationStatus{}, nil).Once()
}

func mockCreateEdgeKVNamespace(m *edgeworkers.Mock, data edgeKVmockData) {
	namespace := edgeworkers.Namespace{
		Name:        data.name,
		GeoLocation: data.geoLocation,
		Retention:   data.retention,
		GroupID:     data.groupID,
	}

	m.On("CreateEdgeKVNamespace", testutils.MockContext, edgeworkers.CreateEdgeKVNamespaceRequest{
		Network:   data.network,
		Namespace: namespace,
	}).Return(&namespace, nil).Once()
}

func mockUpsertItems(m *edgeworkers.Mock, data edgeKVmockData) {
	for _, item := range data.items {
		m.On("UpsertItem", testutils.MockContext, edgeworkers.UpsertItemRequest{
			ItemID:   item.key,
			ItemData: edgeworkers.Item(item.value),
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				Network:     edgeworkers.ItemNetwork(data.network),
				NamespaceID: data.name,
				GroupID:     item.group,
			},
		}).Return(ptr.To("OK"), nil).Once()
	}
}

func mockUpdateEdgeKVNamespace(m *edgeworkers.Mock, data edgeKVmockData) {
	m.On("UpdateEdgeKVNamespace", testutils.MockContext, edgeworkers.UpdateEdgeKVNamespaceRequest{
		Network: data.network,
		UpdateNamespace: edgeworkers.UpdateNamespace{
			Name:      data.name,
			Retention: data.retention,
			GroupID:   data.groupID,
		},
	}).Return(&edgeworkers.Namespace{
		Name:      data.name,
		Retention: data.retention,
		GroupID:   data.groupID,
	}, nil).Once()
}
