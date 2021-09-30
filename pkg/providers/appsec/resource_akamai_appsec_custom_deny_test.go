package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiCustomDeny_res_basic(t *testing.T) {
	t.Run("match by CustomDeny ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateCustomDenyResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetCustomDenyResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDeny.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		crd := appsec.RemoveCustomDenyResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDeny.json"))
		json.Unmarshal([]byte(expectJSD), &crd)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("RemoveCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.RemoveCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&crd, nil)

		crc := appsec.CreateCustomDenyResponse{}
		expectJSC := compactJSON(loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyCreate.json"))
		json.Unmarshal([]byte(expectJSC), &crc)

		client.On("GetCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918"},
		).Return(&cr, nil)

		customDenyWithPreventBrowserCacheFalseJSON := loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheFalse.json")
		client.On("UpdateCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateCustomDenyRequest{ConfigID: 43253, Version: 7, ID: "deny_custom_622918", JsonPayloadRaw: customDenyWithPreventBrowserCacheFalseJSON},
		).Return(&cu, nil)

		customDenyWithPreventBrowserCacheTrueJSON := loadFixtureBytes("testdata/TestResCustomDeny/CustomDenyWithPreventBrowserCacheTrue.json")
		client.On("CreateCustomDeny",
			mock.Anything, // ctx is irrelevant for this test
			appsec.CreateCustomDenyRequest{ConfigID: 43253, Version: 7, JsonPayloadRaw: customDenyWithPreventBrowserCacheTrueJSON},
		).Return(&crc, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResCustomDeny/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
						),
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString("testdata/TestResCustomDeny/update_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_custom_deny.test", "id", "43253:deny_custom_622918"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
