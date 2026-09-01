package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAkamaiEvalPenaltyBox_res_basic(t *testing.T) {
	t.Parallel()
	t.Run("match by EvalPenaltyBox ID", func(t *testing.T) {
		t.Parallel()
		client := edgegrid.NewTestClient()

		updResp := appsec.UpdatePenaltyBoxResponse{}
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalPenaltyBox/PenaltyBox.json"), &updResp)
		require.NoError(t, err)

		getResp := appsec.GetPenaltyBoxResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResEvalPenaltyBox/PenaltyBox.json"), &getResp)
		require.NoError(t, err)

		config := appsec.GetConfigurationResponse{}
		err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.APPSEC.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.APPSEC.On("GetEvalPenaltyBox",
			testutils.MockContext,
			appsec.GetPenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&getResp, nil)

		client.APPSEC.On("UpdateEvalPenaltyBox",
			testutils.MockContext,
			appsec.UpdatePenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "none", PenaltyBoxProtection: false},
		).Return(&updResp, nil)

		mockGetConfigurationVersionDefault(client.APPSEC)
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box.test", "id", "43253:AAAA_81230"),
					),
				},
			},
		})

		client.APPSEC.AssertExpectations(t)
	})
}

func TestAkamaiEvalPenaltyBox_Validation(t *testing.T) {
	t.Parallel()

	setupMocks := func(client *appsec.Mock, action string, protection bool) {
		var config appsec.GetConfigurationResponse
		err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			testutils.MockContext,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		client.On("GetEvalPenaltyBox",
			testutils.MockContext,
			appsec.GetPenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230"},
		).Return(&appsec.GetPenaltyBoxResponse{Action: action, PenaltyBoxProtection: protection}, nil)

		// mock the create call
		client.On("UpdateEvalPenaltyBox",
			testutils.MockContext,
			appsec.UpdatePenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: action, PenaltyBoxProtection: protection},
		).Return(&appsec.UpdatePenaltyBoxResponse{Action: action, PenaltyBoxProtection: protection}, nil)

		// mock the delete call (resets to action=none, protection=false)
		client.On("UpdateEvalPenaltyBox",
			testutils.MockContext,
			appsec.UpdatePenaltyBoxRequest{ConfigID: 43253, Version: 7, PolicyID: "AAAA_81230", Action: "none", PenaltyBoxProtection: false},
		).Return(&appsec.UpdatePenaltyBoxResponse{Action: "none", PenaltyBoxProtection: false}, nil)
	}

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"missing required config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/missing_config_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"missing required security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"missing required penalty_box_action": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/missing_penalty_box_action.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"missing required penalty_box_protection": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/missing_penalty_box_protection.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"invalid action - block is not allowed": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/invalid_action.tf"),
					ExpectError: regexp.MustCompile(`may only contain alert, deny, deny_custom_\{custom_deny_id\}, none`),
				},
			},
		},
		"valid action - alert": {
			init: func(client *appsec.Mock) {
				setupMocks(client, "alert", true)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/action_alert.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box.test", "penalty_box_action", "alert"),
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box.test", "penalty_box_protection", "true"),
					),
				},
			},
		},
		"valid action - deny": {
			init: func(client *appsec.Mock) {
				setupMocks(client, "deny", true)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/action_deny.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box.test", "penalty_box_action", "deny"),
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box.test", "penalty_box_protection", "true"),
					),
				},
			},
		},
		"valid action - deny_custom_abc": {
			init: func(client *appsec.Mock) {
				setupMocks(client, "deny_custom_abc", true)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/action_deny_custom.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box.test", "penalty_box_action", "deny_custom_abc"),
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box.test", "penalty_box_protection", "true"),
					),
				},
			},
		},
		"valid action - none": {
			init: func(client *appsec.Mock) {
				setupMocks(client, "none", false)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResEvalPenaltyBox/match_by_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box.test", "penalty_box_action", "none"),
						resource.TestCheckResourceAttr("akamai_appsec_eval_penalty_box.test", "penalty_box_protection", "false"),
					),
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()

			if tc.init != nil {
				tc.init(client.APPSEC)
			}

			mockGetConfigurationVersionDefault(client.APPSEC)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    tc.steps,
			})
			client.APPSEC.AssertExpectations(t)
		})
	}
}
