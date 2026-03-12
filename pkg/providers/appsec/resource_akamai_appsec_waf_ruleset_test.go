package appsec

import (
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v10/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWAFRulesetResource(t *testing.T) {
	t.Parallel()

	// Load fixture data
	wafRulesetResponse := appsec.CompositeRulesetResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFRuleset/WAFRuleset.json"), &wafRulesetResponse)
	require.NoError(t, err)

	// Load additional fixture responses for specific test cases
	rulesOnlyResponse := appsec.CompositeRulesetResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFRuleset/WAFRuleset_rules_only.json"), &rulesOnlyResponse)
	require.NoError(t, err)

	attackGroupsOnlyResponse := appsec.CompositeRulesetResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFRuleset/WAFRuleset_attack_groups_only.json"), &attackGroupsOnlyResponse)
	require.NoError(t, err)

	rulesAndGroupsResponse := appsec.CompositeRulesetResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFRuleset/WAFRuleset_rules_and_groups.json"), &rulesAndGroupsResponse)
	require.NoError(t, err)

	// Base state checker for resource tests
	baseChecker := test.NewStateChecker("akamai_appsec_waf_ruleset.test").
		CheckEqual("config_id", "111111").
		CheckEqual("security_policy_id", "2222_333333")

	// Base import checker for import tests
	baseImportChecker := test.NewImportChecker().
		CheckEqual("config_id", "111111").
		CheckEqual("security_policy_id", "2222_333333")

	var tests = map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"check schema - missing required attribute config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_missing_config_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},
		"check schema - missing required attribute security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile("Missing required argument"),
				},
			},
		},

		"check validation - invalid rule action": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_invalid_rule_action.tf"),
					ExpectError: regexp.MustCompile("invalid rule_action"),
				},
			},
		},
		"check validation - invalid attack group action": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_invalid_attack_group_action.tf"),
					ExpectError: regexp.MustCompile("invalid attack_group_action"),
				},
			},
		},
		"check validation - invalid condition exception JSON for rule": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_invalid_rule_condition_exception.tf"),
					ExpectError: regexp.MustCompile(`condition_exception must\s+be\s+valid`),
				},
			},
		},
		"check validation - invalid condition exception JSON for attack group": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_invalid_attack_group_condition_exception.tf"),
					ExpectError: regexp.MustCompile(`condition_exception must\s+be\s+valid`),
				},
			},
		},
		"check validation - duplicate rules": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_with_duplicate_rules.tf"),
					ExpectError: regexp.MustCompile("duplicate rule"),
				},
			},
		},
		"check validation - duplicate attack groups": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_with_duplicate_attack_groups.tf"),
					ExpectError: regexp.MustCompile("duplicate attack_group"),
				},
			},
		},

		"create waf ruleset - invalid rule ID not in ruleset": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_with_invalid_rule_id.tf"),
					ExpectError: regexp.MustCompile("Rule .* does not exist"),
				},
			},
		},
		"create waf ruleset - invalid attack group not in ruleset": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_with_invalid_attack_group.tf"),
					ExpectError: regexp.MustCompile("Attack group .* does not exist"),
				},
			},
		},

		"create waf ruleset - with rules": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 2)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, rulesOnlyResponse, 1)
				mockGetWAFCompositeRulesetForResource(m, rulesOnlyResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_with_rules.tf"),
					Check:  baseChecker.CheckEqual("rules.#", "3").Build(),
				},
			},
		},
		"create waf ruleset - with attack groups": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 2)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, attackGroupsOnlyResponse, 1)
				mockGetWAFCompositeRulesetForResource(m, attackGroupsOnlyResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_with_attack_groups.tf"),
					Check:  baseChecker.CheckEqual("attack_groups.#", "3").Build(),
				},
			},
		},
		"create waf ruleset - with rules and attack groups": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 2)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, rulesAndGroupsResponse, 1)
				mockGetWAFCompositeRulesetForResource(m, rulesAndGroupsResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_update_and_add_rule_and_group.tf"),
					Check: baseChecker.
						// Verify rules count and all fields inside the rule
						CheckEqual("rules.#", "1").
						CheckEqual("rules.0.rule_id", "950002").
						CheckEqual("rules.0.rule_action", "alert").
						CheckEqual("rules.0.condition_exception", `{"advancedExceptions":{"conditionOperator":"AND","conditions":[{"paths":["/catssssszzz*"],"positiveMatch":true,"type":"pathMatch"}]}}`).
						// Verify attack_groups count and all fields inside the attack group
						CheckEqual("attack_groups.#", "1").
						CheckEqual("attack_groups.0.attack_group", "CMD").
						CheckEqual("attack_groups.0.attack_group_action", "alert").
						CheckEqual("attack_groups.0.condition_exception", `{"advancedExceptions":{"conditionOperator":"AND","conditions":[{"paths":["/catssssszzz*"],"positiveMatch":true,"type":"pathMatch"}]}}`).
						Build(),
				},
			},
		},

		"create waf ruleset - minimal (config_id and policy_id only)": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 2)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_minimal.tf"),
					Check:  baseChecker.Build(),
				},
			},
		},
		"create waf ruleset - apply and apply again (idempotent)": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 4)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 4)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_minimal.tf"),
					Check:  baseChecker.Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_minimal.tf"),
					Check:  baseChecker.Build(),
				},
			},
		},

		"import WAF ruleset - verify rules and groups loaded": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 2)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 2)
			},
			steps: []resource.TestStep{
				{
					Config:           testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					ImportState:      true,
					ImportStateId:    "111111:2222_333333",
					ResourceName:     "akamai_appsec_waf_ruleset.test",
					ImportStateCheck: baseImportChecker.CheckEqual("rules.#", "309").CheckEqual("attack_groups.#", "10").Build(),
				},
			},
		},
		"import WAF ruleset - successful": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 2)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 2)
			},
			steps: []resource.TestStep{
				{
					Config:           testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					ImportState:      true,
					ImportStateId:    "111111:2222_333333",
					ResourceName:     "akamai_appsec_waf_ruleset.test",
					ImportStateCheck: baseImportChecker.Build(),
				},
			},
		},
		"create waf ruleset - Unable to read latest config version from API": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_minimal.tf"),
					ExpectError: regexp.MustCompile("Unable to read latest config version from API"),
				},
			},
		},
		"create waf ruleset - Unable to get WAF composite ruleset": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetFailureForResource(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_minimal.tf"),
					ExpectError: regexp.MustCompile("Unable to read WAF ruleset"),
				},
			},
		},
		"create waf ruleset - Unable to update WAF composite ruleset": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetFailureForResource(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_create_minimal.tf"),
					ExpectError: regexp.MustCompile("Unable to update WAF ruleset"),
				},
			},
		},

		"import WAF ruleset - Unable to read latest config version from API": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					ImportState:   true,
					ImportStateId: "111111:2222_333333",
					ResourceName:  "akamai_appsec_waf_ruleset.test",
					ExpectError:   regexp.MustCompile("invalid config version"),
				},
			},
		},
		"import WAF ruleset - Unable to read WAF composite ruleset": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetFailureForResource(m)
			},
			steps: []resource.TestStep{
				{
					Config:        testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					ImportState:   true,
					ImportStateId: "111111:2222_333333",
					ResourceName:  "akamai_appsec_waf_ruleset.test",
					ExpectError:   regexp.MustCompile("Unable to read WAF ruleset"),
				},
			},
		},
		"import WAF ruleset - invalid ID format": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					ImportState:        true,
					ImportStateId:      "12345",
					ResourceName:       "akamai_appsec_waf_ruleset.test",
					ExpectError:        regexp.MustCompile("ID '12345' incorrectly formatted: should be 'CONFIG_ID:SECURITY_POLICY_ID'"),
					ImportStatePersist: false,
				},
			},
		},
		"import WAF ruleset - invalid security policy ID format": {
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					ImportState:        true,
					ImportStateId:      "111111:",
					ResourceName:       "akamai_appsec_waf_ruleset.test",
					ExpectError:        regexp.MustCompile("invalid security policy id"),
					ImportStatePersist: false,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			client := &appsec.Mock{}
			if tc.init != nil {
				tc.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    tc.steps,
				})
			})

			client.AssertExpectations(t)
		})
	}
}

