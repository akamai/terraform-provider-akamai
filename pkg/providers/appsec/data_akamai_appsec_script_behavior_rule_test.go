package appsec

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v12/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v12/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestScriptBehaviorRuleDataSource(t *testing.T) {
	t.Parallel()

	for name, test := range scriptBehaviorRuleDSTests(t) {
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

func scriptBehaviorRuleDSTests(t *testing.T) map[string]struct {
	init  func(*appsec.Mock)
	steps []resource.TestStep
} {
	t.Helper()

	fixture := newScriptBehaviorRuleFixture()
	dsName := "data.akamai_appsec_script_behavior_rule.test"

	return map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"read single rule": {
			init: func(m *appsec.Mock) {
				m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 40576}).
					Return(&fixture.configResp, nil).Times(3)
				m.On("ListClientSideProtectionRules", mock.Anything, listRulesReq).
					Return(&fixture.listResp, nil).Times(3)
			},
			steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestDSScriptBehaviorRule/script_behavior_rule.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "rule_id", "3000601"),
					resource.TestCheckResourceAttr(dsName, "security_policy_id", "ab_86352"),
					resource.TestCheckResourceAttr(dsName, "rule_name", "Block PCI Data Access"),
					resource.TestCheckResourceAttr(dsName, "rule_description", "Prevents third-party scripts from accessing payment card data"),
					resource.TestCheckResourceAttr(dsName, "activity_type", "DATA_ACCESS"),
					resource.TestCheckResourceAttr(dsName, "action", "deny"),
					resource.TestCheckResourceAttr(dsName, "page_condition", "MATCHES"),
					resource.TestCheckResourceAttr(dsName, "page_values.#", "2"),
					resource.TestCheckResourceAttr(dsName, "page_values.0", "/checkout"),
					resource.TestCheckResourceAttr(dsName, "script_type", "THIRD_PARTY"),
					resource.TestCheckResourceAttr(dsName, "script_condition", "MATCHES"),
					resource.TestCheckResourceAttr(dsName, "script_values.#", "1"),
					resource.TestCheckResourceAttr(dsName, "activity_values_condition", "MATCHES"),
					resource.TestCheckResourceAttr(dsName, "activity_values.#", "2"),
				),
			}},
		},
		"rule not found": {
			init: func(m *appsec.Mock) {
				m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 40576}).
					Return(&fixture.configResp, nil)
				m.On("ListClientSideProtectionRules", mock.Anything, listRulesReq).
					Return(&appsec.ListClientSideProtectionRulesResponse{}, nil)
			},
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSScriptBehaviorRule/script_behavior_rule.tf"),
				ExpectError: regexp.MustCompile(`rule 3000601 not found in policy "ab_86352"`),
			}},
		},
		"error listing rules": {
			init: func(m *appsec.Mock) {
				m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 40576}).
					Return(&fixture.configResp, nil)
				m.On("ListClientSideProtectionRules", mock.Anything, listRulesReq).
					Return(nil, fmt.Errorf("listing rules failed"))
			},
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSScriptBehaviorRule/script_behavior_rule.tf"),
				ExpectError: regexp.MustCompile(`listing script behavior rules`),
			}},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSScriptBehaviorRule/script_behavior_rule_missing_config_id.tf"),
				ExpectError: regexp.MustCompile(`The argument "config_id" is required`),
			}},
		},
		"missing required argument security_policy_id": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSScriptBehaviorRule/script_behavior_rule_missing_security_policy_id.tf"),
				ExpectError: regexp.MustCompile(`The argument "security_policy_id" is required`),
			}},
		},
		"missing required argument rule_id": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSScriptBehaviorRule/script_behavior_rule_missing_rule_id.tf"),
				ExpectError: regexp.MustCompile(`The argument "rule_id" is required`),
			}},
		},
	}
}
