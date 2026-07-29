package appsec

import (
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v10/internal/edgegrid"
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
	importNoChangesResponse := appsec.CompositeRulesetResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFRuleset/WAFRuleset_import_no_changes.json"), &importNoChangesResponse)
	require.NoError(t, err)

	// Load delete request for import no changes test (resets rules/groups to "none" action)
	importNoChangesDeleteRequest := appsec.UpdateWAFCompositeRulesetRequest{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFRuleset/WAFRuleset_import_no_changes_delete_request.json"), &importNoChangesDeleteRequest)
	require.NoError(t, err)

	getSecurityPoliciesResponse := appsec.GetSecurityPoliciesResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestResWAFRuleset/SecurityPolicyWafRuleset.json"), &getSecurityPoliciesResponse)
	require.NoError(t, err)

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

		"update security_policy_id - should fail with error": {
			init: func(m *appsec.Mock) {
				// Step 1: create
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 1: read (post-create)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 1: read (plan before step 2)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: plan fails with error — framework still runs delete cleanup for step 1 resource
				mockGetConfiguration(m, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					// Create with original policy ID
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					Check:  baseChecker.Build(),
				},
				{
					// Attempt to change security_policy_id to a different known value —
					// PreventStringUpdateIfKnown must reject this during plan
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_updated_security_policy_id.tf"),
					ExpectError: regexp.MustCompile("updating 'security_policy_id' is not allowed"),
				},
			},
		},
		"update security_policy_id with variable reference - should fail": {
			init: func(m *appsec.Mock) {
				// Step 1: create with variable reference (policy_id is known value "2222_333333")
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 1: read (post-create)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 1: read (plan before step 2)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: plan fails with error — attempt to update to different policy_id
				// PreventStringUpdateIfKnown detects both old and new values are known
				mockGetConfiguration(m, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					// Create with original policy ID via variable reference
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_with_variable_reference.tf"),
					Check:  baseChecker.Build(),
				},
				{
					// Attempt to change security_policy_id via different variable reference —
					// PreventStringUpdateIfKnown must reject this during plan when both values are known
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_change_variable_reference.tf"),
					ExpectError: regexp.MustCompile("updating 'security_policy_id' is not allowed"),
				},
			},
		},
		"update security_policy_id to local reference with same value - no update triggered": {
			init: func(m *appsec.Mock) {
				// Step 1: create with inline literal
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 1: read (post-create)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: read (refresh)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: plan — PreventStringUpdateIfKnown sees equal values, no diff, no update
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Delete
				mockGetConfiguration(m, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					// Create with inline literal "2222_333333"
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					Check:  baseChecker.Build(),
				},
				{
					// Switch to local reference resolving to the same "2222_333333" —
					// PreventStringUpdateIfKnown must allow this (values are equal) and no API update should be called
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_security_policy_id_from_local.tf"),
					Check:  baseChecker.Build(),
				},
			},
		},
		"security_policy_id idempotent with local reference - UseStateForUnknown preserves value": {
			init: func(m *appsec.Mock) {
				// Step 1: create with local reference (value is "2222_333333", known at plan time)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 1: read (post-create)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: read (refresh) — UseStateForUnknown ensures state value is carried into plan
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: plan — no diff, no update
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Delete
				mockGetConfiguration(m, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_security_policy_id_from_local.tf"),
					Check:  baseChecker.Build(),
				},
				{
					// Re-apply identical config — UseStateForUnknown + PreventStringUpdateIfKnown must not trigger any update
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_security_policy_id_from_local.tf"),
					Check:  baseChecker.Build(),
				},
			},
		},
		"updating security policy in the waf ruleset resource which is imported already - should fail": {
			init: func(m *appsec.Mock) {
				// Step 1: ImportState call
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 1: Read call (framework reads after import to normalize state)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: refresh (reads imported state with policy "2222_333333")
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: data source plan — GetSecurityPolicies returns policy "AAAA_81230"
				mockGetConfiguration(m, 1)
				mockGetSecurityPolicies(m, getSecurityPoliciesResponse, 1)

				// Step 2: plan fails — PreventStringUpdateIfKnown detects "2222_333333" → "AAAA_81230"
				// framework runs delete cleanup for the imported resource
				mockGetConfiguration(m, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					ImportState:        true,
					ImportStateId:      "111111:2222_333333",
					ResourceName:       "akamai_appsec_waf_ruleset.test",
					ImportStatePersist: true,
					ImportStateCheck:   baseImportChecker.CheckEqual("rules.#", "308").CheckEqual("attack_groups.#", "9").Build(),
				},
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_reference_changed.tf"),
					ExpectError: regexp.MustCompile("updating 'security_policy_id' is not allowed"),
				},
			},
		},

		"same security policy in the waf ruleset resource which is imported already - should not fail": {
			init: func(m *appsec.Mock) {
				// Step 1: ImportState call
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 1: Read call (framework reads after import to normalize state)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: refresh (reads imported state with policy "2222_333333")
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: data source read (plan) — GetSecurityPolicies returns policy "2222_333333"
				mockGetConfiguration(m, 1)
				mockGetSecurityPolicies(m, getSecurityPoliciesResponse, 1)

				// Step 2: plan — PreventStringUpdateIfKnown sees equal values, no diff, no update
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Step 2: data source read (apply) — Terraform re-reads data sources during apply
				mockGetConfiguration(m, 1)
				mockGetSecurityPolicies(m, getSecurityPoliciesResponse, 1)

				// Step 2: data source read (post-apply idempotency check)
				mockGetConfiguration(m, 1)
				mockGetSecurityPolicies(m, getSecurityPoliciesResponse, 1)

				// Step 2: resource read (post-apply idempotency check)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// Delete
				mockGetConfiguration(m, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					ImportState:        true,
					ImportStateId:      "111111:2222_333333",
					ResourceName:       "akamai_appsec_waf_ruleset.test",
					ImportStatePersist: true,
					ImportStateCheck:   baseImportChecker.CheckEqual("rules.#", "308").CheckEqual("attack_groups.#", "9").Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_reference_unchanged.tf"),
					Check:  baseChecker.Build(),
				},
			},
		},

		"create waf ruleset - invalid rule ID not in ruleset": {
			init: func(m *appsec.Mock) {
				// create (fails validation, no update/delete needed)
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
				// create (fails validation, no update/delete needed)
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
				// create
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, rulesOnlyResponse, 1)

				// read
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, rulesOnlyResponse, 1)

				// delete
				mockGetConfiguration(m, 1)
				mockUpdateWAFCompositeRulesetForResource(m, rulesOnlyResponse, 1)
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
				// create
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, attackGroupsOnlyResponse, 1)

				// read
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, attackGroupsOnlyResponse, 1)

				// delete
				mockGetConfiguration(m, 1)
				mockUpdateWAFCompositeRulesetForResource(m, attackGroupsOnlyResponse, 1)
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
				// create
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, rulesAndGroupsResponse, 1)

				// read
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, rulesAndGroupsResponse, 1)

				// delete
				mockGetConfiguration(m, 1)
				mockUpdateWAFCompositeRulesetForResource(m, rulesAndGroupsResponse, 1)
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
				// create
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// read
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// delete
				mockGetConfiguration(m, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
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
				// first apply - create
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
				mockUpdateWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// first apply - read
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// second apply - read (refresh)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// second apply - read (plan, no update since idempotent)
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// delete
				mockGetConfiguration(m, 1)
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
				// import
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// read after import
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config:           testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset.tf"),
					ImportState:      true,
					ImportStateId:    "111111:2222_333333",
					ResourceName:     "akamai_appsec_waf_ruleset.test",
					ImportStateCheck: baseImportChecker.CheckEqual("rules.#", "308").CheckEqual("attack_groups.#", "9").Build(),
				},
			},
		},
		"import WAF ruleset - successful": {
			init: func(m *appsec.Mock) {
				// import
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)

				// read after import
				mockGetConfiguration(m, 1)
				mockGetWAFCompositeRulesetForResource(m, wafRulesetResponse, 1)
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
				// create (fails to get config version)
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
				// create (fails to get ruleset)
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
				// create (fails to update ruleset)
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
				// import (fails to get config version)
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
				// import (fails to get ruleset)
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
		"import WAF ruleset - import then apply no changes": {
			init: func(m *appsec.Mock) {
				// Step 1: Import
				mockGetConfigurationForImportNoChanges(m, 1)
				mockGetWAFCompositeRulesetForImportNoChanges(m, importNoChangesResponse, 1)

				// Step 2: Apply (refresh + plan, no update needed since no changes)
				mockGetConfigurationForImportNoChanges(m, 1)
				mockGetWAFCompositeRulesetForImportNoChanges(m, importNoChangesResponse, 1)
				mockGetConfigurationForImportNoChanges(m, 1)
				mockGetWAFCompositeRulesetForImportNoChanges(m, importNoChangesResponse, 1)
				mockGetConfigurationForImportNoChanges(m, 1)
				mockGetWAFCompositeRulesetForImportNoChanges(m, importNoChangesResponse, 1)

				// Step 3: Delete (cleanup at end of test)
				mockGetConfigurationForImportNoChanges(m, 1)
				mockUpdateWAFCompositeRulesetForImportNoChanges(m, importNoChangesDeleteRequest, importNoChangesResponse, 1)
			},
			steps: []resource.TestStep{
				{
					Config:             testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_import_no_changes.tf"),
					ImportState:        true,
					ImportStateId:      "63905:p563_113699",
					ResourceName:       "akamai_appsec_waf_ruleset.import_no_changes",
					ImportStatePersist: true,
					ImportStateCheck: test.NewImportChecker().
						CheckEqual("config_id", "63905").
						CheckEqual("security_policy_id", "p563_113699").
						CheckEqual("rules.#", "1").
						CheckEqual("attack_groups.#", "2").
						Build(),
				},
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestResWAFRuleset/waf_ruleset_import_no_changes.tf"),
					Check: test.NewStateChecker("akamai_appsec_waf_ruleset.import_no_changes").
						CheckEqual("config_id", "63905").
						CheckEqual("security_policy_id", "p563_113699").
						CheckEqual("rules.#", "1").
						CheckEqual("attack_groups.#", "2").
						Build(),
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

func mockGetSecurityPolicies(m *appsec.Mock, response appsec.GetSecurityPoliciesResponse, times int) {
	m.On("GetSecurityPolicies",
		mock.Anything,
		appsec.GetSecurityPoliciesRequest{ConfigID: 111111, Version: 2},
	).Return(&response, nil).Times(times)
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

// Mock functions for import no changes test (config_id: 63905, policy_id: p563_113699)

// mockGetConfigurationForImportNoChanges mocks the GetConfiguration API call for import no changes test
func mockGetConfigurationForImportNoChanges(m *appsec.Mock, times int) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{
		ConfigID: 63905,
	}).Return(&appsec.GetConfigurationResponse{
		ID:            63905,
		LatestVersion: 2,
	}, nil).Times(times)
}

// mockGetWAFCompositeRulesetForImportNoChanges mocks the GetWAFCompositeRuleset API call for import no changes test
func mockGetWAFCompositeRulesetForImportNoChanges(m *appsec.Mock, response appsec.CompositeRulesetResponse, times int) {
	m.On("GetWAFCompositeRuleset", mock.Anything, appsec.GetWAFCompositeRulesetRequest{
		ConfigID: 63905,
		Version:  2,
		PolicyID: "p563_113699",
	}).Return(&response, nil).Times(times)
}

// mockUpdateWAFCompositeRulesetForImportNoChanges mocks the UpdateWAFCompositeRuleset API call for import no changes test
func mockUpdateWAFCompositeRulesetForImportNoChanges(m *appsec.Mock, request appsec.UpdateWAFCompositeRulesetRequest, response appsec.CompositeRulesetResponse, times int) {
	m.On("UpdateWAFCompositeRuleset", mock.Anything, request).Return(&response, nil).Times(times)
}
