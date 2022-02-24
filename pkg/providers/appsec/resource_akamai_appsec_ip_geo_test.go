package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiIPGeo_res_basic(t *testing.T) {
	t.Run("match by IPGeo ID", func(t *testing.T) {
		client := &mockappsec{}

		getConfigurationResponse := appsec.GetConfigurationResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json")), &getConfigurationResponse)
		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&getConfigurationResponse, nil)

		updateIPGeoResponse := appsec.UpdateIPGeoResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResIPGeo/IPGeo.json")), &updateIPGeoResponse)
		client.On("UpdateIPGeo",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Block: "blockSpecificIPGeo", GeoControls: struct {
				BlockedIPNetworkLists struct {
					NetworkList []string "json:\"networkList\""
				} "json:\"blockedIPNetworkLists\""
			}{BlockedIPNetworkLists: struct {
				NetworkList []string "json:\"networkList\""
			}{NetworkList: []string{"40731_BMROLLOUTGEO", "44831_ECSCGEOBLACKLIST"}}}, IPControls: struct {
				AllowedIPNetworkLists struct {
					NetworkList []string "json:\"networkList\""
				} "json:\"allowedIPNetworkLists\""
				BlockedIPNetworkLists struct {
					NetworkList []string "json:\"networkList\""
				} "json:\"blockedIPNetworkLists\""
			}{AllowedIPNetworkLists: struct {
				NetworkList []string "json:\"networkList\""
			}{NetworkList: []string{"69601_ADYENPRODWHITELIST", "68762_ADYEN"}}, BlockedIPNetworkLists: struct {
				NetworkList []string "json:\"networkList\""
			}{NetworkList: []string{"49185_ADTWAFBYPASSLIST", "49181_ADTIPBLACKLIST"}}}},
		).Return(&updateIPGeoResponse, nil)

		getIPGeoResponse := appsec.GetIPGeoResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResIPGeo/IPGeo.json")), &getIPGeoResponse)
		client.On("GetIPGeo",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		client.On("GetIPGeo",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getIPGeoResponse, nil)

		updateIPGeoProtectionResponseAllProtectionsFalse := appsec.UpdateIPGeoProtectionResponse{}
		json.Unmarshal([]byte(loadFixtureBytes("testdata/TestResIPGeoProtection/PolicyProtections.json")), &updateIPGeoProtectionResponseAllProtectionsFalse)
		client.On("UpdateIPGeoProtection",
			mock.Anything,
			appsec.UpdateIPGeoProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&updateIPGeoProtectionResponseAllProtectionsFalse, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResIPGeo/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test", "id", "43253:AAAA_81230"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
