package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiIPGeo_data_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by IPGeo ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		getIPGeoResponse := appsec.GetIPGeoResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSIPGeo/IPGeo.json"), &getIPGeoResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetIPGeo",
			testutils.MockContext,
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSIPGeo/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_appsec_ip_geo.test", "id", "43253"),
						resource.TestCheckResourceAttr("data.akamai_appsec_ip_geo.test", "asn_controls.#", "1"),
						resource.TestCheckResourceAttr("data.akamai_appsec_ip_geo.test", "geo_controls.#", "1"),
						resource.TestCheckResourceAttr("data.akamai_appsec_ip_geo.test", "ip_controls.#", "1"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
