package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSlowPostProtectionSetting_res_basic(t *testing.T) {
	t.Run("match by SlowPostProtectionSetting ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateSlowPostProtectionSettingResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtectionSetting.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetSlowPostProtectionSettingResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResSlowPostProtectionSetting/SlowPostProtectionSetting.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetSlowPostProtectionSetting",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSlowPostProtectionSettingRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdateSlowPostProtectionSetting",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateSlowPostProtectionSettingRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", SlowRateThreshold: struct {
				Rate   int "json:\"rate\""
				Period int "json:\"period\""
			}{Rate: 10, Period: 30}, DurationThreshold: struct {
				Timeout int "json:\"timeout\""
			}{Timeout: 20}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSlowPostProtectionSetting/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_slow_post.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
