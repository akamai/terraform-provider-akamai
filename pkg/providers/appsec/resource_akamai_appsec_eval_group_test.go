package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalGroup_res_basic(t *testing.T) {
	t.Run("match by AttackGroup ID", func(t *testing.T) {
		client := &appsec.Mock{}

		conditionExceptionJSON := loadFixtureString("testdata/TestResEvalGroup/ConditionException.json")
		conditionExceptionRawMessage := json.RawMessage(conditionExceptionJSON)

		updateResponse := appsec.UpdateAttackGroupResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResEvalGroup/AttackGroup.json"), &updateResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAttackGroupResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResEvalGroup/AttackGroup.json"), &getResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEvalGroup",
			mock.Anything,
			appsec.GetAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&getResponse, nil)

		client.On("UpdateEvalGroup",
			mock.Anything,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: conditionExceptionRawMessage},
		).Return(&updateResponse, nil)

		client.On("UpdateEvalGroup",
			mock.Anything,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "none"},
		).Return(&updateResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEvalGroup/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval_group.test", "id", "43253:AAAA_81230:SQL"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiEvalGroup_res_error_updating_eval_group(t *testing.T) {
	t.Run("match by AttackGroup ID", func(t *testing.T) {
		client := &appsec.Mock{}

		conditionExceptionJSON := loadFixtureString("testdata/TestResEvalGroup/ConditionException.json")
		conditionExceptionRawMessage := json.RawMessage(conditionExceptionJSON)

		updateResponse := appsec.UpdateAttackGroupResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResEvalGroup/AttackGroup.json"), &updateResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAttackGroupResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResEvalGroup/AttackGroup.json"), &getResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("UpdateEvalGroup",
			mock.Anything,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: conditionExceptionRawMessage},
		).Return(nil, fmt.Errorf("UpdateEvalGroup failed"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestResEvalGroup/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_eval_group.test", "id", "43253:AAAA_81230:SQL"),
						),
						ExpectError: regexp.MustCompile(`UpdateEvalGroup failed`),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
