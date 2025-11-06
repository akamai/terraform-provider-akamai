package appsec

import (
	"context"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	akalog "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/log"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockLogger is a simple logger implementation for tests
type mockLogger struct{}

func (m *mockLogger) Debugf(_ string, _ ...interface{}) {}
func (m *mockLogger) Debug(_ string, _ ...interface{})  {}
func (m *mockLogger) Infof(_ string, _ ...interface{})  {}
func (m *mockLogger) Info(_ string, _ ...interface{})   {}
func (m *mockLogger) Warnf(_ string, _ ...interface{})  {}
func (m *mockLogger) Warn(_ string, _ ...interface{})   {}
func (m *mockLogger) Errorf(_ string, _ ...interface{}) {}
func (m *mockLogger) Error(_ string, _ ...interface{})  {}
func (m *mockLogger) Fatalf(_ string, _ ...interface{}) {}
func (m *mockLogger) Fatal(_ string, _ ...interface{})  {}
func (m *mockLogger) Panicf(_ string, _ ...interface{}) {}
func (m *mockLogger) Panic(_ string, _ ...interface{})  {}
func (m *mockLogger) Tracef(_ string, _ ...interface{}) {}
func (m *mockLogger) Trace(_ string, _ ...interface{})  {}
func (m *mockLogger) With(_ string, _ akalog.Fields) akalog.Interface {
	return m
}
func (m *mockLogger) WithField(_ string, _ interface{}) akalog.Interface {
	return m
}
func (m *mockLogger) WithFields(_ akalog.Fields) akalog.Interface {
	return m
}

func newMockLogger() akalog.Interface {
	return &mockLogger{}
}

func TestAkamaiActivations_res_basic(t *testing.T) {
	t.Run("match by Activations ID", func(t *testing.T) {
		client := &appsec.Mock{}

		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		getActivationsResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		getActivationsDeleteResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsDeleteResponse)
		require.NoError(t, err)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		getActivationHistoryResponseCreate := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		getActivationHistoryResponseAfter := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfter)
		require.NoError(t, err)

		createActivationsRequest := appsec.CreateActivationsRequest{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/CreateActivationsRequest.json"), &createActivationsRequest)
		require.NoError(t, err)

		// In create method

		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		// Mock GetHostMoveValidation to return no hosts to move (standard activation)
		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{},
		}
		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			createActivationsRequest,
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// In read method & delete
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Times(2)

		// In Delete method

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547693,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsDeleteResponse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("notes field change suppressed when other fields not changed", func(t *testing.T) {
		client := &appsec.Mock{}

		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		getActivationsResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		getActivationsResponseDelete := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDelete)
		require.NoError(t, err)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		getActivationHistoryResponseCreate := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		getActivationHistoryResponseAfter := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfter)
		require.NoError(t, err)

		// Mock GetHostMoveValidation to return no hosts to move (standard activation)
		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{},
		}

		// In create method
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		//In read method 3 times
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Times(3)

		// In Delete method
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Once()

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547693,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil).Once()

		// update only note field change suppressed
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notes.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("notes field change not suppressed when other fields changed", func(t *testing.T) {
		client := &appsec.Mock{}

		// old create
		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		getActivationsResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		getActivationsResponseDelete := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDelete)
		require.NoError(t, err)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		createActivationsUpdatedResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations_Production.json"), &createActivationsUpdatedResponse)
		require.NoError(t, err)

		removeActivationsUpdatedResponse := appsec.RemoveActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Deactivations_Production.json"), &removeActivationsUpdatedResponse)
		require.NoError(t, err)

		getActivationsUpdatedResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations_Production.json"), &getActivationsUpdatedResponse)
		require.NoError(t, err)

		getActivationsResponseDeleteUpdated := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Deactivations_Production.json"), &getActivationsResponseDeleteUpdated)
		require.NoError(t, err)

		getActivationHistoryResponseCreate := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		getActivationHistoryResponseAfter := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfter)
		require.NoError(t, err)

		getActivationHistoryResponseCreateProd := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreateProduction.json"), &getActivationHistoryResponseCreateProd)
		require.NoError(t, err)

		getActivationHistoryResponseAfterProd := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfterProduction.json"), &getActivationHistoryResponseAfterProd)
		require.NoError(t, err)

		createActivationsRequest := appsec.CreateActivationsRequest{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/CreateActivationsRequest.json"), &createActivationsRequest)
		require.NoError(t, err)

		// Mock GetHostMoveValidation to return no hosts to move (standard activation)
		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{},
		}

		// In create method
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			createActivationsRequest,
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// In read method 2 times, one for after create, one for before update
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Once()

		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreateProd, nil).Times(2)

		// In update method
		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("PRODUCTION"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "PRODUCTION",
				Note:               "Test Notes update",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsUpdatedResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsUpdatedResponse, nil).Once()

		// In read method after update & in delete
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfterProd, nil).Times(2)

		// In Delete method
		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547693,
				Action:             "DEACTIVATE",
				Network:            "PRODUCTION",
				Note:               "Test Notes update",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsUpdatedResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "PRODUCTION"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes update"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("notification_emails field change suppressed when other fields are not changed", func(t *testing.T) {
		// Mock TF lifecycle
		client := &appsec.Mock{}

		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		getActivationsResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		getActivationsResponseDelete := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDelete)
		require.NoError(t, err)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		getActivationHistoryResponseCreate := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		getActivationHistoryResponseAfter := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfter)
		require.NoError(t, err)

		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{}, // Empty slice means no host move needed
		}

		// In create method
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// In read method - after create (1x) + in read during update check (1x) + read after update/suppressed (1x) + in delete (1x) = 4x
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Times(4)

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547693,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil).Once()

		// Verify notification_emails field change is suppressed when nothing else changes

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user@example.com"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notification_emails.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							// Even when notification_emails changes, there is nothing to update
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user@example.com"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("notification_emails field change not suppressed when other fields are changed", func(t *testing.T) {
		// Mock TF lifecycle
		client := &appsec.Mock{}

		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		getActivationsResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		getActivationsResponseDelete := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDelete)
		require.NoError(t, err)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		createActivationsUpdatedResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations_Production.json"), &createActivationsUpdatedResponse)
		require.NoError(t, err)

		removeActivationsUpdatedResponse := appsec.RemoveActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Deactivations_Production.json"), &removeActivationsUpdatedResponse)
		require.NoError(t, err)

		getActivationsUpdatedResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations_Production.json"), &getActivationsUpdatedResponse)
		require.NoError(t, err)

		getActivationsResponseDeleteUpdated := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Deactivations_Production.json"), &getActivationsResponseDeleteUpdated)
		require.NoError(t, err)

		getActivationHistoryResponseCreate := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		getActivationHistoryResponseAfter := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfter)
		require.NoError(t, err)

		getActivationHistoryResponseCreateProd := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreateProduction.json"), &getActivationHistoryResponseCreateProd)
		require.NoError(t, err)

		getActivationHistoryResponseAfterProd := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfterProductionEmails.json"), &getActivationHistoryResponseAfterProd)
		require.NoError(t, err)

		// Mock GetHostMoveValidation for STAGING activation
		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{}, // Empty slice means no host move needed
		}

		// In create method

		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// In read method - after create
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Once()

		// Before update - check current state on PRODUCTION
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreateProd, nil).Times(2)

		// In update method
		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("PRODUCTION"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "PRODUCTION",
				Note:               "Test Notes",
				NotificationEmails: []string{"user1@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsUpdatedResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsUpdatedResponse, nil).Once()

		// In read method after update & in delete
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfterProd, nil).Times(2)

		// In Delete method
		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547693,
				Action:             "DEACTIVATE",
				Network:            "PRODUCTION",
				Note:               "Test Notes",
				NotificationEmails: []string{"user1@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsUpdatedResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil).Once()

		// Verify notification_emails field change is NOT suppressed when something else changes

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user@example.com"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notification_emails_and_network.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "PRODUCTION"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							// Since network and notification_emails changes, there is an update to the notification_emails
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user1@example.com"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("Retry create activation on 500x error", func(t *testing.T) {

		CreateActivationRetry = 10 * time.Millisecond

		err500x := &appsec.Error{StatusCode: 502}

		client := &appsec.Mock{}

		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		getActivationsResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		getActivationsDeleteResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsDeleteResponse)
		require.NoError(t, err)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		getActivationHistoryResponseCreate := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		getActivationHistoryResponseAfter := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfter)
		require.NoError(t, err)

		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{}, // Empty slice means no host move needed
		}

		// First attempt to create
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(nil, err500x).Once()

		// Retry CreateActivations succeeds
		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// In read method & delete - after successful create
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Times(2)

		// In delete method
		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547693,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsDeleteResponse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)

	})

	t.Run("reactivate config when manually deactivated from UI", func(t *testing.T) {
		client := &appsec.Mock{}

		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		getActivationsResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		getActivationsResponseDelete := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDelete)
		require.NoError(t, err)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		getActivationsUpdatedResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations_Production.json"), &getActivationsUpdatedResponse)
		require.NoError(t, err)

		getActivationsResponseDeleteUpdated := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Deactivations_Production.json"), &getActivationsResponseDeleteUpdated)
		require.NoError(t, err)

		getActivationsResponseDeactivated := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Manual_Deactivate.json"), &getActivationsResponseDeactivated)
		require.NoError(t, err)

		getActivationHistoryResponseCreate := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		getActivationHistoryResponseAfter := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfter)
		require.NoError(t, err)

		getActivationHistoryResponseDeactivated := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryDeactivated.json"), &getActivationHistoryResponseDeactivated)
		require.NoError(t, err)

		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{}, // Empty slice means no host move needed
		}

		// First step - create and read

		// In create method
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// In read method after first create
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Once()

		// Second Step : Config deactivated from UI
		// Refresh Read and Create both check activation history for deactivated status
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseDeactivated, nil).Times(2)

		// Reactivation - create activation again

		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// In read method after reactivation & in delete
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Times(2)

		// In Delete method
		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547693,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("Do not retry on 409x error", func(t *testing.T) {

		err409 := &appsec.Error{StatusCode: 409}

		client := &appsec.Mock{}

		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		getActivationsResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		getActivationsDeleteResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsDeleteResponse)
		require.NoError(t, err)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{}, // Empty slice means no host move needed
		}

		// Mock GetActivationHistory for Create - returns different version active, so version 7 needs activation
		getActivationHistoryResponse := appsec.GetActivationHistoryResponse{
			ActivationHistory: []appsec.Activation{
				{
					ActivationID:       547693,
					Version:            6, // Different version is active
					Network:            "STAGING",
					Status:             string(appsec.StatusActive),
					Notes:              "Previous Notes",
					NotificationEmails: []string{"user@example.com"},
				},
			},
		}

		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(nil, err409).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
						),
						ExpectError: regexp.MustCompile("Error: create activation failed"),
					},
				},
			})
		})

		client.AssertExpectations(t)

	})

	t.Run("notification_emails not suppressed when removing notification email from the list", func(t *testing.T) {
		// Mock TF lifecycle
		client := &appsec.Mock{}

		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		getActivationsResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		getActivationsResponseDelete := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDelete)
		require.NoError(t, err)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		createActivationsUpdatedResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations_Production.json"), &createActivationsUpdatedResponse)
		require.NoError(t, err)

		removeActivationsUpdatedResponse := appsec.RemoveActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Deactivations_Production.json"), &removeActivationsUpdatedResponse)
		require.NoError(t, err)

		getActivationsUpdatedResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations_Production.json"), &getActivationsUpdatedResponse)
		require.NoError(t, err)

		getActivationsResponseDeleteUpdated := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Deactivations_Production.json"), &getActivationsResponseDeleteUpdated)
		require.NoError(t, err)

		getActivationHistoryResponseCreate := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		getActivationHistoryResponseAfterV7 := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfterV7)
		require.NoError(t, err)

		getActivationHistoryResponseAfterV8 := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfterV8.json"), &getActivationHistoryResponseAfterV8)
		require.NoError(t, err)

		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{}, // Empty slice means no host move needed
		}

		// In create method
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user1@example.com", "user2@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// In read method after create
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfterV7, nil).Once()

		// Before update - check if version 8 needs activation (should return V7 still active)
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfterV7, nil).Times(2)

		// In update method - activate version 8
		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 8,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user1@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 8}}},
		).Return(&createActivationsUpdatedResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsUpdatedResponse, nil).Once()

		// In read method after update & in delete - should return V8 data
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfterV8, nil).Times(2)

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547693,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user1@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 8}}},
		).Return(&removeActivationsUpdatedResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil)

		// Verify notification_emails field change is NOT suppressed when removing an email from the list

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id_multiple_emails.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.#", "2"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user1@example.com"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.1", "user2@example.com"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notification_emails_remove_and_network.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "8"),
							// Since version and notification_emails change, there is an update to the notification_emails
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.#", "1"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user1@example.com"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("notification_emails change suppressed when removing email from list and other fields are not changed", func(t *testing.T) {
		// Mock TF lifecycle
		client := &appsec.Mock{}

		removeActivationsResponse := appsec.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		getActivationsResponse := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		getActivationsResponseDelete := appsec.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDelete)
		require.NoError(t, err)

		createActivationsResponse := appsec.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		getActivationHistoryResponseCreate := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		getActivationHistoryResponseAfter := appsec.GetActivationHistoryResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfterAdditionalEmails.json"), &getActivationHistoryResponseAfter)
		require.NoError(t, err)

		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{}, // Empty slice means no host move needed
		}
		// In create method
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		client.On("GetHostMoveValidation",
			testutils.MockContext,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user1@example.com", "user2@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// In read method after create (1x) + in read during update check (1x) + read after update (1x) + in delete (1x) = 4x
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Times(4)

		// In delete method
		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547693,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user1@example.com", "user2@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil).Once()

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil).Once()

		// Verify notification_emails field change is suppressed when only removing an email but keeping other fields the same

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id_multiple_emails.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.#", "2"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user1@example.com"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.1", "user2@example.com"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notification_emails_remove_only.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
							// Even when an email is removed from notification_emails, since nothing else changes,
							// the change should be suppressed and the original emails should remain
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.#", "2"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user1@example.com"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.1", "user2@example.com"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("version already active - no activation needed", func(t *testing.T) {
		// This test verifies the scenario where the requested version is already active
		// and no new activation is needed. The resource should use the existing activation
		// without making a CreateActivations API call.

		client := &appsec.Mock{}

		// Response structures
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		var removeActivationsResponse appsec.RemoveActivationsResponse
		var getActivationsResponseDelete appsec.GetActivationsResponse

		// Load test data
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDelete)
		require.NoError(t, err)

		// Mock for Create - check if version 7 is already active (returns ActivationHistoryAfter.json)
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Maybe()

		// No CreateActivations call should be made since version is already active

		// Mocks for cleanup (Delete)
		client.On("GetActivationHistory",
			testutils.MockContext,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Maybe()

		client.On("RemoveActivations",
			testutils.MockContext,
			mock.AnythingOfType("appsec.RemoveActivationsRequest"),
		).Return(&removeActivationsResponse, nil).Maybe()

		client.On("GetActivations",
			testutils.MockContext,
			mock.AnythingOfType("appsec.GetActivationsRequest"),
		).Return(&getActivationsResponseDelete, nil).Maybe()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						// Try to activate version 7 which is already active
						// Should succeed without making a CreateActivations call
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547693"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.#", "1"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user@example.com"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("different version in_progress - cannot activate new version", func(t *testing.T) {
		// Scenario: Version 6 is IN_PROGRESS, trying to activate version 7
		// Expected: Should return error - cannot activate while another version is in progress

		client := &appsec.Mock{}

		// Create in_progress activation history response
		getActivationHistoryResponseInProgress := appsec.GetActivationHistoryResponse{
			ActivationHistory: []appsec.Activation{
				{
					ActivationID:       547690,
					Version:            6,
					Network:            "STAGING",
					Status:             "ACTIVATION_IN_PROGRESS",
					Notes:              "Previous Notes",
					NotificationEmails: []string{"user@example.com"},
				},
			},
		}

		// Mock for Create - check activation history, finds version 6 in progress
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseInProgress, nil).Once()

		// No CreateActivations call should be made - should error before that

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config:      testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						ExpectError: regexp.MustCompile("cannot activate version 7 while version 6 is ACTIVATION_IN_PROGRESS on STAGING for config 43253"),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("update - version already active - no activation needed", func(t *testing.T) {
		// Scenario:
		// - Step 1: Activate version 7
		// - Step 2: Try to "update" but version 7 is still active (no changes)
		// Expected: No new activation should be created in step 2

		client := &appsec.Mock{}

		// Response structures
		var getActivationsResponse appsec.GetActivationsResponse
		var createActivationsResponse appsec.CreateActivationsResponse
		var getActivationHistoryResponseCreate appsec.GetActivationHistoryResponse
		var getActivationHistoryResponseAfter appsec.GetActivationHistoryResponse
		var removeActivationsResponse appsec.RemoveActivationsResponse
		var getActivationsResponseDelete appsec.GetActivationsResponse

		// Load test data
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfter)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDelete)
		require.NoError(t, err)

		// Mock GetHostMoveValidation response
		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{},
		}

		// Step 1: Initial create
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		client.On("GetHostMoveValidation",
			mock.Anything,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       appsec.NetworkValue("STAGING"),
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			mock.Anything,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			mock.Anything,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// Read after step 1 + before update check + after update read + in delete
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Maybe()

		// No CreateActivations call for step 2 - version already active

		// Mocks for cleanup (Delete)
		client.On("RemoveActivations",
			mock.Anything,
			mock.AnythingOfType("appsec.RemoveActivationsRequest"),
		).Return(&removeActivationsResponse, nil).Maybe()

		client.On("GetActivations",
			mock.Anything,
			mock.AnythingOfType("appsec.GetActivationsRequest"),
		).Return(&getActivationsResponseDelete, nil).Maybe()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
						),
					},
					{
						// Same config - version 7 is already active, so no update needed
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveVersion - match by network when multiple networks exist", func(t *testing.T) {
		// Scenario: Version 7 is active on both STAGING and PRODUCTION
		// Expected: When requesting STAGING, should find the STAGING activation (ID 547700)

		client := &appsec.Mock{}

		// Load test data - has version 7 active on both STAGING and PRODUCTION
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryMultipleNetworks.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547700, result.ActivationID)
			require.Equal(t, 7, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "ACTIVATED", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveVersion - match by status and network", func(t *testing.T) {
		// Scenario: Version 7 has ACTIVATED and PENDING_ACTIVATION statuses on STAGING
		// Expected: Should find only the ACTIVATED one (ID 547706), ignoring PENDING

		client := &appsec.Mock{}

		// Load test data - has version 7 with mixed statuses
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryMixedStatuses.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547706, result.ActivationID)
			require.Equal(t, 7, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "ACTIVATED", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveVersion - multiple active versions return latest", func(t *testing.T) {
		// Scenario: Version 7 and 8 are both active on STAGING (history ordered newest first)
		// Expected: Should return version 8 (ID 547702) as it's the latest

		client := &appsec.Mock{}

		// Load test data - has v7 and v8 both active on STAGING, v8 first (latest)
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryMultipleActiveVersions.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547702, result.ActivationID)
			require.Equal(t, 8, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "ACTIVATED", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveVersion - no activation match for network", func(t *testing.T) {
		// Scenario: Version 7 is active on PRODUCTION but not STAGING
		// Expected: Should return nil when searching for STAGING

		client := &appsec.Mock{}

		// Load test data - only PRODUCTION is active
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryNoStagingMatch.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.Nil(t, result, "Should return nil when no activation exists for the requested network")
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveVersion - no active status match returns nil", func(t *testing.T) {
		// Scenario: Version 7 exists on STAGING but with PENDING_DEACTIVATION status, not ACTIVE
		// Expected: Should return nil since no ACTIVE status exists

		client := &appsec.Mock{}

		// Load test data - has v7 but status is PENDING_DEACTIVATION not ACTIVE
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryNoActiveStatus.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.Nil(t, result, "Should return nil when no ACTIVE status exists")
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingVersion - finds in-progress activation", func(t *testing.T) {
		// Scenario: Version 7 is ACTIVATION_IN_PROGRESS on STAGING
		// Expected: Should find and return the in-progress activation (ID 547705)

		client := &appsec.Mock{}

		// Load test data - v7 is IN_PROGRESS
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryInProgress.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveOrPendingVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547705, result.ActivationID)
			require.Equal(t, 7, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "ACTIVATION_IN_PROGRESS", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingVersion - match by network when multiple networks exist", func(t *testing.T) {
		// Scenario: Version 7 is active on both STAGING and PRODUCTION
		// Expected: When requesting STAGING, should find the STAGING activation (ID 547700)

		client := &appsec.Mock{}

		// Load test data - has version 7 active on both STAGING and PRODUCTION
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryMultipleNetworks.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveOrPendingVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547700, result.ActivationID)
			require.Equal(t, 7, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "ACTIVATED", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingVersion - prefers active over pending", func(t *testing.T) {
		// Scenario: Version 7 has both ACTIVATED and PENDING_ACTIVATION statuses on STAGING
		// Expected: Should return ACTIVATED status (ID 547706) since it has priority

		client := &appsec.Mock{}

		// Load test data - has version 7 with mixed statuses
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryMixedStatuses.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveOrPendingVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547706, result.ActivationID)
			require.Equal(t, "ACTIVATED", string(result.Status), "Should prefer ACTIVATED status over PENDING")
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingVersion - multiple active versions return latest", func(t *testing.T) {
		// Scenario: Version 7 and 8 are both active on STAGING (history ordered newest first)
		// Expected: Should return version 8 (ID 547702) as it's the latest active

		client := &appsec.Mock{}

		// Load test data - has v7 and v8 both active on STAGING, v8 first (latest)
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryMultipleActiveVersions.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveOrPendingVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547702, result.ActivationID)
			require.Equal(t, 8, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "ACTIVATED", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingVersion - no activation match for network", func(t *testing.T) {
		// Scenario: Version 7 is active on PRODUCTION but not STAGING
		// Expected: Should return nil when searching for STAGING

		client := &appsec.Mock{}

		// Load test data - only PRODUCTION is active
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryNoStagingMatch.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveOrPendingVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.Nil(t, result, "Should return nil when no activation exists for the requested network")
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingVersion - ignores deactivation statuses", func(t *testing.T) {
		// Scenario: Version 7 exists on STAGING but with PENDING_DEACTIVATION status
		// Expected: Should return nil since deactivation statuses are not considered active or pending

		client := &appsec.Mock{}

		// Load test data - has v7 but status is PENDING_DEACTIVATION
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryNoActiveStatus.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveOrPendingVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.Nil(t, result, "Should return nil when only deactivation statuses exist")
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingVersion - finds pending activation", func(t *testing.T) {
		// Scenario: Version 7 has ACTIVATION_IN_PROGRESS status on STAGING
		// Expected: Should find and return the pending activation

		client := &appsec.Mock{}

		// Create test data with ACTIVATION_IN_PROGRESS status (already tested above but testing inline here)
		getActivationHistoryResponse := appsec.GetActivationHistoryResponse{
			ActivationHistory: []appsec.Activation{
				{
					ActivationID:       547708,
					Version:            7,
					Network:            "STAGING",
					Status:             "ACTIVATION_IN_PROGRESS",
					Notes:              "Activation in progress",
					NotificationEmails: []string{"user@example.com"},
				},
			},
		}

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveOrPendingVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547708, result.ActivationID)
			require.Equal(t, 7, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "ACTIVATION_IN_PROGRESS", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingVersion - finds new status activation", func(t *testing.T) {
		// Scenario: Version 7 has NEW status on STAGING
		// Expected: Should find and return the new activation

		client := &appsec.Mock{}

		// Create test data with NEW status
		getActivationHistoryResponse := appsec.GetActivationHistoryResponse{
			ActivationHistory: []appsec.Activation{
				{
					ActivationID:       547709,
					Version:            7,
					Network:            "STAGING",
					Status:             "NEW",
					Notes:              "New activation",
					NotificationEmails: []string{"user@example.com"},
				},
			},
		}

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			result, err := findCurrentActiveOrPendingVersion(ctx, client, 43253, "STAGING")
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547709, result.ActivationID)
			require.Equal(t, 7, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "NEW", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingDeactivation - finds pending deactivation", func(t *testing.T) {
		// Scenario: Version 7 has PENDING_DEACTIVATION status on STAGING
		// Expected: Should find and return the pending deactivation

		client := &appsec.Mock{}

		// Load test data - has v7 with PENDING_DEACTIVATION status
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryNoActiveStatus.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			logger := newMockLogger()
			result, err := findCurrentActiveOrPendingDeactivation(ctx, client, 43253, "STAGING", logger)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 7, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "PENDING_DEACTIVATION", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingDeactivation - finds deactivation in progress", func(t *testing.T) {
		// Scenario: Version 7 has DEACTIVATION_IN_PROGRESS status on STAGING
		// Expected: Should find and return the deactivation in progress

		client := &appsec.Mock{}

		// Create test data with DEACTIVATION_IN_PROGRESS status
		getActivationHistoryResponse := appsec.GetActivationHistoryResponse{
			ActivationHistory: []appsec.Activation{
				{
					ActivationID:       547710,
					Version:            7,
					Network:            "STAGING",
					Status:             "DEACTIVATION_IN_PROGRESS",
					Notes:              "Deactivation in progress",
					NotificationEmails: []string{"user@example.com"},
				},
			},
		}

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			logger := newMockLogger()
			result, err := findCurrentActiveOrPendingDeactivation(ctx, client, 43253, "STAGING", logger)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547710, result.ActivationID)
			require.Equal(t, 7, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "DEACTIVATION_IN_PROGRESS", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingDeactivation - finds active version when no pending deactivation", func(t *testing.T) {
		// Scenario: Version 7 is ACTIVE on STAGING (no pending deactivation)
		// Expected: Should return the active version

		client := &appsec.Mock{}

		// Load test data - has v7 with ACTIVE status
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryMultipleNetworks.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			logger := newMockLogger()
			result, err := findCurrentActiveOrPendingDeactivation(ctx, client, 43253, "STAGING", logger)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 7, result.Version)
			require.Equal(t, "STAGING", result.Network)
			require.Equal(t, "ACTIVATED", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingDeactivation - returns nil when no match", func(t *testing.T) {
		// Scenario: No activation exists for STAGING network
		// Expected: Should return nil

		client := &appsec.Mock{}

		// Load test data - has v7 only on PRODUCTION, not STAGING
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryNoStagingMatch.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			logger := newMockLogger()
			result, err := findCurrentActiveOrPendingDeactivation(ctx, client, 43253, "STAGING", logger)
			require.NoError(t, err)
			require.Nil(t, result, "Should return nil when no activation exists for the requested network")
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingDeactivation - prioritizes pending deactivation over active", func(t *testing.T) {
		// Scenario: Both PENDING_DEACTIVATION and ACTIVE versions exist on STAGING
		// Expected: Should return the pending deactivation (higher priority)

		client := &appsec.Mock{}

		// Create test data with both PENDING_DEACTIVATION and ACTIVE
		getActivationHistoryResponse := appsec.GetActivationHistoryResponse{
			ActivationHistory: []appsec.Activation{
				{
					ActivationID:       547711,
					Version:            8,
					Network:            "STAGING",
					Status:             "PENDING_DEACTIVATION",
					Notes:              "Pending deactivation",
					NotificationEmails: []string{"user@example.com"},
				},
				{
					ActivationID:       547693,
					Version:            7,
					Network:            "STAGING",
					Status:             "ACTIVE",
					Notes:              "Active version",
					NotificationEmails: []string{"user@example.com"},
				},
			},
		}

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			logger := newMockLogger()
			result, err := findCurrentActiveOrPendingDeactivation(ctx, client, 43253, "STAGING", logger)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 547711, result.ActivationID)
			require.Equal(t, 8, result.Version)
			require.Equal(t, "PENDING_DEACTIVATION", string(result.Status))
		})

		client.AssertExpectations(t)
	})

	t.Run("findCurrentActiveOrPendingDeactivation - ignores activation pending statuses", func(t *testing.T) {
		// Scenario: Version 7 has ACTIVATION_IN_PROGRESS status on STAGING
		// Expected: Should return nil since activation pending statuses are not relevant for deactivation

		client := &appsec.Mock{}

		// Load test data - has v7 with ACTIVATION_IN_PROGRESS
		var getActivationHistoryResponse appsec.GetActivationHistoryResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryInProgress.json"), &getActivationHistoryResponse)
		require.NoError(t, err)

		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponse, nil).Once()

		useClient(client, func() {
			ctx := context.Background()
			logger := newMockLogger()
			result, err := findCurrentActiveOrPendingDeactivation(ctx, client, 43253, "STAGING", logger)
			require.NoError(t, err)
			require.Nil(t, result, "Should return nil when only activation pending statuses exist")
		})

		client.AssertExpectations(t)
	})

	t.Run("deactivation already in progress - wait for existing deactivation", func(t *testing.T) {
		// Scenario: User creates resource, then tries to destroy while deactivation is already in progress
		// Expected: Should detect pending deactivation and wait for it instead of creating a new one

		client := &appsec.Mock{}

		// Response structures for create phase
		var getActivationsResponse appsec.GetActivationsResponse
		var createActivationsResponse appsec.CreateActivationsResponse
		var getActivationHistoryResponseCreate appsec.GetActivationHistoryResponse
		var getActivationHistoryResponseAfter appsec.GetActivationHistoryResponse

		// Response structures for delete phase
		var getActivationHistoryResponsePendingDeactivation appsec.GetActivationHistoryResponse
		var getActivationsResponsePendingDeactivation appsec.GetActivationsResponse
		var getActivationsResponseDeactivated appsec.GetActivationsResponse

		// Host move validation response
		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{},
		}

		// Load test data for create
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfter)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsPendingDeactivation.json"), &getActivationsResponsePendingDeactivation)
		require.NoError(t, err)

		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDeactivated)
		require.NoError(t, err)

		// Create test data for delete - deactivation already in progress
		getActivationHistoryResponsePendingDeactivation = appsec.GetActivationHistoryResponse{
			ActivationHistory: []appsec.Activation{
				{
					ActivationID:       547695,
					Version:            7,
					Network:            "STAGING",
					Status:             "PENDING_DEACTIVATION",
					Notes:              "Test Notes",
					NotificationEmails: []string{"user@example.com"},
				},
			},
		}

		// ===== CREATE PHASE =====

		// Mock for Create - check if version needs activation
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		// Host move validation
		client.On("GetHostMoveValidation",
			mock.Anything,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       "STAGING",
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		// Create activation
		client.On("CreateActivations",
			mock.Anything,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		// Poll activation
		client.On("GetActivations",
			mock.Anything,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// Read after create
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfter, nil).Once()

		// ===== DELETE PHASE =====

		// Check for pending deactivation - finds one already in progress
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponsePendingDeactivation, nil).Once()

		// No RemoveActivations call should be made - deactivation already exists

		// Initial lookup returns deactivated immediately (to avoid waiting in test)
		// In a real scenario, this would poll multiple times, but for testing we skip the wait
		client.On("GetActivations",
			mock.Anything,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDeactivated, nil).Maybe()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
						),
					},
					// When destroy happens, deactivation is already in progress
					// Should detect it and wait instead of creating a new deactivation
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("read finds different activation - version changed externally", func(t *testing.T) {
		// Scenario:
		// - Step 1: Terraform activates version 7 (activation ID 547694)
		// - Someone manually activates version 8 outside Terraform (activation ID 547693)
		// - Step 2: Terraform tries to activate version 8, finds it already active
		// Expected: No new activation needed, uses existing activation ID 547693

		client := &appsec.Mock{}

		// Host move validation response
		getHostMoveValidationResponse := appsec.GetHostMoveValidationResponse{
			HostsToMove: []appsec.HostToMove{},
		}

		// Load responses
		var getActivationsResponse appsec.GetActivationsResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponse)
		require.NoError(t, err)

		var createActivationsResponse appsec.CreateActivationsResponse
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &createActivationsResponse)
		require.NoError(t, err)

		var getActivationHistoryResponseCreate appsec.GetActivationHistoryResponse
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryCreate.json"), &getActivationHistoryResponseCreate)
		require.NoError(t, err)

		var getActivationHistoryResponseAfterV7 appsec.GetActivationHistoryResponse
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfter.json"), &getActivationHistoryResponseAfterV7)
		require.NoError(t, err)

		var getActivationHistoryResponseAfterV8 appsec.GetActivationHistoryResponse
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationHistoryAfterV8.json"), &getActivationHistoryResponseAfterV8)
		require.NoError(t, err)

		var removeActivationsResponse appsec.RemoveActivationsResponse
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &removeActivationsResponse)
		require.NoError(t, err)

		var getActivationsResponseDelete appsec.GetActivationsResponse
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &getActivationsResponseDelete)
		require.NoError(t, err)

		// Create response for v8 activation (ID 547693)
		var getActivationsResponseV8 appsec.GetActivationsResponse
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &getActivationsResponseV8)
		require.NoError(t, err)
		// Modify for v8
		getActivationsResponseV8.ActivationID = 547693

		// Step 1: Create activation for version 7
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseCreate, nil).Once()

		// Host move validation
		client.On("GetHostMoveValidation",
			mock.Anything,
			appsec.GetHostMoveValidationRequest{
				ConfigID:      43253,
				ConfigVersion: 7,
				Network:       "STAGING",
			},
		).Return(&getHostMoveValidationResponse, nil).Once()

		client.On("CreateActivations",
			mock.Anything,
			appsec.CreateActivationsRequest{
				Action:             "ACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&createActivationsResponse, nil).Once()

		client.On("GetActivations",
			mock.Anything,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Once()

		// Read after step 1
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfterV7, nil).Once()

		// Step 2: Update to v8
		// Before update - check current state, v8 is already active (manual activation)
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfterV8, nil).Once()

		// Read after step 2 (2nd call)
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfterV8, nil).Once()

		// No CreateActivations call for step 2 - version 8 already active

		// Cleanup mocks - Delete phase also calls GetActivationHistory
		client.On("GetActivationHistory",
			mock.Anything,
			appsec.GetActivationHistoryRequest{ConfigID: 43253},
		).Return(&getActivationHistoryResponseAfterV8, nil).Once()

		client.On("RemoveActivations",
			mock.Anything,
			mock.AnythingOfType("appsec.RemoveActivationsRequest"),
		).Return(&removeActivationsResponse, nil).Once()

		client.On("GetActivations",
			mock.Anything,
			mock.AnythingOfType("appsec.GetActivationsRequest"),
		).Return(&getActivationsResponseDelete, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
						),
					},
					{
						// Update to activate v8
						// v8 was already activated manually (activation ID 547693)
						// No new activation should be created
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id_v8.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							// Should use existing activation for v8
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "id", "547694"), // For now we do not update id field in read method. So in this case, the id remains unchanged.
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "8"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user1@example.com"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
