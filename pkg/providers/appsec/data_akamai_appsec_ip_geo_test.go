package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiIPGeo_data_basic(t *testing.T) {
	t.Run("match by IPGeo ID", func(t *testing.T) {
		client := &appsec.Mock{}

		getIPGeoResponse := appsec.GetIPGeoResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSIPGeo/IPGeo.json"), &getIPGeoResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetIPGeo",
			mock.Anything,
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestDSIPGeo/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_ip_geo.test", "id", "43253"),
							resource.TestCheckResourceAttr("data.akamai_appsec_ip_geo.test", "asn_network_lists.#", "1"),
							resource.TestCheckResourceAttr("data.akamai_appsec_ip_geo.test", "geo_network_lists.#", "1"),
							resource.TestCheckResourceAttr("data.akamai_appsec_ip_geo.test", "ip_network_lists.#", "1"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
