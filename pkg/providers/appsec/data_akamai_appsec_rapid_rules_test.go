package appsec

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/ptr"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v8/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDataRapidRules(t *testing.T) {

	getRapidRulesResponse := appsec.GetRapidRulesResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSRapidRules/RapidRules.json"), &getRapidRulesResponse)
	require.NoError(t, err)

	getRapidRulesResponseForSingleRule := appsec.GetRapidRulesResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSRapidRules/RapidRule.json"), &getRapidRulesResponseForSingleRule)
	require.NoError(t, err)

	getRapidRulesDefaultActionResponse := appsec.GetRapidRulesDefaultActionResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSRapidRules/RapidRulesDefaultAction.json"), &getRapidRulesDefaultActionResponse)
	require.NoError(t, err)

	getAttackGroupsResponse := appsec.GetAttackGroupsResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSRapidRules/AttackGroups.json"), &getAttackGroupsResponse)
	require.NoError(t, err)

	getRapidRulesWithExpiry := appsec.GetRapidRulesResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSRapidRules/RapidRulesWithExpiry.json"), &getRapidRulesWithExpiry)
	require.NoError(t, err)

	baseChecker := test.NewStateChecker("data.akamai_appsec_rapid_rules.test").
		CheckEqual("enabled", "true")

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - rapid rules disabled status should be returned": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockGetRapidRulesStatus(m, false, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSRapidRules/rapid_rules_no_rule_id.tf"),
					Check: baseChecker.
						CheckEqual("enabled", "false").
						CheckEqual("default_action", "No default action. Rapid rules is turned off.").
						CheckEqual("output_text", "Rapid rules is turned off.").
						Build(),
				},
			},
		},
		"happy path - return rapid rules": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockGetRapidRulesStatus(m, true, 3)
				mockGetRapidRules(m, getRapidRulesResponse, 3)
				mockGetRapidRulesDefaultAction(m, "alert", 3)
				mockGetAttackGroups(m, getAttackGroupsResponse, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSRapidRules/rapid_rules_no_rule_id.tf"),
					Check: baseChecker.
						CheckEqual("rapid_rules.0.action", "alert").
						CheckEqual("rapid_rules.0.attack_group", "LFI").
						CheckEqual("rapid_rules.0.attack_group_exception", indentJSON(t, "testdata/TestDSRapidRules/AttackGroupException.json")).
						CheckEqual("rapid_rules.0.condition_exception", "{}").
						CheckEqual("rapid_rules.0.id", "300050014").
						CheckEqual("rapid_rules.0.lock", "false").
						CheckEqual("rapid_rules.0.name", "Local File Inclusion (LFI) Attack on Linux files").
						CheckMissing("rapid_rules.0.expired").
						CheckMissing("rapid_rules.0.expire_in_days").
						CheckEqual("rapid_rules.1.action", "alert").
						CheckEqual("rapid_rules.1.attack_group", "PLATFORM").
						CheckEqual("rapid_rules.1.attack_group_exception", "{}").
						CheckEqual("rapid_rules.1.condition_exception", "{}").
						CheckEqual("rapid_rules.1.id", "777888999").
						CheckEqual("rapid_rules.1.lock", "false").
						CheckEqual("rapid_rules.1.name", "Citrix Virtual Apps and Desktops (XEN) Unauthenticated RCE (CVE-2024-8068 and CVE-2024-8069) Attack Detected").
						CheckEqual("rapid_rules.2.action", "alert").
						CheckEqual("rapid_rules.2.attack_group", "XSS").
						CheckEqual("rapid_rules.2.attack_group_exception", "{}").
						CheckEqual("rapid_rules.2.condition_exception", indentJSON(t, "testdata/TestDSRapidRules/ConditionException.json")).
						CheckEqual("rapid_rules.2.id", "999997").
						CheckEqual("rapid_rules.2.lock", "false").
						CheckEqual("rapid_rules.2.name", "Cross-site Scripting (XSS) Attack - Dummy test rule").
						Build(),
				},
			},
		},
		"happy path - return single rapid rule": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockGetRapidRulesStatus(m, true, 3)
				mockGetRapidRule(m, getRapidRulesResponseForSingleRule, 3)
				mockGetRapidRulesDefaultAction(m, "alert", 3)
				mockGetAttackGroups(m, getAttackGroupsResponse, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSRapidRules/rapid_rules_with_rule_id.tf"),
					Check: baseChecker.
						CheckEqual("rapid_rules.0.action", "alert").
						CheckEqual("rapid_rules.0.attack_group", "XSS").
						CheckEqual("rapid_rules.0.attack_group_exception", "{}").
						CheckEqual("rapid_rules.0.condition_exception", indentJSON(t, "testdata/TestDSRapidRules/ConditionException.json")).
						CheckEqual("rapid_rules.0.id", "999997").
						CheckEqual("rapid_rules.0.lock", "false").
						CheckEqual("rapid_rules.0.name", "Cross-site Scripting (XSS) Attack - Dummy test rule").
						CheckMissing("rapid_rules.0.expired").
						CheckMissing("rapid_rules.0.expire_in_days").
						Build(),
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSRapidRules/rapid_rules_missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "config_id" is required, but no definition was found`),
				},
			},
		},
		"missing required argument security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSRapidRules/rapid_rules_missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "security_policy_id" is required, but no definition was found`),
				},
			},
		},
		"error response from GetConfiguration api": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSRapidRules/rapid_rules_no_rule_id.tf"),
					ExpectError: regexp.MustCompile("Error: invalid config version"),
				},
			},
		},
		"error response from GetRapidRulesStatus api": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockGetRapidRulesStatusFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSRapidRules/rapid_rules_no_rule_id.tf"),
					ExpectError: regexp.MustCompile("Error: calling 'GetRapidRulesStatus'"),
				},
			},
		},
		"error response from GetRapidRules api": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockGetRapidRulesStatus(m, true, 1)
				mockGetRapidRulesFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSRapidRules/rapid_rules_no_rule_id.tf"),
					ExpectError: regexp.MustCompile("Error: calling 'GetRapidRules'"),
				},
			},
		},
		"error response from GetRapidRulesDefaultAction api": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockGetRapidRulesStatus(m, true, 1)
				mockGetRapidRules(m, getRapidRulesResponse, 1)
				mockGetRapidRulesDefaultActionFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSRapidRules/rapid_rules_no_rule_id.tf"),
					ExpectError: regexp.MustCompile("Error: calling 'getRapidRulesDefaultAction'"),
				},
			},
		},
		"error response from GetAttackGroups api": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockGetRapidRulesStatus(m, true, 1)
				mockGetRapidRules(m, getRapidRulesResponse, 1)
				mockGetRapidRulesDefaultAction(m, "alert", 1)
				mockGetAttackGroupsFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSRapidRules/rapid_rules_no_rule_id.tf"),
					ExpectError: regexp.MustCompile("Error: calling 'GetAttackGroups'"),
				},
			},
		},
		"happy path – include_expiry_details = true": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockGetRapidRulesStatus(m, true, 3)
				mockGetRapidRulesWithExpiry(m, getRapidRulesWithExpiry, 3) // flag = true
				mockGetRapidRulesDefaultAction(m, "alert", 3)
				mockGetAttackGroups(m, getAttackGroupsResponse, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t,
						"testdata/TestDSRapidRules/rapid_rules_with_expiry.tf"),
					Check: baseChecker.
						CheckEqual("include_expiry_details", "true").
						CheckEqual("rapid_rules.0.expired", "true").
						CheckMissing("rapid_rules.0.expire_in_days").
						CheckEqual("rapid_rules.1.expire_in_days", "5").
						Build(),
				},
			},
		},
		"include_expiry_details = true but rapid rules disabled": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockGetRapidRulesStatus(m, false, 3) // disabled ⇒ provider never calls GetRapidRules
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t,
						"testdata/TestDSRapidRules/rapid_rules_with_expiry.tf"),
					Check: baseChecker.
						CheckEqual("enabled", "false").
						CheckEqual("default_action", "No default action. Rapid rules is turned off.").
						CheckEqual("output_text", "Rapid rules is turned off.").
						CheckEqual("include_expiry_details", "true").
						CheckMissing("rapid_rules.0.expired").
						CheckMissing("rapid_rules.0.expire_in_days").
						Build(),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &appsec.Mock{}
			if test.init != nil {
				test.init(client)
			}

			useClient(client, func() {
				resource.UnitTest(t, resource.TestCase{
					ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
					IsUnitTest:               true,
					Steps:                    test.steps,
				})
			})

			client.AssertExpectations(t)
		})
	}
}

