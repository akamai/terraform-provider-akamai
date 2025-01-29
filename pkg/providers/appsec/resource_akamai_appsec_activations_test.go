package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

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

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Times(3)

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
		).Return(&createActivationsResponse, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsDeleteResponse, nil)

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547694,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsDeleteResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
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

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Times(3)

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
		).Return(&createActivationsResponse, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponseDelete, nil)

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547694,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil)

		// update only note field change suppressed

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notes.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
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

	t.Run("notes field change not suppressed when other fields  changed", func(t *testing.T) {
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

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Times(3)

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
		).Return(&getActivationsUpdatedResponse, nil).Times(3)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponseDeleteUpdated, nil)

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547694,
				Action:             "DEACTIVATE",
				Network:            "PRODUCTION",
				Note:               "Test Notes update",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsUpdatedResponse, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&removeActivationsUpdatedResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
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

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Times(3)

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
		).Return(&createActivationsResponse, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponseDelete, nil)

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547694,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil)

		// Verify notification_emails field change is suppressed when nothing else changes

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user@example.com"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notification_emails.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
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

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Times(3)

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
		).Return(&getActivationsUpdatedResponse, nil).Times(3)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponseDeleteUpdated, nil)

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547694,
				Action:             "DEACTIVATE",
				Network:            "PRODUCTION",
				Note:               "Test Notes",
				NotificationEmails: []string{"user1@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsUpdatedResponse, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsResponseDelete, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&removeActivationsUpdatedResponse, nil)

		// Verify notification_emails field change is NOT suppressed when something else changes

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "notification_emails.0", "user@example.com"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notification_emails_and_network.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
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

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponse, nil).Times(3)

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

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsDeleteResponse, nil)

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547694,
				Action:             "DEACTIVATE",
				Network:            "STAGING",
				Note:               "Test Notes",
				NotificationEmails: []string{"user@example.com"},
				ActivationConfigs: []struct {
					ConfigID      int `json:"configId"`
					ConfigVersion int `json:"configVersion"`
				}{{ConfigID: 43253, ConfigVersion: 7}}},
		).Return(&removeActivationsResponse, nil)

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547695},
		).Return(&getActivationsDeleteResponse, nil).Once()

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

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
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

		// First step - create and read

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
		).Return(&getActivationsResponse, nil).Times(3)

		// Second Step : Config deactivated from UI

		client.On("GetActivations",
			testutils.MockContext,
			appsec.GetActivationsRequest{ActivationID: 547694},
		).Return(&getActivationsResponseDeactivated, nil).Once()

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
		).Return(&getActivationsResponse, nil).Times(3)

		client.On("RemoveActivations",
			testutils.MockContext,
			appsec.RemoveActivationsRequest{
				ActivationID:       547694,
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
		).Return(&getActivationsResponseDelete, nil).Times(1)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "config_id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "note", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_appsec_activations.test", "version", "7"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
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

}
