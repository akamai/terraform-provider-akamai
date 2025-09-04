package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiRatePolicyAction_res_basic(t *testing.T) {
	t.Run("match by RatePolicyAction ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		actionAfterCreate := appsec.UpdateRatePolicyActionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicyAction/ActionAfterCreate.json"), &actionAfterCreate)
		require.NoError(t, err)

		allActionsAfterCreate := appsec.GetRatePolicyActionsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicyAction/AllActionsAfterCreate.json"), &allActionsAfterCreate)
		require.NoError(t, err)

		actionAfterUpdate := appsec.UpdateRatePolicyActionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicyAction/ActionAfterUpdate.json"), &actionAfterUpdate)
		require.NoError(t, err)

		allActionsAfterUpdate := appsec.GetRatePolicyActionsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicyAction/AllActionsAfterUpdate.json"), &allActionsAfterUpdate)
		require.NoError(t, err)

		actionAfterChallengeUpdate := appsec.UpdateRatePolicyActionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicyAction/ActionAfterChallengeUpdate.json"), &actionAfterChallengeUpdate)
		require.NoError(t, err)

		allActionsAfterChallengeUpdate := appsec.GetRatePolicyActionsResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicyAction/AllActionsAfterChallengeUpdate.json"), &allActionsAfterChallengeUpdate)
		require.NoError(t, err)

		actionAfterDelete := appsec.UpdateRatePolicyActionResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRatePolicyAction/ActionAfterDelete.json"), &actionAfterDelete)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		// Create ("none", "none")
		client.On("UpdateRatePolicyAction",
			testutils.MockContext,
			appsec.UpdateRatePolicyActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RatePolicyID: 135355, Ipv4Action: "none", Ipv6Action: "none"},
		).Return(&actionAfterCreate, nil).Once()

		// Read called from Create, returns ("none, "none") for this policy
		client.On("GetRatePolicyActions",
			testutils.MockContext,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterCreate, nil).Once()

		// Read called for TestStep 1, returns ("none", "none") for this policy
		client.On("GetRatePolicyActions",
			testutils.MockContext,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterCreate, nil).Once()

		// Read called by TF to determine whether Update will be called (diff check), returns ("none", "none") for this policy
		client.On("GetRatePolicyActions",
			testutils.MockContext,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterCreate, nil).Once()

		// Update with ("alert", "deny")
		client.On("UpdateRatePolicyAction",
			testutils.MockContext,
			appsec.UpdateRatePolicyActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RatePolicyID: 135355, Ipv4Action: "alert", Ipv6Action: "deny"},
		).Return(&actionAfterUpdate, nil).Once()

		// Read called from Update, returns ("alert", "deny") for this policy
		client.On("GetRatePolicyActions",
			testutils.MockContext,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterUpdate, nil).Once()

		// Read called for TestStep 2, returns ("alert", "deny") for this policy
		client.On("GetRatePolicyActions",
			testutils.MockContext,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterUpdate, nil).Once()

		// Read called by TF to determine whether Update will be called (diff check), returns ("alert", "deny") for this policy
		client.On("GetRatePolicyActions",
			testutils.MockContext,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterUpdate, nil).Once()

		// Update with ("challenge_1234", "challenge_5678")
		client.On("UpdateRatePolicyAction",
			testutils.MockContext,
			appsec.UpdateRatePolicyActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RatePolicyID: 135355, Ipv4Action: "challenge_1234", Ipv6Action: "challenge_5678"},
		).Return(&actionAfterChallengeUpdate, nil).Once()

		// Read called from Update, returns ("challenge_1234", "challenge_5678") for this policy
		client.On("GetRatePolicyActions",
			testutils.MockContext,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterChallengeUpdate, nil).Once()

		// Read called for TestStep 3, returns ("challenge_1234", "challenge_5678") for this policy
		client.On("GetRatePolicyActions",
			testutils.MockContext,
			appsec.GetRatePolicyActionsRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&allActionsAfterChallengeUpdate, nil).Once()

		// Delete, returns ("none", "none") and does NOT generate a Read
		client.On("UpdateRatePolicyAction",
			testutils.MockContext,
			appsec.UpdateRatePolicyActionRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", RatePolicyID: 135355, Ipv4Action: "none", Ipv6Action: "none"},
		).Return(&actionAfterDelete, nil).Once()

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResRatePolicyAction/create.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "id", "43253:AAAA_81230:135355"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv4_action", "none"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv6_action", "none"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResRatePolicyAction/update.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "id", "43253:AAAA_81230:135355"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv4_action", "alert"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv6_action", "deny"),
						),
					},
					{
						Config: testutils.LoadFixtureString(t, "testdata/TestResRatePolicyAction/update-challenge_action.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "id", "43253:AAAA_81230:135355"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv4_action", "challenge_1234"),
							resource.TestCheckResourceAttr("akamai_appsec_rate_policy_action.test", "ipv6_action", "challenge_5678"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
