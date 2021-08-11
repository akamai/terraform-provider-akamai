package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAttackGroup_res_basic(t *testing.T) {
	t.Run("match by AttackGroup ID", func(t *testing.T) {
		client := &mockappsec{}

		cu := appsec.UpdateAttackGroupResponse{}
		expectJSU := compactJSON(loadFixtureBytes("testdata/TestResAttackGroup/AttackGroup.json"))
		json.Unmarshal([]byte(expectJSU), &cu)

		cr := appsec.GetAttackGroupResponse{}
		expectJS := compactJSON(loadFixtureBytes("testdata/TestResAttackGroup/AttackGroup.json"))
		json.Unmarshal([]byte(expectJS), &cr)

		cd := appsec.UpdateAttackGroupResponse{}
		expectJSD := compactJSON(loadFixtureBytes("testdata/TestResAttackGroup/AttackGroup.json"))
		json.Unmarshal([]byte(expectJSD), &cd)

		config := appsec.GetConfigurationResponse{}
		expectConfigs := compactJSON(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"))
		json.Unmarshal([]byte(expectConfigs), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetAttackGroup",
			mock.Anything, // ctx is irrelevant for this test
			appsec.GetAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&cr, nil)

		conditionExceptionJSON := loadFixtureBytes("testdata/TestResAttackGroup/ConditionException.json")
		client.On("UpdateAttackGroup",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: conditionExceptionJSON},
		).Return(&cu, nil)

		client.On("UpdateAttackGroup",
			mock.Anything, // ctx is irrelevant for this test
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "none"},
		).Return(&cd, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest: true,
				Providers:  testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResAttackGroup/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_attack_group.test", "id", "43253:AAAA_81230:SQL"),
						),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
