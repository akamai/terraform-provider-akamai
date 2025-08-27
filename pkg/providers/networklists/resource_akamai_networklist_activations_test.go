package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccAkamaiActivations_res_basic(t *testing.T) {
	t.Run("create and update notes and network field in activations resource", func(t *testing.T) {
		client := &networklists.Mock{}

		cu := networklists.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &cu)
		require.NoError(t, err)

		ga := networklists.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &ga)
		require.NoError(t, err)

		cr := networklists.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &cr)
		require.NoError(t, err)

		ar := networklists.GetActivationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &ar)
		require.NoError(t, err)

		client.On("CreateActivations",
			testutils.MockContext,
			networklists.CreateActivationsRequest{UniqueID: "86093_AGEOLIST", Action: "ACTIVATE", Network: "STAGING", Comments: "Test Notes", NotificationRecipients: []string{"user@example.com"}},
		).Return(&cr, nil)

		client.On("GetActivation",
			testutils.MockContext,
			networklists.GetActivationRequest{ActivationID: 547694},
		).Return(&ar, nil)

		client.On("CreateActivations",
			testutils.MockContext,
			networklists.CreateActivationsRequest{UniqueID: "86093_AGEOLIST", Action: "ACTIVATE", Network: "PRODUCTION", Comments: "Test Notes Updated", NotificationRecipients: []string{"user@example.com"}},
		).Return(&cr, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network_list_id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails.0", "user@example.com"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "sync_point", "0"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network_list_id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "PRODUCTION"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "Test Notes Updated"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails.0", "user@example.com"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "sync_point", "1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("update notes field only in activations resource", func(t *testing.T) {
		client := &networklists.Mock{}

		cu := networklists.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &cu)
		require.NoError(t, err)

		ga := networklists.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &ga)
		require.NoError(t, err)

		cr := networklists.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &cr)
		require.NoError(t, err)

		ar := networklists.GetActivationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &ar)
		require.NoError(t, err)

		client.On("CreateActivations",
			testutils.MockContext,
			networklists.CreateActivationsRequest{UniqueID: "86093_AGEOLIST", Action: "ACTIVATE", Network: "STAGING", Comments: "Test Notes", NotificationRecipients: []string{"user@example.com"}},
		).Return(&cr, nil)

		client.On("GetActivation",
			testutils.MockContext,
			networklists.GetActivationRequest{ActivationID: 547694},
		).Return(&ar, nil)

		// update only note field change suppressed

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network_list_id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails.0", "user@example.com"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "sync_point", "0"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notes.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network_list_id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails.0", "user@example.com"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "sync_point", "0"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("notification_emails field change suppressed when other fields are not changed", func(t *testing.T) {
		// Mock TF lifecycle
		client := &networklists.Mock{}

		cu := networklists.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &cu)
		require.NoError(t, err)

		ga := networklists.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &ga)
		require.NoError(t, err)

		cr := networklists.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &cr)
		require.NoError(t, err)

		ar := networklists.GetActivationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &ar)
		require.NoError(t, err)

		client.On("CreateActivations",
			testutils.MockContext,
			networklists.CreateActivationsRequest{UniqueID: "86093_AGEOLIST", Action: "ACTIVATE", Network: "STAGING", Comments: "Test Notes", NotificationRecipients: []string{"user@example.com"}},
		).Return(&cr, nil)

		client.On("GetActivation",
			testutils.MockContext,
			networklists.GetActivationRequest{ActivationID: 547694},
		).Return(&ar, nil)

		// Verify notification_emails field change is suppressed when nothing else changes
		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network_list_id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails.0", "user@example.com"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "sync_point", "0"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notification_emails.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network_list_id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "Test Notes"),
							// Even when notification_emails changes, there is nothing to update
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails.0", "user@example.com"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "sync_point", "0"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("notification_emails field change not suppressed when other fields are changed", func(t *testing.T) {
		// Mock TF lifecycle
		client := &networklists.Mock{}

		cu := networklists.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &cu)
		require.NoError(t, err)

		ga := networklists.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &ga)
		require.NoError(t, err)

		cr := networklists.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &cr)
		require.NoError(t, err)

		ar := networklists.GetActivationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &ar)
		require.NoError(t, err)

		client.On("CreateActivations",
			testutils.MockContext,
			networklists.CreateActivationsRequest{UniqueID: "86093_AGEOLIST", Action: "ACTIVATE", Network: "STAGING", Comments: "Test Notes", NotificationRecipients: []string{"user@example.com"}},
		).Return(&cr, nil)

		client.On("GetActivation",
			testutils.MockContext,
			networklists.GetActivationRequest{ActivationID: 547694},
		).Return(&ar, nil)

		client.On("CreateActivations",
			testutils.MockContext,
			networklists.CreateActivationsRequest{UniqueID: "86093_AGEOLIST", Action: "ACTIVATE", Network: "PRODUCTION", Comments: "Test Notes", NotificationRecipients: []string{"user1@example.com"}},
		).Return(&cr, nil)

		// Verify notification_emails field change is NOT suppressed when something else changes

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network_list_id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails.0", "user@example.com"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "sync_point", "0"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_notification_emails_and_network.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network_list_id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "PRODUCTION"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "Test Notes"),
							// Since network and notification_emails changes, there is an update to the notification_emails
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails.0", "user1@example.com"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "sync_point", "0"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("Retry create activation on 500x error", func(t *testing.T) {
		client := &networklists.Mock{}

		cu := networklists.RemoveActivationsResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/ActivationsDelete.json"), &cu)
		require.NoError(t, err)

		ga := networklists.GetActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &ga)
		require.NoError(t, err)

		cr := networklists.CreateActivationsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &cr)
		require.NoError(t, err)

		ar := networklists.GetActivationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResActivations/Activations.json"), &ar)
		require.NoError(t, err)

		client.On("CreateActivations",
			testutils.MockContext,
			networklists.CreateActivationsRequest{UniqueID: "86093_AGEOLIST", Action: "ACTIVATE", Network: "STAGING", Comments: "Test Notes", NotificationRecipients: []string{"user@example.com"}},
		).Return(nil, &networklists.Error{StatusCode: 500}).Once()

		client.On("CreateActivations",
			testutils.MockContext,
			networklists.CreateActivationsRequest{UniqueID: "86093_AGEOLIST", Action: "ACTIVATE", Network: "STAGING", Comments: "Test Notes", NotificationRecipients: []string{"user@example.com"}},
		).Return(&cr, nil)

		client.On("GetActivation",
			testutils.MockContext,
			networklists.GetActivationRequest{ActivationID: 547694},
		).Return(&ar, nil)

		client.On("CreateActivations",
			testutils.MockContext,
			networklists.CreateActivationsRequest{UniqueID: "86093_AGEOLIST", Action: "ACTIVATE", Network: "PRODUCTION", Comments: "Test Notes Updated", NotificationRecipients: []string{"user@example.com"}},
		).Return(&cr, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network_list_id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "Test Notes"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails.0", "user@example.com"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "sync_point", "0"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network_list_id", "86093_AGEOLIST"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "PRODUCTION"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "Test Notes Updated"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails.0", "user@example.com"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "sync_point", "1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