// mockGetWAFCompositeRulesetForResource mocks the GetWAFCompositeRuleset API call for resource tests
func mockGetWAFCompositeRulesetForResource(m *appsec.Mock, response appsec.CompositeRulesetResponse, times int) {
	m.On("GetWAFCompositeRuleset", mock.Anything, appsec.GetWAFCompositeRulesetRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(&response, nil).Times(times)
}

// mockGetWAFCompositeRulesetFailureForResource mocks a failed GetWAFCompositeRuleset API call
func mockGetWAFCompositeRulesetFailureForResource(m *appsec.Mock) {
	m.On("GetWAFCompositeRuleset", mock.Anything, appsec.GetWAFCompositeRulesetRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(nil, &wafRulesetResourceServerError).Once()
}

// mockUpdateWAFCompositeRulesetForResource mocks the UpdateWAFCompositeRuleset API call for resource tests
func mockUpdateWAFCompositeRulesetForResource(m *appsec.Mock, response appsec.CompositeRulesetResponse, times int) {
	m.On("UpdateWAFCompositeRuleset", mock.Anything, mock.MatchedBy(func(req appsec.UpdateWAFCompositeRulesetRequest) bool {
		return req.ConfigID == 111111 && req.PolicyID == "2222_333333"
	})).Return(&response, nil).Times(times)
}

// mockUpdateWAFCompositeRulesetFailureForResource mocks a failed UpdateWAFCompositeRuleset API call
func mockUpdateWAFCompositeRulesetFailureForResource(m *appsec.Mock) {
	m.On("UpdateWAFCompositeRuleset", mock.Anything, mock.MatchedBy(func(req appsec.UpdateWAFCompositeRulesetRequest) bool {
		return req.ConfigID == 111111 && req.PolicyID == "2222_333333"
	})).Return(nil, &wafRulesetResourceServerError).Once()
}

var wafRulesetResourceServerError = appsec.Error{
	Type:       "internal_error",
	Title:      "Internal Server Error",
	Detail:     "Error Fetching Composite Ruleset",
	StatusCode: http.StatusInternalServerError,
}
