package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/edgeworkers"
)

func Test_populateEKV(t *testing.T) {
	upsertWindow = time.Millisecond
	maxUpsertAttempts = 1
	namespaceName, staging, ekvGroup := "DevExpTest", edgeworkers.ItemStagingNetwork, "greetings"
	anError, nsCreationError := "an error", "The requested namespace does not exist or namespace type is not configured for 12345"
	success := tools.StringPtr("Item was upserted in KV store with database 123456, namespace DevExpTest, group greetings, and key FR.")
	tests := map[string]struct {
		data      []interface{}
		network   edgeworkers.ItemNetwork
		init      func(*edgeworkers.Mock)
		withError error
	}{
		"no insert": {
			init: func(m *edgeworkers.Mock) {},
		},
		"one upsert": {
			network: staging,
			data: []interface{}{
				map[string]interface{}{
					"key":   "FR",
					"value": "bonjour",
					"group": ekvGroup,
				},
			},
			init: func(m *edgeworkers.Mock) {
				m.On("UpsertItem", mock.Anything, edgeworkers.UpsertItemRequest{
					ItemID:   "FR",
					ItemData: "bonjour",
					ItemsRequestParams: edgeworkers.ItemsRequestParams{
						NamespaceID: namespaceName,
						Network:     staging,
						GroupID:     ekvGroup,
					},
				}).Return(success, nil).Once()
			},
		},
		"server error on upsert": {
			network: staging,
			data: []interface{}{
				map[string]interface{}{
					"key":   "FR",
					"value": "bonjour",
					"group": ekvGroup,
				},
			},
			init: func(m *edgeworkers.Mock) {
				m.On("UpsertItem", mock.Anything, edgeworkers.UpsertItemRequest{
					ItemID:   "FR",
					ItemData: "bonjour",
					ItemsRequestParams: edgeworkers.ItemsRequestParams{
						NamespaceID: namespaceName,
						Network:     staging,
						GroupID:     ekvGroup,
					},
				}).Return(nil, errors.New(anError)).Once()
			},
			withError: fmt.Errorf(anError),
		},
		"max attempts not reached": {
			network: staging,
			data: []interface{}{
				map[string]interface{}{
					"key":   "FR",
					"value": "bonjour",
					"group": ekvGroup,
				},
			},
			init: func(m *edgeworkers.Mock) {
				maxUpsertAttempts = 2
				m.On("UpsertItem", mock.Anything, edgeworkers.UpsertItemRequest{
					ItemID:   "FR",
					ItemData: "bonjour",
					ItemsRequestParams: edgeworkers.ItemsRequestParams{
						NamespaceID: namespaceName,
						Network:     staging,
						GroupID:     ekvGroup,
					},
				}).Return(nil, errors.New(nsCreationError)).Once()
				m.On("UpsertItem", mock.Anything, edgeworkers.UpsertItemRequest{
					ItemID:   "FR",
					ItemData: "bonjour",
					ItemsRequestParams: edgeworkers.ItemsRequestParams{
						NamespaceID: namespaceName,
						Network:     staging,
						GroupID:     ekvGroup,
					},
				}).Return(success, nil).Once()
			},
		},
		"max attempts reached": {
			network: staging,
			data: []interface{}{
				map[string]interface{}{
					"key":   "FR",
					"value": "bonjour",
					"group": ekvGroup,
				},
			},
			init: func(m *edgeworkers.Mock) {
				maxUpsertAttempts = 1
				m.On("UpsertItem", mock.Anything, edgeworkers.UpsertItemRequest{
					ItemID:   "FR",
					ItemData: "bonjour",
					ItemsRequestParams: edgeworkers.ItemsRequestParams{
						NamespaceID: namespaceName,
						Network:     staging,
						GroupID:     ekvGroup,
					},
				}).Return(nil, errors.New(nsCreationError)).Once()
			},
			withError: fmt.Errorf("The requested namespace does not exist or namespace type is not configured for 12345"),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &edgeworkers.Mock{}
			test.init(client)
			err := populateEKV(context.Background(), client, test.data, &edgeworkers.Namespace{Name: namespaceName}, staging)
			client.AssertExpectations(t)
			if test.withError != nil {
				require.Error(t, err)
				require.Equal(t, err.Error(), test.withError.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestResourceEdgeKV(t *testing.T) {
	initWindow = time.Millisecond
	anError := "an error"
	initialized, inProgress := &edgeworkers.EdgeKVInitializationStatus{
		AccountStatus:    "INITIALIZED",
		CPCode:           "1027828",
		ProductionStatus: "INITIALIZED",
		StagingStatus:    "INITIALIZED",
	}, &edgeworkers.EdgeKVInitializationStatus{
		AccountStatus:    "IN PROGRESS",
		CPCode:           "123456",
		ProductionStatus: "IN PROGRESS",
		StagingStatus:    "IN PROGRESS",
	}
	ewUpsertBadRequestError := edgeworkers.Error{
		Status: http.StatusBadRequest,
		Detail: "The requested namespace does not exist or namespace type is not configured for 1193952 and teststaging.",
	}
	initStatusOneAttempt := []*edgeworkers.EdgeKVInitializationStatus{initialized}
	initNoErrorsOneAttempt := []error{nil}
	namespaceName, net, retention, retentionUpdated, groupID := "DevExpTest", "staging", tools.IntPtr(86401), tools.IntPtr(88401), tools.IntPtr(1234)
	var noData []map[string]interface{}
	id := fmt.Sprintf("%s:%s", namespaceName, net)
	tests := map[string]struct {
		init  func(*edgeworkers.Mock)
		steps []resource.TestStep
	}{
		"basic": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, noData)
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", id),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", namespaceName),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", net),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", *groupID)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", *retention)),
					),
				},
			},
		},
		"basic - retention 0": {
			init: func(m *edgeworkers.Mock) {
				// create
				retention := tools.IntPtr(0)
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, noData)
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic_retention_0.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", id),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", namespaceName),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", net),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", *groupID)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", 0)),
					),
				},
			},
		},
		"with some data to upsert": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, []map[string]interface{}{
						{"key": "es", "value": "hola", "group": "greetings"},
					})
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/ekv_with_data.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", id),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", namespaceName),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", net),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", *groupID)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", *retention)),
					),
				},
			},
		},
		"data upsert error": {
			init: func(m *edgeworkers.Mock) {
				upsertWindow = 1
				maxUpsertAttempts = 2
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, []map[string]interface{}{
						{"key": "es", "value": "hola", "group": "greetings", "error": ewUpsertBadRequestError},
						{"key": "es", "value": "hola", "group": "greetings", "error": ewUpsertBadRequestError},
					})
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/ekv_with_data.tf"),
					ExpectError: regexp.MustCompile("The requested namespace does not exist"),
				},
			},
		},
		"error in namespace initialization": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, anError, "",
					[]*edgeworkers.EdgeKVInitializationStatus{}, []error{}, noData)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					ExpectError: regexp.MustCompile(anError),
				},
			},
		},
		"namespace not initialized": {
			init: func(m *edgeworkers.Mock) {
				maxInitDuration = time.Duration(2) * time.Millisecond
				// create
				m.On("InitializeEdgeKV", mock.Anything).Return(inProgress, nil).Once()
				m.On("GetEdgeKVInitializationStatus", mock.Anything).Return(inProgress, nil)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					ExpectError: regexp.MustCompile("there was a timeout initializing"),
				},
			},
		},
		"error create namespace": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", anError,
					initStatusOneAttempt, initNoErrorsOneAttempt, noData)
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					ExpectError: regexp.MustCompile(anError),
				},
			},
		},
		"basic error read": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, noData)
				// read
				m.On("GetEdgeKVNamespace", mock.Anything, edgeworkers.GetEdgeKVNamespaceRequest{Network: edgeworkers.NamespaceNetwork(net), Name: namespaceName}).Return(nil, fmt.Errorf(anError)).Once()
			},
			steps: []resource.TestStep{
				{
					Config:      loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					ExpectError: regexp.MustCompile(anError),
				},
			},
		},
		"basic no diff no update": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, noData)
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "").Times(4)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", id),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", namespaceName),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", net),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", *groupID)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", *retention)),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", id),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", namespaceName),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", net),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", *groupID)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", *retention)),
					),
				},
			},
		},
		"ignore diff on group_id": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, noData)
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "").Times(4)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", id),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", namespaceName),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", net),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", *groupID)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", *retention)),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/update_diff_group_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", id),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", namespaceName),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", net),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", *groupID)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", *retention)),
					),
				},
			},
		},
		"basic diff retention update": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, noData)
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "").Times(3)
				// update
				stubResourceEdgeKVUpdatePhase(m, namespaceName, retentionUpdated, groupID, nil)
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retentionUpdated, groupID, "").Times(2)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", id),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", namespaceName),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", net),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", *groupID)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", *retention)),
					),
				},
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/update_retention.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", id),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", namespaceName),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", net),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", *groupID)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", *retentionUpdated)),
					),
				},
			},
		},
		"basic diff initial_data upsert attempt -> error": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, noData)
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "").Times(3)
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_edgekv.test", "id", id),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "namespace_name", namespaceName),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "network", net),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "group_id", fmt.Sprintf("%d", *groupID)),
						resource.TestCheckResourceAttr("akamai_edgekv.test", "retention_in_seconds", fmt.Sprintf("%d", *retention)),
					),
				},
				{
					Config:      loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/update_data.tf"),
					ExpectError: regexp.MustCompile("the field \"initial_data\" cannot be updated after resource creation"),
				},
			},
		},
		"test import": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, noData)
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
				// 2nd step: import
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
				},
				{
					ImportState:       true,
					ImportStateId:     id,
					ResourceName:      "akamai_edgekv.test",
					ImportStateVerify: true,
				},
			},
		},
		"test import - invalid ID": {
			init: func(m *edgeworkers.Mock) {
				// create
				stubResourceEdgeKVCreatePhase(m, namespaceName, net, retention, groupID, "", "",
					initStatusOneAttempt, initNoErrorsOneAttempt, noData)
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
				// read
				stubResourceEdgeKVReadPhase(m, namespaceName, net, retention, groupID, "")
			},
			steps: []resource.TestStep{
				{
					Config: loadFixtureString("./testdata/TestResourceEdgeWorkersEdgeKV/basic.tf"),
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
					ProviderFactories: testAccProviders,
					IsUnitTest:        true,
					Steps:             test.steps,
				})
			})
			client.AssertExpectations(t)
			maxInitDuration = time.Minute
		})
	}
}

