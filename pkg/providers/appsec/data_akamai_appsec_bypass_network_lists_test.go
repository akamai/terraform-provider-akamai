package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiBypassNetworkLists_data_basic(t *testing.T) {
	t.Run("match by BypassNetworkLists ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getBypassNetworkListsResponse := appsec.GetWAPBypassNetworkListsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSBypassNetworkLists/BypassNetworkLists.json"), &getBypassNetworkListsResponse)
		require.NoError(t, err)

		client.On("GetWAPBypassNetworkLists",
			testutils.MockContext,
			appsec.GetWAPBypassNetworkListsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getBypassNetworkListsResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSBypassNetworkLists/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_bypass_network_lists.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
