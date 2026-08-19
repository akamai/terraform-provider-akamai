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
	"github.com/stretchr/testify/require"
)

// notFoundError simulates the 404 INVALID-INPUT-ERROR response returned by the AI rules
// API for policies whose enrollment state doesn't support a given endpoint.
var notFoundError = appsec.Error{
	Type:       "INVALID-INPUT-ERROR",
	Title:      "Invalid Input Error",
	Detail:     "HTTP 404 Not Found",
	StatusCode: http.StatusNotFound,
}

func TestDataWAFAIRules(t *testing.T) {
	t.Parallel()

	getAIRulesResp := appsec.ListAIRulesResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSWafAIRules/AIRules.json"), &getAIRulesResp)
	require.NoError(t, err)

	// NOT_ENROLLED policies return 200 from GET /ai-rules with status + rules populated.
	getAIRulesNotEnrolledResp := appsec.ListAIRulesResponse{}
	err = json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSWafAIRules/AIRulesNotEnrolled.json"), &getAIRulesNotEnrolledResp)
	require.NoError(t, err)

	baseChecker := test.NewStateChecker("data.akamai_appsec_waf_ai_rules.test").
		CheckEqual("config_id", "111111").
		CheckEqual("security_policy_id", "2222_333333").
		CheckEqual("ai_rule_status", "ENABLED")

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - return all AI rules": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockListWAFAIRules(m, getAIRulesResp, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSWafAIRules/waf_ai_rules.tf"),
					Check: baseChecker.
						CheckEqual("ai_rules.#", "2").
						CheckEqual("ai_rules.0.rule_id", "3001000").
						CheckEqual("ai_rules.0.rule_version", "1").
						CheckEqual("ai_rules.0.title", "AI-Detected SQL Injection Attack").
						CheckEqual("ai_rules.0.risk_score_group", "SQL_AI").
						CheckEqual("ai_rules.0.rule_description", "A SQL Injection attack consists of insertion or injection of a SQL query via the input data from the client to the application.").
						CheckEqual("ai_rules.0.action", "alert").
						CheckEqual("ai_rules.0.condition_exception", "[\n  {\n    \"exception\": {\n      \"specificHeaderCookieParamXmlOrJsonNames\": [\n        {\n          \"names\": [\n            \"test\"\n          ],\n          \"selector\": \"REQUEST_COOKIES_NAMES\",\n          \"wildcard\": true\n        }\n      ]\n    }\n  }\n]").
						CheckEqual("ai_rules.1.rule_id", "3001001").
						CheckEqual("ai_rules.1.risk_score_group", "XSS_AI").
						CheckEqual("ai_rules.1.rule_description", "Cross-Site Scripting attacks are a type of injection problem, in which malicious scripts are injected into otherwise benign web sites.").
						CheckEqual("ai_rules.1.action", "deny").
						CheckEqual("ai_rules.1.condition_exception", "[\n  {\n    \"exception\": {\n      \"specificHeaderCookieParamXmlOrJsonNames\": [\n        {\n          \"names\": [\n            \"test\"\n          ],\n          \"selector\": \"REQUEST_COOKIES_NAMES\",\n          \"wildcard\": true\n        }\n      ]\n    }\n  }\n]").
						CheckEqual("output_text", "\n+----------------------------------------------------------------------------------+\n| AIRulesDS                                                                        |\n+---------+---------+----------------------------------+------------------+--------+\n| RULE ID | VERSION | TITLE                            | RISK SCORE GROUP | ACTION |\n+---------+---------+----------------------------------+------------------+--------+\n| 3001000 | 1       | AI-Detected SQL Injection Attack | SQL_AI           | alert  |\n| 3001001 | 1       | AI-Detected XSS Attack           | XSS_AI           | deny   |\n+---------+---------+----------------------------------+------------------+--------+\n").
						Build(),
				},
			},
		},
		"not enrolled - GET /ai-rules returns 200 with NOT_ENROLLED and rules": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockListWAFAIRules(m, getAIRulesNotEnrolledResp, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSWafAIRules/waf_ai_rules.tf"),
					Check: test.NewStateChecker("data.akamai_appsec_waf_ai_rules.test").
						CheckEqual("config_id", "111111").
						CheckEqual("security_policy_id", "2222_333333").
						CheckEqual("ai_rule_status", "NOT_ENROLLED").
						CheckEqual("ai_rules.#", "0").
						Build(),
				},
			},
		},
		"enrolled policy - GET /ai-rules 404, falls back to GET /ai-rules/status": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockListWAFAIRules404(m, 3)
				mockGetAIRulesStatus(m, appsec.GetAIRulesStatusResponse{AIRuleStatus: "ENABLED"}, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSWafAIRules/waf_ai_rules.tf"),
					Check: test.NewStateChecker("data.akamai_appsec_waf_ai_rules.test").
						CheckEqual("config_id", "111111").
						CheckEqual("security_policy_id", "2222_333333").
						CheckEqual("ai_rule_status", "ENABLED").
						CheckEqual("ai_rules.#", "0").
						Build(),
				},
			},
		},
		"both endpoints 404 - returns NOT_ENROLLED with empty rules": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 3)
				mockListWAFAIRules404(m, 3)
				mockGetAIRulesStatus404(m, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSWafAIRules/waf_ai_rules.tf"),
					Check: test.NewStateChecker("data.akamai_appsec_waf_ai_rules.test").
						CheckEqual("config_id", "111111").
						CheckEqual("security_policy_id", "2222_333333").
						CheckEqual("ai_rule_status", "NOT_ENROLLED").
						CheckEqual("ai_rules.#", "0").
						Build(),
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSWafAIRules/waf_ai_rules_missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "config_id" is required, but no definition was found`),
				},
			},
		},
		"missing required argument security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSWafAIRules/waf_ai_rules_missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "security_policy_id" is required, but no definition was found`),
				},
			},
		},
		"error from GetConfiguration api": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationFailure(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSWafAIRules/waf_ai_rules.tf"),
					ExpectError: regexp.MustCompile("Error: fetching latest config version"),
				},
			},
		},
		"error from GetAIRules api (non-404)": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockListWAFAIRulesFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSWafAIRules/waf_ai_rules.tf"),
					ExpectError: regexp.MustCompile("Error: calling 'ListAIRules'"),
				},
			},
		},
		"error from GetAIRulesStatus api (after GetAIRules 404)": {
			init: func(m *appsec.Mock) {
				mockGetConfiguration(m, 1)
				mockListWAFAIRules404(m, 1)
				mockListWAFAIRulesStatusFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSWafAIRules/waf_ai_rules.tf"),
					ExpectError: regexp.MustCompile("Error: calling 'GetAIRulesStatus'"),
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

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewTestProtoV6ProviderFactory(client, NewSubprovider()),
				Steps:                    tc.steps,
			})

			client.APPSEC.AssertExpectations(t)
		})
	}
}

func mockListWAFAIRules(client *appsec.Mock, resp appsec.ListAIRulesResponse, times int) {
	client.On("ListAIRules", testutils.MockContext, appsec.ListAIRulesRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(&resp, nil).Times(times)
}

func mockListWAFAIRules404(client *appsec.Mock, times int) {
	client.On("ListAIRules", testutils.MockContext, appsec.ListAIRulesRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(nil, &notFoundError).Times(times)
}

func mockListWAFAIRulesFailure(client *appsec.Mock) {
	client.On("ListAIRules", testutils.MockContext, appsec.ListAIRulesRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(nil, &serverError).Once()
}

func mockGetAIRulesStatus404(client *appsec.Mock, times int) {
	client.On("GetAIRulesStatus", testutils.MockContext, appsec.GetAIRulesStatusRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(nil, &notFoundError).Times(times)
}

func mockListWAFAIRulesStatusFailure(client *appsec.Mock) {
	client.On("GetAIRulesStatus", testutils.MockContext, appsec.GetAIRulesStatusRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(nil, &serverError).Once()
}
