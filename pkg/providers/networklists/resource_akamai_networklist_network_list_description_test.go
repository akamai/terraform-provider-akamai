package networklists

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/networklists"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAccAkamaiNetworkListDescription_res_basic(t *testing.T) {
	t.Run("match by NetworkListDescription ID", func(t *testing.T) {
		client := &networklists.Mock{}

		cu := networklists.UpdateNetworkListDescriptionResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkListDescription/NetworkListDescription.json"), &cu)
		require.NoError(t, err)

		cr := networklists.GetNetworkListDescriptionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResNetworkListDescription/NetworkListDescription.json"), &cr)
		require.NoError(t, err)

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
				IsUnitTest:               false,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResNetworkListDescription/match_by_id.tf"),
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
