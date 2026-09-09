package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestWAFAIRulesResource(t *testing.T) {
	t.Parallel()

	aiRulesStatusEnabledResp := appsec.GetAIRulesStatusResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWafAIRules/AIRulesStatus.json"), &aiRulesStatusEnabledResp)
	require.NoError(t, err)

	aiRulesStatusDisabledResp := appsec.GetAIRulesStatusResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWafAIRules/AIRulesStatusDisabled.json"), &aiRulesStatusDisabledResp)
	require.NoError(t, err)

	aiRulesResp := appsec.ListAIRulesResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWafAIRules/AIRules.json"), &aiRulesResp)
	require.NoError(t, err)

	aiRulesRespDeny := aiRulesRespWithAction(aiRulesResp, 3001000, "deny")
	aiRulesRespAlert := aiRulesRespWithAction(aiRulesResp, 3001000, "alert")

	statusChecker := test.NewStateChecker("akamai_appsec_waf_ai_rules.test").
		CheckEqual("config_id", "111111").
		CheckEqual("security_policy_id", "2222_333333")

	actionChecker := test.NewStateChecker("akamai_appsec_waf_ai_rules.test").
		CheckEqual("config_id", "111111").
		CheckEqual("security_policy_id", "2222_333333").
		CheckEqual("rule_id", "3001000").
		CheckEqual("rule_description", `A SQL Injection attack consists of insertion or "injection" of a SQL query via the input data from the client to the application.`)

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"status mode - create and read ENABLED": {
			init: func(m *appsec.Mock) {
				// Create: modifiable version + update status
				mockGetConfiguration(m, 1)
				mockUpdateAIRulesStatus(m, "ENABLED", 1)
				// Post-create Read
				mockGetConfiguration(m, 1)
				mockGetAIRulesStatus(m, aiRulesStatusEnabledResp, 1)
				// Destroy
				mockGetConfiguration(m, 1)
				mockUpdateAIRulesStatusToDisabled(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_status.tf"),
					Check: statusChecker.
						CheckEqual("ai_rule_status", "ENABLED").
						Build(),
				},
			},
		},
		"status mode - update from ENABLED to DISABLED": {
			init: func(m *appsec.Mock) {
				// Step 1: Create ENABLED + post-create Read
				mockGetConfiguration(m, 1)
				mockUpdateAIRulesStatus(m, "ENABLED", 1)
				mockGetConfiguration(m, 1)
				mockGetAIRulesStatus(m, aiRulesStatusEnabledResp, 1)
				// Step 2: plan Read + Update DISABLED + post-update Read + Destroy
				mockGetConfiguration(m, 1) // plan refresh
				mockGetAIRulesStatus(m, aiRulesStatusEnabledResp, 1)
				mockGetConfiguration(m, 1) // update
				mockUpdateAIRulesStatus(m, "DISABLED", 1)
				mockGetConfiguration(m, 1) // post-update Read
				mockGetAIRulesStatus(m, aiRulesStatusDisabledResp, 1)
				mockGetConfiguration(m, 1) // Destroy
				mockUpdateAIRulesStatusToDisabled(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_status.tf"),
					Check: statusChecker.
						CheckEqual("ai_rule_status", "ENABLED").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/update_waf_ai_rules_status.tf"),
					Check: statusChecker.
						CheckEqual("ai_rule_status", "DISABLED").
						Build(),
				},
			},
		},
		"status mode - drift: out-of-band unenrollment triggers 404 on refresh": {
			init: func(m *appsec.Mock) {
				// Step 1: Create ENABLED + post-create Read
				mockGetConfiguration(m, 1)
				mockUpdateAIRulesStatus(m, "ENABLED", 1)
				mockGetConfiguration(m, 1)
				mockGetAIRulesStatus(m, aiRulesStatusEnabledResp, 1)
				// Post-step-1 plan verification Read (framework checks for drift after step 1)
				mockGetConfiguration(m, 1)
				mockGetAIRulesStatus(m, aiRulesStatusEnabledResp, 1)
				// Step 2 RefreshState: GetAIRulesStatus returns 404 (unenrolled) — resource removed from state
				mockGetConfiguration(m, 1)
				mockGetAIRulesStatus404(m, 1)
				// Post-test cleanup: framework destroys what remains in state after refresh
				mockGetConfiguration(m, 1)
				mockUpdateAIRulesStatusToDisabled(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_status.tf"),
					Check:  statusChecker.CheckEqual("ai_rule_status", "ENABLED").Build(),
				},
				{
					// RefreshState simulates out-of-band unenrollment: Read gets 404,
					// resource is removed from state, and refresh plan is non-empty.
					RefreshState:       true,
					ExpectNonEmptyPlan: true,
				},
			},
		},
		"action mode - create and read rule action": {
			init: func(m *appsec.Mock) {
				// Create: modifiable version + findAIRule (version lookup) + update action
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesForResource(m, aiRulesResp, 1)
				mockUpdateAIRuleAction(m, 3001000, 1, "deny", 1)
				// Post-create Read: findAIRule returns current action
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesForResource(m, aiRulesRespDeny, 1)
				// Destroy: reset action to none
				mockGetConfiguration(m, 1)
				mockUpdateAIRuleAction(m, 3001000, 1, "none", 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/waf_ai_rule_action.tf"),
					Check: actionChecker.
						CheckEqual("action", "deny").
						Build(),
				},
			},
		},
		"action mode - update from deny to alert": {
			init: func(m *appsec.Mock) {
				// Step 1: Create deny + post-create Read
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesForResource(m, aiRulesResp, 1)
				mockUpdateAIRuleAction(m, 3001000, 1, "deny", 1)
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesForResource(m, aiRulesRespDeny, 1)
				// Step 2: plan Read + Update alert + post-update Read
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesForResource(m, aiRulesRespDeny, 1)
				mockGetConfiguration(m, 1)
				mockUpdateAIRuleAction(m, 3001000, 1, "alert", 1)
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesForResource(m, aiRulesRespAlert, 1)
				// Destroy: reset action to none
				mockGetConfiguration(m, 1)
				mockUpdateAIRuleAction(m, 3001000, 1, "none", 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/waf_ai_rule_action.tf"),
					Check: actionChecker.
						CheckEqual("action", "deny").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/update_waf_ai_rules_action.tf"),
					Check: actionChecker.
						CheckEqual("action", "alert").
						Build(),
				},
			},
		},
		"import status mode": {
			init: func(m *appsec.Mock) {
				// ImportState: getLatestConfigVersion + GetAIRulesStatus
				mockGetConfiguration(m, 1)
				mockGetAIRulesStatus(m, aiRulesStatusEnabledResp, 1)
				// Post-import Read
				mockGetConfiguration(m, 1)
				mockGetAIRulesStatus(m, aiRulesStatusEnabledResp, 1)
				// Destroy
				mockGetConfiguration(m, 1)
				mockUpdateAIRulesStatusToDisabled(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_status.tf"),
					ImportState:   true,
					ImportStateId: "111111:2222_333333",
					ResourceName:  "akamai_appsec_waf_ai_rules.test",
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("config_id", "111111").
						CheckEqual("security_policy_id", "2222_333333").
						CheckEqual("ai_rule_status", "ENABLED").
						Build(),
					ImportStatePersist: true,
				},
			},
		},
		"import action mode": {
			init: func(m *appsec.Mock) {
				// ImportState: getLatestConfigVersion + findAIRule (ListAIRules)
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesForResource(m, aiRulesRespDeny, 1)
				// Post-import Read
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesForResource(m, aiRulesRespDeny, 1)
				// Destroy: reset action to none
				mockGetConfiguration(m, 1)
				mockUpdateAIRuleAction(m, 3001000, 1, "none", 1)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/waf_ai_rule_action.tf"),
					ImportState:   true,
					ImportStateId: "111111:2222_333333:3001000",
					ResourceName:  "akamai_appsec_waf_ai_rules.test",
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("config_id", "111111").
						CheckEqual("security_policy_id", "2222_333333").
						CheckEqual("rule_id", "3001000").
						CheckEqual("action", "deny").
						CheckEqual("rule_description", `A SQL Injection attack consists of insertion or "injection" of a SQL query via the input data from the client to the application.`).
						Build(),
					ImportStatePersist: true,
				},
			},
		},
		"import status mode - policy not enrolled": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockGetAIRulesStatus404(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_status.tf"),
					ImportState:        true,
					ImportStateId:      "111111:2222_333333",
					ResourceName:       "akamai_appsec_waf_ai_rules.test",
					ExpectError:        regexp.MustCompile("policy not enrolled in AI rules"),
					ImportStatePersist: true,
				},
			},
		},
		"import - incorrectly formatted id": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_status.tf"),
					ImportState:        true,
					ImportStateId:      "111111",
					ResourceName:       "akamai_appsec_waf_ai_rules.test",
					ExpectError:        regexp.MustCompile(`incorrectly formatted`),
					ImportStatePersist: true,
				},
			},
		},
		"import - invalid config id": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_status.tf"),
					ImportState:        true,
					ImportStateId:      "abc:2222_333333",
					ResourceName:       "akamai_appsec_waf_ai_rules.test",
					ExpectError:        regexp.MustCompile(`invalid config id`),
					ImportStatePersist: true,
				},
			},
		},
		"import action mode - rule not found": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesForResource(m, aiRulesResp, 1)
			},
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/waf_ai_rule_action.tf"),
					ImportState:        true,
					ImportStateId:      "111111:2222_333333:9999999",
					ResourceName:       "akamai_appsec_waf_ai_rules.test",
					ExpectError:        regexp.MustCompile(`AI rule 9999999 not found`),
					ImportStatePersist: true,
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_missing_config_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"validate config - conflicting ai_rule_status and rule_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_conflict_status_and_rule.tf"),
					ExpectError: regexp.MustCompile("Conflicting attributes"),
				},
			},
		},
		"validate config - missing action when rule_id is set": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_missing_action.tf"),
					ExpectError: regexp.MustCompile("Incomplete rule action configuration"),
				},
			},
		},
		"validate config - no mode specified": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_no_mode.tf"),
					ExpectError: regexp.MustCompile("Missing required attributes"),
				},
			},
		},
		"status mode - error from GetConfiguration api": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/create_waf_ai_rules_status.tf"),
					ExpectError: regexp.MustCompile("Error: fetching modifiable config version"),
				},
			},
		},
		"action mode - error from GetAIRules api during create": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesForResourceFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWafAIRules/waf_ai_rule_action.tf"),
					ExpectError: regexp.MustCompile("Error: resolving rule version"),
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

