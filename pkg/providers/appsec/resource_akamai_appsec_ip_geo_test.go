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

		cu := appsec.UpdateIPGeoResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResIPGeo/IPGeo.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetIPGeoResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResIPGeo/IPGeo.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetIPGeo",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetIPGeoRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

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
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResIPGeo/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_ip_geo.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
