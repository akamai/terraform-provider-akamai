package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiBypassNetworkLists_res_basic(t *testing.T) {
	t.Run("match by BypassNetworkLists ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		updateWAPBypassNetworkListsResponse := appsec.UpdateWAPBypassNetworkListsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResBypassNetworkLists/BypassNetworkLists.json"), &updateWAPBypassNetworkListsResponse)
		require.NoError(t, err)

		getWAPBypassNetworkListsResponse := appsec.GetWAPBypassNetworkListsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResBypassNetworkLists/BypassNetworkLists.json"), &getWAPBypassNetworkListsResponse)
		require.NoError(t, err)

		removeWAPBypassNetworkListsResponse := appsec.RemoveWAPBypassNetworkListsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResBypassNetworkLists/RemoveNetworkLists.json"), &removeWAPBypassNetworkListsResponse)
		require.NoError(t, err)

		client.On("GetWAPBypassNetworkLists",
			testutils.MockContext,
			appsec.GetWAPBypassNetworkListsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getWAPBypassNetworkListsResponse, nil)

		client.On("UpdateWAPBypassNetworkLists",
			testutils.MockContext,
			appsec.UpdateWAPBypassNetworkListsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", NetworkLists: []string{"1304427_AAXXBBLIST", "888518_ACDDCKERS"}},
		).Return(&updateWAPBypassNetworkListsResponse, nil)

		client.On("RemoveWAPBypassNetworkLists",
			testutils.MockContext,
			appsec.RemoveWAPBypassNetworkListsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", NetworkLists: []string{}},
		).Return(&removeWAPBypassNetworkListsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResBypassNetworkLists/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_bypass_network_lists.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