func stubResourceEdgeKVUpdatePhase(m *edgeworkers.Mock, namespaceName string, retention, groupID *int, err *edgeworkers.Error) *mock.Call {
	call := m.On("UpdateEdgeKVNamespace", mock.Anything, mock.AnythingOfType("edgeworkers.UpdateEdgeKVNamespaceRequest"))
	if err != nil {
		return call.Return(nil, err).Once()
	}
	return call.Return(&edgeworkers.Namespace{
		Name:      namespaceName,
		Retention: retention,
		GroupID:   groupID,
	}, nil).Once()
}

func stubResourceEdgeKVReadPhase(m *edgeworkers.Mock, namespaceName, net string, retention, groupID *int, anError string) *mock.Call {
	on := m.On("GetEdgeKVNamespace", mock.Anything, edgeworkers.GetEdgeKVNamespaceRequest{Network: edgeworkers.NamespaceNetwork(net), Name: namespaceName})
	if anError != "" {
		return on.Return(nil, fmt.Errorf(anError)).Once()
	}
	return on.Return(&edgeworkers.Namespace{
		Name:      namespaceName,
		Retention: retention,
		GroupID:   groupID,
	}, nil).Once()
}

func stubResourceEdgeKVCreatePhase(m *edgeworkers.Mock, namespaceName, net string, retention, groupID *int,
	errorInit, errorCreate string, initStatus []*edgeworkers.EdgeKVInitializationStatus,
	errInitStatus []error, data []map[string]interface{}) {
	onInit := m.On("InitializeEdgeKV", mock.Anything)
	if errorInit != "" {
		onInit.Return(nil, fmt.Errorf(errorInit)).Once()
		return
	}
	onInit.Return(&edgeworkers.EdgeKVInitializationStatus{}, nil).Once()
	if len(errInitStatus) != len(initStatus) {
		panic(fmt.Sprintf("len(errInitStatus)=%d && len(errInitStatus)=%d", len(errInitStatus), len(initStatus)))
	}
	for i, err := range errInitStatus {
		status := initStatus[i]
		m.On("GetEdgeKVInitializationStatus", mock.Anything).Return(status, err).Once()
	}
	onCreate := m.On("CreateEdgeKVNamespace", mock.Anything, edgeworkers.CreateEdgeKVNamespaceRequest{
		Network: edgeworkers.NamespaceNetwork(net),
		Namespace: edgeworkers.Namespace{
			Name:      namespaceName,
			Retention: retention,
			GroupID:   groupID,
		},
	})
	if errorCreate != "" {
		onCreate.Return(nil, fmt.Errorf(errorCreate))
		return
	}
	onCreate.Return(&edgeworkers.Namespace{Name: namespaceName}, nil).Once()

	for _, item := range data {
		onUpsert := m.On("UpsertItem", mock.Anything, edgeworkers.UpsertItemRequest{
			ItemID:   getStringValue(item, "key"),
			ItemData: edgeworkers.Item(getStringValue(item, "value")),
			ItemsRequestParams: edgeworkers.ItemsRequestParams{
				NamespaceID: namespaceName,
				Network:     edgeworkers.ItemNetwork(net),
				GroupID:     getStringValue(item, "group"),
			},
		})
		if err, ok := item["error"]; ok {
			onUpsert.Return(nil, fmt.Errorf("%s: %s", edgeworkers.ErrUpsertItem, err)).Once()
		} else {
			onUpsert.Return(tools.StringPtr("OK"), nil).Once()
		}
	}
}
