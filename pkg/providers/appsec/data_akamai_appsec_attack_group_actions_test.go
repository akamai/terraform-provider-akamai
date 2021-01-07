package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAttackGroupActions_data_basic(t *testing.T) {
	t.Run("match by AttackGroupActions ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetAttackGroupActionsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSAttackGroupActions/AttackGroupActions.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetAttackGroupActions",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAttackGroupActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAttackGroupActions/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_attack_group_actions.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
