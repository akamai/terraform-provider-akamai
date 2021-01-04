package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAttackGroupConditionException_data_basic(t *testing.T) {
	t.Run("match by AttackGroupConditionException ID", func(t *testing.T) {
		client := &mockappsec{}

		cv := appsec.GetAttackGroupConditionExceptionResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestDSAttackGroupConditionException/AttackGroupConditionException.json"))
		json.Unmarshal([]byte(expectJS), &cv)

		client.On("GetAttackGroupConditionException",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAttackGroupConditionExceptionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&cv, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAttackGroupConditionException/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_attack_group_condition_exception.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
