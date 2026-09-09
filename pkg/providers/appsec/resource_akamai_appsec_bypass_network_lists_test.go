package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiBypassNetworkLists_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by BypassNetworkLists ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
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

		client.APPSEC.On("GetWAPBypassNetworkLists",
			testutils.MockContext,
			appsec.GetWAPBypassNetworkListsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getWAPBypassNetworkListsResponse, nil)

		client.APPSEC.On("UpdateWAPBypassNetworkLists",
			testutils.MockContext,
			appsec.UpdateWAPBypassNetworkListsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", NetworkLists: []string{"1304427_AAXXBBLIST", "888518_ACDDCKERS"}},
		).Return(&updateWAPBypassNetworkListsResponse, nil)

		client.APPSEC.On("RemoveWAPBypassNetworkLists",
			testutils.MockContext,
			appsec.RemoveWAPBypassNetworkListsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", NetworkLists: []string{}},
		).Return(&removeWAPBypassNetworkListsResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResBypassNetworkLists/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_bypass_network_lists.test", "id", "43253"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
