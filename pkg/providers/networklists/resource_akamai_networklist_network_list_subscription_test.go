package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAccAkamaiNetworkListSubscription_res_basic(t *testing.T) {
	t.Run("match by NetworkListSubscription ID", func(t *testing.T) {
		client := &networklists.Mock{}

		cu := networklists.UpdateNetworkListSubscriptionResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkListSubscription/NetworkListSubscription.json"), &cu)
		require.NoError(t, err)

		cr := networklists.GetNetworkListSubscriptionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkListSubscription/NetworkListSubscription.json"), &cr)
		require.NoError(t, err)

		cd := networklists.RemoveNetworkListSubscriptionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkListSubscription/NetworkListSubscription.json"), &cd)
		require.NoError(t, err)

		client.On("GetNetworkListSubscription",
			testutils.MockContext,
			networklists.GetNetworkListSubscriptionRequest{Recipients: []string{"test@email.com"}, UniqueIDs: []string{"79536_MARTINNETWORKLIST"}},
		).Return(&cr, nil)

		client.On("UpdateNetworkListSubscription",
			testutils.MockContext,
			networklists.UpdateNetworkListSubscriptionRequest{Recipients: []string{"test@email.com"}, UniqueIDs: []string{"79536_MARTINNETWORKLIST"}},
		).Return(&cu, nil)

		client.On("RemoveNetworkListSubscription",
			testutils.MockContext,
			networklists.RemoveNetworkListSubscriptionRequest{Recipients: []string{"test@email.com"}, UniqueIDs: []string{"79536_MARTINNETWORKLIST"}},
		).Return(&cd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResNetworkListSubscription/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_subscription.test", "id", "f7a36129f691baa1201d963b8537eb69caa28863:dd6085a7b8c8f8efaecbd420aff85a3e865ad5ca"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