func mockGetAIRulesStatus(client *appsec.Mock, resp appsec.GetAIRulesStatusResponse, times int) {
	client.On("GetAIRulesStatus", testutils.MockContext, appsec.GetAIRulesStatusRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(&resp, nil).Times(times)
}

func mockUpdateAIRulesStatus(client *appsec.Mock, status string, times int) {
	client.On("UpdateAIRulesStatus", testutils.MockContext, appsec.UpdateAIRulesStatusRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
		Body:     appsec.UpdateAIRulesStatusRequestBody{AIRuleStatus: status},
	}).Return(&appsec.UpdateAIRulesStatusResponse{AIRuleStatus: status}, nil).Times(times)
}

func mockUpdateAIRulesStatusToDisabled(client *appsec.Mock, times int) {
	mockUpdateAIRulesStatus(client, "DISABLED", times)
}

func mockListWAFAIRulesForResource(client *appsec.Mock, resp appsec.ListAIRulesResponse, times int) {
	client.On("ListAIRules", testutils.MockContext, appsec.ListAIRulesRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(&resp, nil).Times(times)
}

func mockListWAFAIRulesForResourceFailure(client *appsec.Mock) {
	client.On("ListAIRules", testutils.MockContext, appsec.ListAIRulesRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(nil, &serverError).Once()
}

func mockUpdateAIRuleAction(client *appsec.Mock, ruleID, ruleVersionID int64, action string, times int) {
	client.On("UpdateAIRuleAction", testutils.MockContext, appsec.UpdateAIRuleActionRequest{
		ConfigID:    111111,
		Version:     2,
		PolicyID:    "2222_333333",
		RuleID:      ruleID,
		RuleVersion: ruleVersionID,
		Body:        appsec.UpdateAIRuleActionRequestBody{Action: action},
	}).Return(&appsec.UpdateAIRuleActionResponse{Action: action}, nil).Times(times)
}

func aiRulesRespWithAction(base appsec.ListAIRulesResponse, ruleID int64, action string) appsec.ListAIRulesResponse {
	result := base
	result.AIRules = make([]appsec.PolicyAIRule, len(base.AIRules))
	copy(result.AIRules, base.AIRules)
	for i := range result.AIRules {
		if result.AIRules[i].RuleID == ruleID {
			result.AIRules[i].Action = action
			break
		}
	}
	return result
}
