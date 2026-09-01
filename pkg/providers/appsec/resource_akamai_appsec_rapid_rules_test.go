package appsec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRapidRulesResource(t *testing.T) {
	t.Parallel()

	rapidRulesPriorState := appsec.GetRapidRulesResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRapidRules/RapidRulesResponsePriorState.json"), &rapidRulesPriorState)
	require.NoError(t, err)

	rapidRulesUpdatedState := appsec.GetRapidRulesResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRapidRules/RapidRulesResponseUpdatedState.json"), &rapidRulesUpdatedState)
	require.NoError(t, err)

	ruleConditionException := appsec.RuleConditionException{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRapidRules/ConditionException.json"), &ruleConditionException)
	require.NoError(t, err)

	rapidRulesOrderTestUpdated := appsec.GetRapidRulesResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResRapidRules/RapidRulesOrderTest.json"), &rapidRulesOrderTestUpdated)
	require.NoError(t, err)

	baseChecker := test.NewStateChecker(resourceReferenceName).
		CheckEqual("id", "111111:2222_333333").
		CheckEqual("config_id", "111111").
		CheckEqual("security_policy_id", "2222_333333").
		CheckEqual("default_action", "akamai_managed").
		CheckEqual("rule_definitions", "null")

	orderedChecker := test.NewStateChecker(resourceReferenceName).
		CheckEqual("id", "111111:2222_333333").
		CheckEqual("config_id", "111111").
		CheckEqual("security_policy_id", "2222_333333").
		CheckEqual("enabled", "true").
		CheckEqual("default_action", "deny")

	createRapidRulesInit := func(m *appsec.Mock) {
		mockGetConfiguration(m, 3)
		mockUpdateRapidRulesStatus(m, true, 1)
		mockUpdateRapidRulesDefaultAction(m, "deny")
		mockGetRapidRules(m, rapidRulesPriorState, 1)
		mockUpdateRapidRuleActionLock(m, 3000214, false, 2)
		mockUpdateRapidRuleAction(m, 3000214, "deny_custom_858380")
		mockUpdateRapidRuleException(m, 3000214, ruleConditionException)

		mockGetRapidRulesStatus(m, true, 1)
		mockGetRapidRulesDefaultAction(m, "deny", 1)
		mockGetRapidRules(m, rapidRulesUpdatedState, 1)
		mockUpdateRapidRulesStatus(m, false, 1)
	}

	createRapidRulesStateChecker := baseChecker.
		CheckEqual("default_action", "deny").
		CheckEqual("rule_definitions", toDefinitionsJSON("RuleDefinitions.json")).
		Build()

	var tests = map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"check schema - missing required attribute config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_missing_config_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"check schema - missing required attribute security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"validate config - invalid default_action configuration attribute": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_invalid_default_action.tf"),
					ExpectError: regexp.MustCompile("Error: Invalid configuration attribute"),
				},
			},
		},
		"validate config - invalid (empty string) default_action configuration attribute": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_invalid_empty_default_action.tf"),
					ExpectError: regexp.MustCompile("Error: Invalid configuration attribute"),
				},
			},
		},
		"validate config - invalid rule action configuration attribute": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_invalid_rule_action.tf"),
					ExpectError: regexp.MustCompile("Error: Invalid configuration attribute"),
				},
			},
		},
		"validate config - missing rule ID configuration attribute": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_missing_rule_id.tf"),
					ExpectError: regexp.MustCompile("Error: Invalid configuration attribute"),
				},
			},
		},
		"validate config - missing rule action lock configuration attribute": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_missing_rule_action_lock.tf"),
					ExpectError: regexp.MustCompile("Error: Invalid configuration attribute"),
				},
			},
		},
		"validate config - unknown field in rule definitions json file": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_with_invalid_rule_definitions_json.tf"),
					ExpectError: regexp.MustCompile("Error: json: unknown field \"unknownFiled\""),
				},
			},
		},
		"validate config - empty array in rule definitions json file": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_with_empty_array_rule_definitions_json.tf"),
					ExpectError: regexp.MustCompile("Error: JSON cannot be empty"),
				},
			},
		},
		"validate config - empty rule definitions json file": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_with_empty_rule_definitions_json.tf"),
					ExpectError: regexp.MustCompile("Error: Attribute cannot be empty"),
				},
			},
		},
		"create rapid rules - required only": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockGetRapidRulesDefaultAction(m, "akamai_managed", 2)

				mockGetRapidRulesStatus(m, true, 1)
				mockGetRapidRules(m, rapidRulesPriorState, 1)
				mockUpdateRapidRulesStatus(m, false, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_with_minimum_params.tf"),
					Check:  baseChecker.Build(),
				},
			},
		},
		"create rapid rules": {
			init: createRapidRulesInit,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules.tf"),
					Check:  createRapidRulesStateChecker,
				},
			},
		},
		"create rapid rules - attribute accept values from variables": {
			init: createRapidRulesInit,
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_use_vars.tf"),
					Check:  createRapidRulesStateChecker,
				},
			},
		},
		"create rapid rules - Unable to read latest config version from API": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules.tf"),
					ExpectError: regexp.MustCompile("Unable to read latest config version from API"),
				},
			},
		},
		"create rapid rules - Unable to update rapid rules status": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockUpdateRapidRulesStatusFailure(m, true)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules.tf"),
					ExpectError: regexp.MustCompile("Unable to update rapid rules status"),
				},
			},
		},
		"create rapid rules - Unable to update rapid rules default action": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockUpdateRapidRulesDefaultActionFailure(m, "deny")
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules.tf"),
					ExpectError: regexp.MustCompile("Unable to update rapid rules default action"),
				},
			},
		},
		"create rapid rules - Unable to read rapid rules": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockUpdateRapidRulesDefaultAction(m, "deny")
				mockGetRapidRulesFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules.tf"),
					ExpectError: regexp.MustCompile("Unable to read rapid rules"),
				},
			},
		},
		"create rapid rules - Update rapid rule action lock failure": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockUpdateRapidRulesDefaultAction(m, "deny")
				mockGetRapidRules(m, rapidRulesPriorState, 1)
				mockUpdateRapidRuleActionLockFailure(m, 3000214, false, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules.tf"),
					ExpectError: regexp.MustCompile("Update rapid rule action lock failure"),
				},
			},
		},
		"create rapid rules - Update rapid rule action failure": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockUpdateRapidRulesDefaultAction(m, "deny")
				mockGetRapidRules(m, rapidRulesPriorState, 1)
				mockUpdateRapidRuleActionLock(m, 3000214, false, 1)
				mockUpdateRapidRuleActionFailure(m, 3000214, "deny_custom_858380")
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules.tf"),
					ExpectError: regexp.MustCompile("Update rapid rule action failure"),
				},
			},
		},
		"create rapid rules - Update rapid rule exception failure": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockUpdateRapidRulesDefaultAction(m, "deny")
				mockGetRapidRules(m, rapidRulesPriorState, 1)
				mockUpdateRapidRuleActionLock(m, 3000214, false, 2)
				mockUpdateRapidRuleAction(m, 3000214, "deny_custom_858380")
				mockUpdateRapidRuleExceptionFailure(m, 3000214, ruleConditionException)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules.tf"),
					ExpectError: regexp.MustCompile("Update rapid rule exception failure"),
				},
			},
		},
		"create rapid rules - Unable to read rapid rules status": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockUpdateRapidRulesDefaultAction(m, "deny")
				mockGetRapidRules(m, rapidRulesPriorState, 1)
				mockUpdateRapidRuleActionLock(m, 3000214, false, 2)
				mockUpdateRapidRuleAction(m, 3000214, "deny_custom_858380")
				mockUpdateRapidRuleException(m, 3000214, ruleConditionException)

				mockGetRapidRulesStatusFailure(m)
				mockUpdateRapidRulesStatus(m, false, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules.tf"),
					ExpectError: regexp.MustCompile("Unable to read rapid rules status"),
				},
			},
		},
		"update rapid rules": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 6)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockGetRapidRulesDefaultAction(m, "akamai_managed", 3)

				mockGetRapidRulesStatus(m, true, 3)
				mockGetRapidRules(m, rapidRulesPriorState, 3)

				mockUpdateRapidRulesDefaultAction(m, "deny")
				mockUpdateRapidRuleActionLock(m, 3000214, false, 2)
				mockUpdateRapidRuleAction(m, 3000214, "deny_custom_858380")
				mockUpdateRapidRuleException(m, 3000214, ruleConditionException)

				mockGetRapidRulesDefaultAction(m, "deny", 1)
				mockGetRapidRules(m, rapidRulesUpdatedState, 1)
				mockUpdateRapidRulesStatus(m, false, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_with_minimum_params.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules.tf"),
					Check: baseChecker.
						CheckEqual("default_action", "deny").
						CheckEqual("rule_definitions", toDefinitionsJSON("RuleDefinitions.json")).
						Build(),
				},
			},
		},
		"delete rapid rules": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 4)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockGetRapidRulesDefaultAction(m, "akamai_managed", 3)

				mockGetRapidRulesStatus(m, true, 2)
				mockGetRapidRules(m, rapidRulesPriorState, 2)
				mockUpdateRapidRulesStatus(m, false, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_with_minimum_params.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResRapidRules/delete_rapid_rules.tf"),
					Check: func(s *terraform.State) error {
						_, ok := s.RootModule().Resources[resourceReferenceName]
						if ok {
							return fmt.Errorf("resource %s is still present in the Terraform state", resourceReferenceName)
						}
						return nil
					},
				},
			},
		},
		"import rapid rules": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockGetRapidRulesDefaultAction(m, "akamai_managed", 2)
				mockGetRapidRules(m, rapidRulesPriorState, 2)

				mockGetRapidRulesStatus(m, true, 1)
				mockUpdateRapidRulesStatus(m, false, 1)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_with_minimum_params.tf"),
					ImportState:   true,
					ImportStateId: "111111:2222_333333",
					ResourceName:  resourceReferenceName,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("id", "111111:2222_333333").
						CheckEqual("config_id", "111111").
						CheckEqual("security_policy_id", "2222_333333").
						CheckEqual("default_action", "akamai_managed").
						CheckEqual("rule_definitions", "[\n  {\n    \"id\": 3000998,\n    \"action\": \"alert\",\n    \"lock\": false\n  },\n  {\n    \"id\": 3000214,\n    \"action\": \"alert\",\n    \"lock\": false,\n    \"conditionException\": {\n      \"exception\": {\n        \"specificHeaderCookieParamXmlOrJsonNames\": [\n          {\n            \"names\": [\n              \"Auth\"\n            ],\n            \"selector\": \"REQUEST_HEADERS\"\n          }\n        ]\n      }\n    }\n  },\n  {\n    \"id\": 999999,\n    \"action\": \"alert\",\n    \"lock\": false\n  }\n]").
						Build(),
					ImportStatePersist: true,
				},
			},
		},
		"import rapid rules - invalid id format": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_with_minimum_params.tf"),
					ImportState:        true,
					ImportStateId:      "12345",
					ResourceName:       resourceReferenceName,
					ExpectError:        regexp.MustCompile("Error: ID '12345' incorrectly formatted: should be 'CONFIG_ID:SECURITY_POLICY_ID'"),
					ImportStatePersist: true,
				},
			},
		},
		"import rapid rules - invalid security policy id value": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResRapidRules/create_rapid_rules_with_minimum_params.tf"),
					ImportState:        true,
					ImportStateId:      "123:",
					ResourceName:       resourceReferenceName,
					ExpectError:        regexp.MustCompile("Error: invalid security policy id ''"),
					ImportStatePersist: true,
				},
			},
		},
		"update rapid rules - suppress diff on rule order": {
			init: func(m *appsec.Mock) {
				// Step 1: Create resource with original rule order
				mockGetConfiguration(m, 6)
				mockUpdateRapidRulesStatus(m, true, 1)
				mockGetRapidRulesDefaultAction(m, "deny", 4)
				mockGetRapidRulesStatus(m, true, 3)

				expectedReq := appsec.GetRapidRulesRequest{
					ConfigID: 111111,
					Version:  2,
					PolicyID: "2222_333333",
				}

				// First read returns prior state (used during Step 1 create)
				m.On("GetRapidRules",
					mock.Anything,
					expectedReq,
				).Return(&rapidRulesPriorState, nil).Once()

				// Subsequent reads return updated state with reordered rules (Step 1 final read + Step 2 reads)
				m.On("GetRapidRules",
					mock.Anything,
					expectedReq,
				).Return(&rapidRulesOrderTestUpdated, nil).Times(4)

				// Apply rule actions, locks, and exceptions (2 rules × 2 locks each = 4)
				// Called During Read and Update in Step 1 and Step 2
				m.On("UpdateRapidRuleActionLock",
					mock.Anything,
					mock.AnythingOfType("appsec.UpdateRapidRuleActionLockRequest"),
				).Return(&appsec.UpdateRapidRuleActionLockResponse{
					Enabled: false,
				}, nil).Times(8)

				// Update actions for 2 rules
				// Called During Read and Update in Step 1 and Step 2
				m.On("UpdateRapidRuleAction",
					mock.Anything,
					mock.AnythingOfType("appsec.UpdateRapidRuleActionRequest"),
				).Return(&appsec.UpdateRapidRuleActionResponse{}, nil).Times(4)

				// Apply exceptions for 2 rules
				// Called During Read and Update in Step 1 and Step 2
				m.On("UpdateRapidRuleException",
					mock.Anything,
					mock.AnythingOfType("appsec.UpdateRapidRuleExceptionRequest"),
				).Return(
					(*appsec.UpdateRapidRuleExceptionResponse)(&ruleConditionException),
					nil,
				).Times(2)

				// Final status cleanup.
				mockUpdateRapidRulesStatus(m, false, 1)
			},
			steps: []resource.TestStep{
				{
					// Step 1: canonical order.
					Config: testutils.LoadFixtureString(t,
						"testdata/TestResRapidRules/create_rapid_rules_test_reordered_list.tf"),
					Check: orderedChecker.Build(),
				},
				{
					// Step 2: same rules, shuffled order. PreventJsonReorder
					// NOTE: ExpectNonEmptyPlan=true due to cosmetic 'id' diff.
					// Terraform marks computed attributes as "to be computed" during any
					// update, even when the plan modifier successfully prevents the real diff.
					Config: testutils.LoadFixtureString(t,
						"testdata/TestResRapidRules/update_rapid_rules_reordered.tf"),
					ExpectNonEmptyPlan: true,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue(
								resourceReferenceName,
								tfjsonpath.New("rule_definitions"),
								knownvalue.StringExact(
									toDefinitionsJSON("RuleDefinitionsOrderTestOriginal.json"),
								),
							),
							// Verify other attributes also didn't change
							plancheck.ExpectKnownValue(
								resourceReferenceName,
								tfjsonpath.New("enabled"),
								knownvalue.Bool(true),
							),
							plancheck.ExpectKnownValue(
								resourceReferenceName,
								tfjsonpath.New("default_action"),
								knownvalue.StringExact("deny"),
							),

							// Only 'id' should be unknown (cosmetic change)
							plancheck.ExpectUnknownValue(
								resourceReferenceName,
								tfjsonpath.New("id"),
							),
						},
					},

					Check: orderedChecker.
						CheckEqual("rule_definitions",
							toDefinitionsJSON("RuleDefinitionsOrderTestOriginal.json")).
						Build(),
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client := edgegrid.NewTestClient()

			if test.init != nil {
				test.init(client.APPSEC)
			}

			mockGetConfigurationVersionDefault(client.APPSEC)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    test.steps,
			})
			client.APPSEC.AssertExpectations(t)
		})
	}
}

