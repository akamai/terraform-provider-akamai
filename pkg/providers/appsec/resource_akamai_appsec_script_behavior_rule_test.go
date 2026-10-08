package appsec

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

type scriptBehaviorRuleTestCase struct {
	init  func(*appsec.Mock)
	steps []resource.TestStep
}

type scriptBehaviorRuleFixture struct {
	configResp     appsec.GetConfigurationResponse
	createResp     *appsec.ClientSideProtectionRuleBody
	updateResp     *appsec.ClientSideProtectionRuleBody
	listResp       appsec.ListClientSideProtectionRulesResponse
	listUpdateResp appsec.ListClientSideProtectionRulesResponse
}

func TestScriptBehaviorRuleResource(t *testing.T) {
	t.Parallel()

	for name, test := range scriptBehaviorRuleTests(t) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := edgegrid.NewTestClient()
			if test.init != nil {
				test.init(client.APPSEC)
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    test.steps,
			})

			client.APPSEC.AssertExpectations(t)
		})
	}
}

func scriptBehaviorRuleTests(t *testing.T) map[string]scriptBehaviorRuleTestCase {
	t.Helper()

	fixture := newScriptBehaviorRuleFixture()
	resourceName := "akamai_appsec_script_behavior_rule.test"

	expectedAttrs := test.AttributeBatch{
		"config_id":                 "40576",
		"security_policy_id":        "ab_86352",
		"rule_id":                   "3000601",
		"rule_name":                 "Block PCI Data Access",
		"rule_description":          "Prevents third-party scripts from accessing payment card data",
		"activity_type":             "DATA_ACCESS",
		"action":                    "deny",
		"page_condition":            "MATCHES",
		"page_values.#":             "2",
		"page_values.0":             "/checkout",
		"page_values.1":             "/billing/*",
		"script_type":               "THIRD_PARTY",
		"script_condition":          "MATCHES",
		"script_values.#":           "1",
		"script_values.0":           "analytics.example.com/*",
		"activity_values_condition": "MATCHES",
		"activity_values.#":         "2",
		"activity_values.0":         "CREDIT_CARD",
		"activity_values.1":         "PII",
	}

	return map[string]scriptBehaviorRuleTestCase{
		"create and destroy": {
			init: func(m *appsec.Mock) {
				mockScriptBehaviorRuleCreateAndDestroy(m, fixture)
			},
			steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestResScriptBehaviorRule/script_behavior_rule.tf"),
				Check:  test.NewStateChecker(resourceName).CheckEqualBatch("", expectedAttrs).Build(),
			}},
		},
		"update": {
			init: func(m *appsec.Mock) {
				mockScriptBehaviorRuleUpdate(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResScriptBehaviorRule/script_behavior_rule.tf"),
					Check:  test.NewStateChecker(resourceName).CheckEqual("action", "deny").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResScriptBehaviorRule/script_behavior_rule_update.tf"),
					Check:  test.NewStateChecker(resourceName).CheckEqual("action", "alert").Build(),
				},
			},
		},
		"import state": {
			init: func(m *appsec.Mock) {
				mockScriptBehaviorRuleImport(m, fixture)
			},
			steps: []resource.TestStep{
				{
					Config:                               testutils.LoadFixtureString(t, "testdata/TestResScriptBehaviorRule/script_behavior_rule.tf"),
					ResourceName:                         resourceName,
					ImportState:                          true,
					ImportStateId:                        "40576:ab_86352:3000601",
					ImportStateVerifyIdentifierAttribute: "rule_id",
					ImportStateCheck:                     newImportCheckerWithBatch(expectedAttrs).Build(),
				},
			},
		},
		"error on create": {
			init: func(m *appsec.Mock) {
				m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 40576}).
					Return(&fixture.configResp, nil).Once()
				m.On("GetConfigurationVersion", mock.Anything, appsec.GetConfigurationVersionRequest{ConfigID: 40576, Version: 6}).
					Return(nil, appsecError(400)).Once()
				m.On("CreateClientSideProtectionRule", mock.Anything, mock.MatchedBy(func(req appsec.CreateClientSideProtectionRuleRequest) bool {
					return req.ConfigID == 40576 && req.PolicyID == "ab_86352"
				})).Return(nil, fmt.Errorf("API error: internal server error")).Once()
			},
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestResScriptBehaviorRule/script_behavior_rule.tf"),
				ExpectError: regexp.MustCompile(`creating script behavior rule`),
			}},
		},
		"error on update": {
			init: func(m *appsec.Mock) {
				mockScriptBehaviorRuleCreateAndDestroy(m, fixture)
				// Second apply: pre-plan refresh + update error
				m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 40576}).
					Return(&fixture.configResp, nil).Times(2)
				m.On("GetConfigurationVersion", mock.Anything, appsec.GetConfigurationVersionRequest{ConfigID: 40576, Version: 6}).
					Return(nil, appsecError(400)).Once()
				m.On("ListClientSideProtectionRules", mock.Anything, listRulesReq).
					Return(&fixture.listResp, nil).Once()
				m.On("UpdateClientSideProtectionRule", mock.Anything, mock.MatchedBy(func(req appsec.UpdateClientSideProtectionRuleRequest) bool {
					return req.ConfigID == 40576 && req.RuleID == 3000601
				})).Return(nil, fmt.Errorf("API error: internal server error")).Once()
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResScriptBehaviorRule/script_behavior_rule.tf"),
					Check:  test.NewStateChecker(resourceName).CheckEqual("action", "deny").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResScriptBehaviorRule/script_behavior_rule_update.tf"),
					ExpectError: regexp.MustCompile(`updating script behavior rule`),
				},
			},
		},
	}
}

