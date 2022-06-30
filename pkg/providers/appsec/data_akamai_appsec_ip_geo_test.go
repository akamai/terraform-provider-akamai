package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiIPGeo_data_basic(t *testing.T) {
	t.Run("match by IPGeo ID", func(t *testing.T) {
		client := &mockappsec{}

		getIPGeoResponse := appsec.GetIPGeoResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestDSIPGeo/IPGeo.json"), &getIPGeoResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

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
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSIPGeo/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_ip_geo.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
