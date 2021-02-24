package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAttackGroupAction_res_basic(t *testing.T) {
	t.Run("match by AttackGroupAction ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateAttackGroupActionResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResAttackGroupAction/AttackGroupActionUpdate.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetAttackGroupActionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResAttackGroupAction/AttackGroupAction.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		client.On("GetAttackGroupAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAttackGroupActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&cr, nil)

		client.On("UpdateAttackGroupAction",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAttackGroupActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "none", Group: "SQL"},
		).Return(&cu, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAttackGroupAction/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_attack_group_action.test", "id", "43253:7:AAAA_81230:SQL"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