func newScriptBehaviorRuleFixture() scriptBehaviorRuleFixture {
	ruleID := int64(3000601)

	createResp := &appsec.ClientSideProtectionRuleBody{
		RuleID:       &ruleID,
		RuleName:     "Block PCI Data Access",
		Description:  "Prevents third-party scripts from accessing payment card data",
		ActivityType: "DATA_ACCESS",
		Action:       "deny",
		Conditions: appsec.ClientSideProtectionConditions{
			Page: appsec.ClientSideProtectionConditionCriteria{
				Condition: "MATCHES",
				Values:    []string{"/checkout", "/billing/*"},
			},
			Script: appsec.ClientSideProtectionScriptConditionCriteria{
				ScriptType: "THIRD_PARTY",
				Condition:  "MATCHES",
				Values:     []string{"analytics.example.com/*"},
			},
			ActivityTypeValues: appsec.ClientSideProtectionConditionCriteria{
				Condition: "MATCHES",
				Values:    []string{"CREDIT_CARD", "PII"},
			},
		},
	}

	updateResp := &appsec.ClientSideProtectionRuleBody{
		RuleID:       &ruleID,
		RuleName:     "Block PCI Data Access",
		Description:  "Prevents third-party scripts from accessing payment card data",
		ActivityType: "DATA_ACCESS",
		Action:       "alert",
		Conditions:   createResp.Conditions,
	}

	listResp := appsec.ListClientSideProtectionRulesResponse{*createResp}
	listUpdateResp := appsec.ListClientSideProtectionRulesResponse{*updateResp}

	return scriptBehaviorRuleFixture{
		configResp: appsec.GetConfigurationResponse{
			ID:                40576,
			LatestVersion:     6,
			StagingVersion:    5,
			ProductionVersion: 5,
		},
		createResp:     createResp,
		updateResp:     updateResp,
		listResp:       listResp,
		listUpdateResp: listUpdateResp,
	}
}

var listRulesReq = appsec.ListClientSideProtectionRulesRequest{
	ConfigID: 40576,
	Version:  6,
	PolicyID: "ab_86352",
}

func mockScriptBehaviorRuleCreateAndDestroy(m *appsec.Mock, fixture scriptBehaviorRuleFixture) {
	// Create: getModifiableConfigVersion
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 40576}).
		Return(&fixture.configResp, nil).Times(4)
	m.On("GetConfigurationVersion", mock.Anything, appsec.GetConfigurationVersionRequest{ConfigID: 40576, Version: 6}).
		Return(nil, appsecError(400)).Times(2)

	// Create
	m.On("CreateClientSideProtectionRule", mock.Anything, mock.MatchedBy(func(req appsec.CreateClientSideProtectionRuleRequest) bool {
		return req.ConfigID == 40576 && req.PolicyID == "ab_86352" && req.Body.RuleName == "Block PCI Data Access"
	})).Return(fixture.createResp, nil).Once()

	// Read (in Create + post-step refresh)
	m.On("ListClientSideProtectionRules", mock.Anything, listRulesReq).
		Return(&fixture.listResp, nil).Times(2)

	// Delete: getModifiableConfigVersion
	m.On("DeleteClientSideProtectionRule", mock.Anything, appsec.DeleteClientSideProtectionRuleRequest{
		ConfigID: 40576, Version: 6, PolicyID: "ab_86352", RuleID: 3000601,
	}).Return(nil).Once()
}

func mockScriptBehaviorRuleUpdate(m *appsec.Mock, fixture scriptBehaviorRuleFixture) {
	// Step 1: create + read-in-create + post-step-refresh
	// Step 2: pre-plan-refresh + update + read-in-update + post-step-refresh
	// Destroy
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 40576}).
		Return(&fixture.configResp, nil).Times(8)
	m.On("GetConfigurationVersion", mock.Anything, appsec.GetConfigurationVersionRequest{ConfigID: 40576, Version: 6}).
		Return(nil, appsecError(400)).Times(3)

	m.On("CreateClientSideProtectionRule", mock.Anything, mock.MatchedBy(func(req appsec.CreateClientSideProtectionRuleRequest) bool {
		return req.ConfigID == 40576 && req.Body.Action == "deny"
	})).Return(fixture.createResp, nil).Once()

	// Step 1: read-in-create + post-step-refresh; Step 2: pre-plan-refresh (old state)
	m.On("ListClientSideProtectionRules", mock.Anything, listRulesReq).
		Return(&fixture.listResp, nil).Times(3)

	m.On("UpdateClientSideProtectionRule", mock.Anything, mock.MatchedBy(func(req appsec.UpdateClientSideProtectionRuleRequest) bool {
		return req.ConfigID == 40576 && req.PolicyID == "ab_86352" && req.RuleID == 3000601 && req.Body.Action == "alert"
	})).Return(fixture.updateResp, nil).Once()

	// Step 2: read-in-update + post-step-refresh (new state)
	m.On("ListClientSideProtectionRules", mock.Anything, listRulesReq).
		Return(&fixture.listUpdateResp, nil).Times(2)

	m.On("DeleteClientSideProtectionRule", mock.Anything, appsec.DeleteClientSideProtectionRuleRequest{
		ConfigID: 40576, Version: 6, PolicyID: "ab_86352", RuleID: 3000601,
	}).Return(nil).Once()
}

func mockScriptBehaviorRuleImport(m *appsec.Mock, fixture scriptBehaviorRuleFixture) {
	// Import only performs Read (no create or destroy)
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 40576}).
		Return(&fixture.configResp, nil)
	m.On("ListClientSideProtectionRules", mock.Anything, listRulesReq).
		Return(&fixture.listResp, nil)
}

func appsecError(status int) *appsec.Error {
	return &appsec.Error{StatusCode: status}
}
