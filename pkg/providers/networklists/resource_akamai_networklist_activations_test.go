package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAccAkamaiActivations_res_basic(t *testing.T) {
	t.Run("match by Activations ID", func(t *testing.T) {
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

		client.On("GetActivations",
			mock.Anything,
			networklists.GetActivationsRequest{ActivationID: 547694},
		).Return(&ga, nil)

		client.On("CreateActivations",
			mock.Anything, // ctx is irrelevant for this test
			networklists.CreateActivationsRequest{Action: "ACTIVATE", Network: "STAGING", Comments: "", NotificationRecipients: []string{"martin@email.io"}},
		).Return(&cr, nil)

		client.On("RemoveActivations",
			mock.Anything, // ctx is irrelevant for this test
			networklists.RemoveActivationsRequest{ActivationID: 547694, Action: "DEACTIVATE", Network: "STAGING", Comments: "", NotificationRecipients: []string{"martin@email.io"}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        false,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "name", "Network list test"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "TEST Notes"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails", "[user.example.com]"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResActivations/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "id", "547694"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "name", "Network list test"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "network", "STAGING"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notes", "TEST Notes updated"),
							resource.TestCheckResourceAttr("akamai_networklist_activations.test", "notification_emails", "[user.example.com]"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