func mockGetRapidRule(m *appsec.Mock, response appsec.GetRapidRulesResponse, times int) {
	m.On("GetRapidRules", mock.Anything, appsec.GetRapidRulesRequest{ConfigID: 111111, Version: 2, PolicyID: "2222_333333", RuleID: ptr.To(int64(999997))}).
		Return(&response, nil).Times(times)
}

func mockGetRapidRulesDefaultActionFailure(m *appsec.Mock) {
	m.On("GetRapidRulesDefaultAction", mock.Anything, appsec.GetRapidRulesDefaultActionRequest{ConfigID: 111111, Version: 2, PolicyID: "2222_333333"}).
		Return(nil, &serverError).Once()
}

func mockGetAttackGroups(m *appsec.Mock, response appsec.GetAttackGroupsResponse, times int) {
	m.On("GetAttackGroups", mock.Anything, appsec.GetAttackGroupsRequest{ConfigID: 111111, Version: 2, PolicyID: "2222_333333"}).
		Return(&response, nil).Times(times)
}

func mockGetAttackGroupsFailure(m *appsec.Mock) {
	m.On("GetAttackGroups", mock.Anything, appsec.GetAttackGroupsRequest{ConfigID: 111111, Version: 2, PolicyID: "2222_333333"}).
		Return(nil, &serverError).Once()
}

func mockGetRapidRulesWithExpiry(m *appsec.Mock, resp appsec.GetRapidRulesResponse, times int) {
	m.On("GetRapidRules", mock.Anything, appsec.GetRapidRulesRequest{ConfigID: 111111, Version: 2, PolicyID: "2222_333333", IncludeExpiryDetails: true, RuleID: nil}).
		Return(&resp, nil).Times(times)
}

// indentJSON converts json file to a JSON-encoded string with Indent
func indentJSON(t *testing.T, path string) string {
	ruleException := appsec.RuleConditionException{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, path), &ruleException)
	require.NoError(t, err)
	jsonBytes, err := json.MarshalIndent(ruleException, "", "  ")
	require.NoError(t, err)
	return string(jsonBytes)
}
