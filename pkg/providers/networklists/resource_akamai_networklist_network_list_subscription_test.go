package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/networklists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiNetworkListSubscription_res_basic(t *testing.T) {
	t.Run("match by NetworkListSubscription ID", func(t *testing.T) {
		client := &mocknetworklists{}

		cu := networklists.UpdateNetworkListSubscriptionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResNetworkListSubscription/NetworkListSubscription.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := networklists.GetNetworkListSubscriptionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResNetworkListSubscription/NetworkListSubscription.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		cd := networklists.RemoveNetworkListSubscriptionResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResNetworkListSubscription/NetworkListSubscription.json"))
		json.Unmarshal([]byte(expectJSD), &cd)

		client.On("GetNetworkListSubscription",
			mock.Anything, // ctx is irrelevant for this test
			networklists.GetNetworkListSubscriptionRequest{Recipients: []string{"test@email.com"}, UniqueIds: []string{"79536_MARTINNETWORKLIST"}}, //Recipients: "79536_MARTINNETWORKLIST", UniqueIds: "test@email.com"},
		).Return(&cr, nil)

		client.On("UpdateNetworkListSubscription",
			mock.Anything, // ctx is irrelevant for this test
			networklists.UpdateNetworkListSubscriptionRequest{Recipients: []string{"test@email.com"}, UniqueIds: []string{"79536_MARTINNETWORKLIST"}}, //Recipients: ,
		).Return(&cu, nil)

		client.On("RemoveNetworkListSubscription",
			mock.Anything, // ctx is irrelevant for this test
			networklists.RemoveNetworkListSubscriptionRequest{Recipients: []string{"test@email.com"}, UniqueIds: []string{"79536_MARTINNETWORKLIST"}}, //Recipients: ,
		).Return(&cd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResNetworkListSubscription/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_network_list_subscription.test", "id", "f7a36129f691baa1201d963b8537eb69caa28863:dd6085a7b8c8f8efaecbd420aff85a3e865ad5ca"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
