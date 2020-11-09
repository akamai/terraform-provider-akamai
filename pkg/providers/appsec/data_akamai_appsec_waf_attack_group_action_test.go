package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiWAFAttackGroupActions_data_basic(t *testing.T) {
	t.Run("match by WAFAttackGroupActions ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetWAFAttackGroupActionsResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSWAFAttackGroupActions/WAFAttackGroupActions.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetWAFAttackGroupActions",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetWAFAttackGroupActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSWAFAttackGroupActions/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_waf_attack_group_actions.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
