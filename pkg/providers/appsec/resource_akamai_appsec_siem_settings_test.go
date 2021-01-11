package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiSiemSettings_res_basic(t *testing.T) {
	t.Run("match by SiemSettings ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateSiemSettingsResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResSiemSettings/SiemSettings.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetSiemSettingsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResSiemSettings/SiemSettings.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetSiemSettings",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetSiemSettingsRequest{ConfigID: 43253, Version: 7},
		).Return(&cr, nil)

		client.On("UpdateSiemSettings",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateSiemSettingsRequest{ConfigID: 43253, Version: 7, EnableForAllPolicies: false, EnableSiem: true, EnabledBotmanSiemEvents: true, SiemDefinitionID: 1, FirewallPolicyIds: []string{"12345"}},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResSiemSettings/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_siem_settings.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
