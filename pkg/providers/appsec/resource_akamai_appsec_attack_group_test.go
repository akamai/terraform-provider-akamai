package appsec

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAttackGroup_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by AttackGroup ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		conditionExceptionJSON := testutils.LoadFixtureString(t, "testdata/TestResAttackGroup/ConditionException.json")
		conditionExceptionRawMessage := json.RawMessage(conditionExceptionJSON)

		updateResponse := appsec.UpdateAttackGroupResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAttackGroup/AttackGroup.json"), &updateResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAttackGroupResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAttackGroup/AttackGroup.json"), &getResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetAttackGroup",
			testutils.MockContext,
			appsec.GetAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL"},
		).Return(&getResponse, nil)

		client.APPSEC.On("UpdateAttackGroup",
			testutils.MockContext,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: conditionExceptionRawMessage},
		).Return(&updateResponse, nil)

		client.APPSEC.On("UpdateAttackGroup",
			testutils.MockContext,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "none"},
		).Return(&updateResponse, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAttackGroup/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_attack_group.test", "id", "43253:AAAA_81230:SQL"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}

func TestAkamaiAttackGroup_res_error_updating_attack_group(t *testing.T) {
	t.Parallel()
	t.Run("match by AttackGroup ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		conditionExceptionJSON := testutils.LoadFixtureString(t, "testdata/TestResAttackGroup/ConditionException.json")
		conditionExceptionRawMessage := json.RawMessage(conditionExceptionJSON)

		updateResponse := appsec.UpdateAttackGroupResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAttackGroup/AttackGroup.json"), &updateResponse)
		require.NoError(t, err)

		getResponse := appsec.GetAttackGroupResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResAttackGroup/AttackGroup.json"), &getResponse)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("UpdateAttackGroup",
			testutils.MockContext,
			appsec.UpdateAttackGroupRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Group: "SQL", Action: "alert", JsonPayloadRaw: conditionExceptionRawMessage},
		).Return(nil, fmt.Errorf("UpdateAttackGroup failed"))

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResAttackGroup/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_attack_group.test", "id", "43253:AAAA_81230:SQL"),
					),
					ExpectError: regexp.MustCompile(`UpdateAttackGroup failed`),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})

}
