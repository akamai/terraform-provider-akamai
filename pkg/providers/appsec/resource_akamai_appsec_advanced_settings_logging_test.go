package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAdvancedSettingsLogging_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsLogging ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateAdvancedSettingsLoggingResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetAdvancedSettingsLoggingResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsLogging/AdvancedSettingsLogging.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetAdvancedSettingsLogging",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateAdvancedSettingsLogging",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAdvancedSettingsLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: "", AllowSampling: false, Cookies: struct {
				Type string "json:\"type\""
			}{Type: ""}, CustomHeaders: struct {
				Type   string   "json:\"type\""
				Values []string "json:\"values\""
			}{Type: "", Values: []string(nil)}, StandardHeaders: struct {
				Type   string   "json:\"type\""
				Values []string "json:\"values\""
			}{Type: "", Values: []string(nil)}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAdvancedSettingsLogging/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
