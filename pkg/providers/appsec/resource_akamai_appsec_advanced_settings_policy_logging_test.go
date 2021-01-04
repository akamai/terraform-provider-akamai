package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAdvancedSettingsPolicyLogging_res_basic(t *testing.T) {
	t.Run("match by AdvancedSettingsPolicyLogging ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateAdvancedSettingsPolicyLoggingResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsPolicyLogging/AdvancedSettingsPolicyLogging.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetAdvancedSettingsPolicyLoggingResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResAdvancedSettingsPolicyLogging/AdvancedSettingsPolicyLogging.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetAdvancedSettingsPolicyLogging",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAdvancedSettingsPolicyLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdateAdvancedSettingsPolicyLogging",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAdvancedSettingsPolicyLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAdvancedSettingsPolicyLogging/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_advanced_settings_policy_logging.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