func mockGetConfiguration(client *appsec.Mock, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 111111}).
		Return(&appsec.GetConfigurationResponse{
			FileType:          "RBAC",
			ID:                43253,
			LatestVersion:     2,
			Name:              "Akamai Tools",
			ProductionVersion: 1,
			StagingVersion:    1,
			TargetProduct:     "KSD",
		}, nil).Times(times)
}

func mockGetConfigurationFailure(client *appsec.Mock, times int) {
	client.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 111111}).
		Return(nil, &serverError).Times(times)
}

func mockUpdateRapidRulesStatus(client *appsec.Mock, enabled bool, times int) {
	client.On("UpdateRapidRulesStatus", mock.Anything, appsec.UpdateRapidRulesStatusRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
		Body:     appsec.UpdateRapidRulesStatusRequestBody{Enabled: ptr.To(enabled)},
	}).Return(&appsec.UpdateRapidRulesStatusResponse{Enabled: true}, nil).Times(times)
}

func mockUpdateRapidRulesStatusFailure(client *appsec.Mock, enabled bool) {
	client.On("UpdateRapidRulesStatus", mock.Anything, appsec.UpdateRapidRulesStatusRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
		Body:     appsec.UpdateRapidRulesStatusRequestBody{Enabled: ptr.To(enabled)},
	}).Return(nil, &serverError).Once()
}

