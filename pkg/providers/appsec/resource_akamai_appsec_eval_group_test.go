package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v7/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalGroup_res_basic(t *testing.T) {
	t.Run("match by AttackGroup ID", func(t *testing.T) {
		client := &appsec.Mock{}

		conditionExceptionJSON := testutils.LoadFixtureString(t, "testdata/TestResEvalGroup/ConditionException.json")
		conditionExceptionRawMessage := json.RawMessage(conditionExceptionJSON)

		updateResponse := appsec.UpdateAttackGroupResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalGroup/AttackGroup.json"), &updateResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAttackGroupResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalGroup/AttackGroup.json"), &getResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEvalGroup",
			testutils.MockContext,
			appsec.GetAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&getResponse, nil)

		client.On("UpdateEvalGroup",
			testutils.MockContext,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: conditionExceptionRawMessage},
		).Return(&updateResponse, nil)

		client.On("UpdateEvalGroup",
			testutils.MockContext,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "none"},
		).Return(&updateResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResEvalGroup/match_by_id.tf"),
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

		conditionExceptionJSON := testutils.LoadFixtureString(t, "testdata/TestResEvalGroup/ConditionException.json")
		conditionExceptionRawMessage := json.RawMessage(conditionExceptionJSON)

		updateResponse := appsec.UpdateAttackGroupResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalGroup/AttackGroup.json"), &updateResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAttackGroupResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalGroup/AttackGroup.json"), &getResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("UpdateEvalGroup",
			testutils.MockContext,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: conditionExceptionRawMessage},
		).Return(nil, fmt.Errorf("UpdateEvalGroup failed"))

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResEvalGroup/match_by_id.tf"),
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
