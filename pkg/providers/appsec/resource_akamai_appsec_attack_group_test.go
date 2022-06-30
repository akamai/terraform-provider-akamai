package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccAkamaiAttackGroup_res_basic(t *testing.T) {
	t.Run("match by AttackGroup ID", func(t *testing.T) {
		client := &mockappsec{}

		conditionExceptionJSON := loadFixtureString("testdata/TestResAttackGroup/ConditionException.json")
		conditionExceptionRawMessage := json.RawMessage(conditionExceptionJSON)

		updateResponse := appsec.UpdateAttackGroupResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAttackGroup/AttackGroup.json"), &updateResponse)

		getResponse := appsec.GetAttackGroupResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAttackGroup/AttackGroup.json"), &getResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetAttackGroup",
			mock.Anything,
			appsec.GetAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&getResponse, nil)

		client.On("UpdateAttackGroup",
			mock.Anything,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: conditionExceptionRawMessage},
		).Return(&updateResponse, nil)

		client.On("UpdateAttackGroup",
			mock.Anything,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "none"},
		).Return(&updateResponse, nil)

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
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAccAkamaiAttackGroup_res_error_updating_attack_group(t *testing.T) {
	t.Run("match by AttackGroup ID", func(t *testing.T) {
		client := &mockappsec{}

		conditionExceptionJSON := loadFixtureString("testdata/TestResAttackGroup/ConditionException.json")
		conditionExceptionRawMessage := json.RawMessage(conditionExceptionJSON)

		updateResponse := appsec.UpdateAttackGroupResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAttackGroup/AttackGroup.json"), &updateResponse)

		getResponse := appsec.GetAttackGroupResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResAttackGroup/AttackGroup.json"), &getResponse)

		config := appsec.GetConfigurationResponse{}
		json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("UpdateAttackGroup",
			mock.Anything,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: conditionExceptionRawMessage},
		).Return(nil, fmt.Errorf("UpdateAttackGroup failed"))

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
						ExpectError: regexp.MustCompile(`UpdateAttackGroup failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
