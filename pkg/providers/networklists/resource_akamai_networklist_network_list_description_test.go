package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/networklists"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiNetworkListDescription_res_basic(t *testing.T) {
	t.Run("match by NetworkListDescription ID", func(t *testing.T) {
		client := &networklists.Mock{}

		cu := networklists.UpdateNetworkListDescriptionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResNetworkListDescription/NetworkListDescription.json"), &cu)

		cr := networklists.GetNetworkListDescriptionResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResNetworkListDescription/NetworkListDescription.json"), &cr)

		client.On("GetNetworkListDescription",
			mock.Anything,
			networklists.GetNetworkListDescriptionRequest{UniqueID: "2275_VOYAGERCALLCENTERWHITELI", Name: ""},
		).Return(&cr, nil)

		client.On("UpdateNetworkListDescription",
			mock.Anything,
			networklists.UpdateNetworkListDescriptionRequest{UniqueID: "2275_VOYAGERCALLCENTERWHITELI", Name: "Voyager Call Center Whitelist", Description: "Notes about this network list"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: false,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResNetworkListDescription/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_networklist_description.test", "id", "2275_VOYAGERCALLCENTERWHITELI"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
