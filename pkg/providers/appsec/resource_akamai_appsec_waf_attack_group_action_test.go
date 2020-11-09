package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiWAFAttackGroupAction_res_basic(t *testing.T) {
	t.Run("match by WAFAttackGroupAction ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateWAFAttackGroupActionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResWAFAttackGroupAction/WAFAttackGroupAction.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetWAFAttackGroupActionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResWAFAttackGroupAction/WAFAttackGroupAction.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetWAFAttackGroupAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetWAFAttackGroupActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&cr, nil)

		client.On("UpdateWAFAttackGroupAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateWAFAttackGroupActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "alert", Group: "SQL"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResWAFAttackGroupAction/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_waf_attack_group_action.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
