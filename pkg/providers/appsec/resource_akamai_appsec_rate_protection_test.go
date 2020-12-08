package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiRateProtection_res_basic(t *testing.T) {
	t.Run("match by RateProtection ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateRateProtectionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResRateProtection/RateProtectionUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetRateProtectionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResRateProtection/RateProtection.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetRateProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cr, nil)

		client.On("UpdateRateProtection",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateRateProtectionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", ApplyRateControls: true},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				PreCheck:   func() { testAccPreCheck(t) },
				IsUnitTest: false,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResRateProtection/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "enabled", "false"),
						),
						ExpectNonEmptyPlan: true,
					},
					{
						Config: loadFixtureString("testdata/TestResRateProtection/update_by_id.tf"),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "id", "43253"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_protection.test", "enabled", "false"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
