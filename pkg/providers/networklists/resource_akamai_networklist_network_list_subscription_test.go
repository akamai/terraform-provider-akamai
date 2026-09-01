package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccAkamaiNetworkListSubscription_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by NetworkListSubscription ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		cu := networklists.UpdateNetworkListSubscriptionResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkListSubscription/NetworkListSubscription.json"), &cu)
		require.NoError(t, err)

		cr := networklists.GetNetworkListSubscriptionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkListSubscription/NetworkListSubscription.json"), &cr)
		require.NoError(t, err)

		cd := networklists.RemoveNetworkListSubscriptionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkListSubscription/NetworkListSubscription.json"), &cd)
		require.NoError(t, err)

		client.NetworkLists.On("GetNetworkListSubscription",
			testutils.MockContext,
			networklists.GetNetworkListSubscriptionRequest{Recipients: []string{"test@email.com"}, UniqueIDs: []string{"79536_MARTINNETWORKLIST"}},
		).Return(&cr, nil)

		client.NetworkLists.On("UpdateNetworkListSubscription",
			testutils.MockContext,
			networklists.UpdateNetworkListSubscriptionRequest{Recipients: []string{"test@email.com"}, UniqueIDs: []string{"79536_MARTINNETWORKLIST"}},
		).Return(&cu, nil)

		client.NetworkLists.On("RemoveNetworkListSubscription",
			testutils.MockContext,
			networklists.RemoveNetworkListSubscriptionRequest{Recipients: []string{"test@email.com"}, UniqueIDs: []string{"79536_MARTINNETWORKLIST"}},
		).Return(&cd, nil)

		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6SDKProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResNetworkListSubscription/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_networklist_subscription.test", "id", "f7a36129f691baa1201d963b8537eb69caa28863:dd6085a7b8c8f8efaecbd420aff85a3e865ad5ca"),
					),
				},
			},
		})

		client.NetworkLists.AssertExpectations(t)
	})

}