func mockGetRapidRulesDefaultAction(client *appsec.Mock, action string, times int) {
	client.On("GetRapidRulesDefaultAction", mock.Anything, appsec.GetRapidRulesDefaultActionRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(&appsec.GetRapidRulesDefaultActionResponse{
		Action: action,
	}, nil).Times(times)
}

func mockUpdateRapidRulesDefaultAction(client *appsec.Mock, action string) {
	client.On("UpdateRapidRulesDefaultAction", mock.Anything, appsec.UpdateRapidRulesDefaultActionRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
		Body:     appsec.UpdateRapidRulesDefaultActionRequestBody{Action: action},
	}).Return(&appsec.UpdateRapidRulesDefaultActionResponse{
		Action: "akamai_managed",
	}, nil).Once()
}

func mockUpdateRapidRulesDefaultActionFailure(client *appsec.Mock, action string) {
	client.On("UpdateRapidRulesDefaultAction", mock.Anything, appsec.UpdateRapidRulesDefaultActionRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
		Body:     appsec.UpdateRapidRulesDefaultActionRequestBody{Action: action},
	}).Return(nil, &serverError).Once()
}

func mockGetRapidRulesStatus(client *appsec.Mock, enabled bool, times int) {
	client.On("GetRapidRulesStatus", mock.Anything, appsec.GetRapidRulesStatusRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(&appsec.GetRapidRulesStatusResponse{
		Enabled: enabled,
	}, nil).Times(times)
}

func mockGetRapidRulesStatusFailure(client *appsec.Mock) {
	client.On("GetRapidRulesStatus", mock.Anything, appsec.GetRapidRulesStatusRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(nil, &serverError).Once()
}

func mockGetRapidRules(client *appsec.Mock, resp appsec.GetRapidRulesResponse, times int) {
	client.On("GetRapidRules", mock.Anything, appsec.GetRapidRulesRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(&resp, nil).Times(times)
}

func mockGetRapidRulesFailure(client *appsec.Mock) {
	client.On("GetRapidRules", mock.Anything, appsec.GetRapidRulesRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(nil, &serverError).Once()
}

func mockUpdateRapidRuleActionLock(client *appsec.Mock, ruleID int64, enabled bool, times int) {
	client.On("UpdateRapidRuleActionLock", mock.Anything, appsec.UpdateRapidRuleActionLockRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
		RuleID:   ruleID,
		Body:     appsec.UpdateRapidRuleActionLockRequestBody{Enabled: ptr.To(enabled)},
	}).Return(&appsec.UpdateRapidRuleActionLockResponse{
		Enabled: enabled,
	}, nil).Times(times)
}

func mockUpdateRapidRuleActionLockFailure(client *appsec.Mock, ruleID int64, enabled bool, times int) {
	client.On("UpdateRapidRuleActionLock", mock.Anything, appsec.UpdateRapidRuleActionLockRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
		RuleID:   ruleID,
		Body:     appsec.UpdateRapidRuleActionLockRequestBody{Enabled: ptr.To(enabled)},
	}).Return(nil, &serverError).Times(times)
}

func mockUpdateRapidRuleAction(client *appsec.Mock, ruleID int64, action string) {
	client.On("UpdateRapidRuleAction", mock.Anything, appsec.UpdateRapidRuleActionRequest{
		ConfigID:    111111,
		Version:     2,
		PolicyID:    "2222_333333",
		RuleID:      ruleID,
		RuleVersion: 3,
		Body:        appsec.UpdateRapidRuleActionRequestBody{Action: action},
	}).Return(&appsec.UpdateRapidRuleActionResponse{
		Action: action,
		Lock:   true,
	}, nil).Once()
}

func mockUpdateRapidRuleActionFailure(client *appsec.Mock, ruleID int64, action string) {
	client.On("UpdateRapidRuleAction", mock.Anything, appsec.UpdateRapidRuleActionRequest{
		ConfigID:    111111,
		Version:     2,
		PolicyID:    "2222_333333",
		RuleID:      ruleID,
		RuleVersion: 3,
		Body:        appsec.UpdateRapidRuleActionRequestBody{Action: action},
	}).Return(nil, &serverError).Once()
}

func mockUpdateRapidRuleException(client *appsec.Mock, ruleID int64, exception appsec.RuleConditionException) {
	client.On("UpdateRapidRuleException", mock.Anything, appsec.UpdateRapidRuleExceptionRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
		RuleID:   ruleID,
		Body:     exception,
	}).Return((*appsec.UpdateRapidRuleExceptionResponse)(&exception), nil).Once()
}

func mockUpdateRapidRuleExceptionFailure(client *appsec.Mock, ruleID int64, exception appsec.RuleConditionException) {
	client.On("UpdateRapidRuleException", mock.Anything, appsec.UpdateRapidRuleExceptionRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
		RuleID:   ruleID,
		Body:     exception,
	}).Return(nil, &serverError).Once()
}

func toDefinitionsJSON(file string) string {
	data, _ := os.ReadFile("testdata/TestResRapidRules/" + file)
	return string(data)
}

var resourceReferenceName = "akamai_appsec_rapid_rules.test"

var serverError = appsec.Error{
	Type:       "internal_error",
	Title:      "Internal Server Error",
	Detail:     "Error creating zone",
	StatusCode: http.StatusInternalServerError,
}
