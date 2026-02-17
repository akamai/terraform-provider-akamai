package appsec

import (
	"encoding/json"
	"net/http"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/test"
	"github.com/akamai/terraform-provider-akamai/v9/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDataWAFRuleset(t *testing.T) {

	getWAFRulesetResponse := appsec.CompositeRulesetResponse{}
	err := json.Unmarshal(testutils.LoadFixtureBytes(t, "testdata/TestDSWAFRuleset/WAFRuleset.json"), &getWAFRulesetResponse)
	require.NoError(t, err)

	baseChecker := test.NewStateChecker("data.akamai_appsec_waf_ruleset.test")

	tests := map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"happy path - return WAF ruleset with attack groups and rules": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationForWAFRuleset(m, 3)
				mockGetWAFCompositeRuleset(m, getWAFRulesetResponse, 3)
			},
			steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSWAFRuleset/waf_ruleset.tf"),
					Check: baseChecker.
						CheckEqual("config_id", "111111").
						CheckEqual("security_policy_id", "2222_333333").

						// Attack groups are sorted alphabetically by group name
						// First attack group: CMD (has condition exception)
						CheckEqual("attack_groups.0.attack_group", "CMD").
						CheckEqual("attack_groups.0.attack_group_action", "alert").
						CheckEqual("attack_groups.0.condition_exception", indentJSON(t, "testdata/TestDSWAFRuleset/ConditionException.json")).

						// Check OUTBOUND attack group with "none" action (3rd after sorting)
						CheckEqual("attack_groups.2.attack_group", "OUTBOUND").
						CheckEqual("attack_groups.2.attack_group_action", "none").
						CheckEqual("attack_groups.2.condition_exception", "{}").

						// POLICY attack group without condition exception (5th after sorting)
						CheckEqual("attack_groups.4.attack_group", "POLICY").
						CheckEqual("attack_groups.4.attack_group_action", "alert").
						CheckEqual("attack_groups.4.condition_exception", "{}").

						// Rules are sorted by rule ID
						// First rule (ID 950002)
						CheckEqual("rules.0.rule_id", "950002").
						CheckEqual("rules.0.rule_action", "deny").
						CheckEqual("rules.0.condition_exception", "{}").

						// Rule ID 950006 with none action (2nd after sorting)
						CheckEqual("rules.1.rule_id", "950006").
						CheckEqual("rules.1.rule_action", "none").
						CheckEqual("rules.1.condition_exception", "{}").

						// Rule ID 950007 with condition exception & alert (3rd after sorting)
						CheckEqual("rules.2.rule_id", "950007").
						CheckEqual("rules.2.rule_action", "alert").
						CheckEqual("rules.2.condition_exception", indentJSON(t, "testdata/TestDSWAFRuleset/ConditionException.json")).

						// Last rule (ID 99999900)
						CheckEqual("rules.308.rule_id", "99999900").
						CheckEqual("rules.308.rule_action", "deny").
						CheckEqual("rules.308.condition_exception", "{}").
						Build(),
				},
			},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSWAFRuleset/waf_ruleset_missing_config_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "config_id" is required, but no definition was found`),
				},
			},
		},
		"missing required argument security_policy_id": {
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSWAFRuleset/waf_ruleset_missing_security_policy_id.tf"),
					ExpectError: regexp.MustCompile(`The argument "security_policy_id" is required, but no definition was found`),
				},
			},
		},
		"error response from GetConfiguration api": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationFailureForWAFRuleset(m, 1)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSWAFRuleset/waf_ruleset.tf"),
					ExpectError: regexp.MustCompile("retrieving config version"),
				},
			},
		},
		"error response from GetWAFCompositeRuleset api": {
			init: func(m *appsec.Mock) {
				mockGetConfigurationForWAFRuleset(m, 1)
				mockGetWAFCompositeRulesetFailure(m)
			},
			steps: []resource.TestStep{
				{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSWAFRuleset/waf_ruleset.tf"),
					ExpectError: regexp.MustCompile("Error: calling 'GetWAFCompositeRuleset'"),
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

func mockGetConfigurationForWAFRuleset(m *appsec.Mock, times int) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 111111}).
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

func mockGetConfigurationFailureForWAFRuleset(m *appsec.Mock, times int) {
	m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 111111}).
		Return(nil, &wafRulesetServerError).Times(times)
}

func mockGetWAFCompositeRuleset(m *appsec.Mock, response appsec.CompositeRulesetResponse, times int) {
	m.On("GetWAFCompositeRuleset", mock.Anything, appsec.GetWAFCompositeRulesetRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(&response, nil).Times(times)
}

func mockGetWAFCompositeRulesetFailure(m *appsec.Mock) {
	m.On("GetWAFCompositeRuleset", mock.Anything, appsec.GetWAFCompositeRulesetRequest{
		ConfigID: 111111,
		Version:  2,
		PolicyID: "2222_333333",
	}).Return(nil, &wafRulesetServerError).Once()
}

var wafRulesetServerError = appsec.Error{
	Type:       "internal_error",
	Title:      "Internal Server Error",
	Detail:     "Error Fetching Composite Ruleset",
	StatusCode: http.StatusInternalServerError,
}
