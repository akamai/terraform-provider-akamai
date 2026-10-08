package appsec

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v14/pkg/appsec"
	"github.com/akamai/terraform-provider-akamai/v11/internal/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v11/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestScriptBehaviorRulesDataSource(t *testing.T) {
	t.Parallel()

	for name, test := range scriptBehaviorRulesDSTests(t) {
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

func scriptBehaviorRulesDSTests(t *testing.T) map[string]struct {
	init  func(*appsec.Mock)
	steps []resource.TestStep
} {
	t.Helper()

	fixture := newScriptBehaviorRuleFixture()
	dsName := "data.akamai_appsec_script_behavior_rules.test"

	return map[string]struct {
		init  func(*appsec.Mock)
		steps []resource.TestStep
	}{
		"read all rules": {
			init: func(m *appsec.Mock) {
				m.On("GetConfiguration", mock.Anything, appsec.GetConfigurationRequest{ConfigID: 40576}).
					Return(&fixture.configResp, nil).Times(3)
				m.On("ListClientSideProtectionRules", mock.Anything, listRulesReq).
					Return(&fixture.listResp, nil).Times(3)
			},
			steps: []resource.TestStep{{
				Config: testutils.LoadFixtureString(t, "testdata/TestDSScriptBehaviorRules/script_behavior_rules.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "security_policy_id", "ab_86352"),
					resource.TestCheckResourceAttr(dsName, "rules.#", "1"),
					resource.TestCheckResourceAttr(dsName, "rules.0.rule_id", "3000601"),
					resource.TestCheckResourceAttr(dsName, "rules.0.rule_name", "Block PCI Data Access"),
					resource.TestCheckResourceAttr(dsName, "rules.0.rule_description", "Prevents third-party scripts from accessing payment card data"),
					resource.TestCheckResourceAttr(dsName, "rules.0.activity_type", "DATA_ACCESS"),
					resource.TestCheckResourceAttr(dsName, "rules.0.action", "deny"),
					resource.TestCheckResourceAttr(dsName, "rules.0.page_condition", "MATCHES"),
					resource.TestCheckResourceAttr(dsName, "rules.0.page_values.#", "2"),
					resource.TestCheckResourceAttr(dsName, "rules.0.page_values.0", "/checkout"),
					resource.TestCheckResourceAttr(dsName, "rules.0.page_values.1", "/billing/*"),
					resource.TestCheckResourceAttr(dsName, "rules.0.script_type", "THIRD_PARTY"),
					resource.TestCheckResourceAttr(dsName, "rules.0.script_condition", "MATCHES"),
					resource.TestCheckResourceAttr(dsName, "rules.0.script_values.#", "1"),
					resource.TestCheckResourceAttr(dsName, "rules.0.script_values.0", "analytics.example.com/*"),
					resource.TestCheckResourceAttr(dsName, "rules.0.activity_values_condition", "MATCHES"),
					resource.TestCheckResourceAttr(dsName, "rules.0.activity_values.#", "2"),
					resource.TestCheckResourceAttr(dsName, "rules.0.activity_values.0", "CREDIT_CARD"),
					resource.TestCheckResourceAttr(dsName, "rules.0.activity_values.1", "PII"),
				),
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
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSScriptBehaviorRules/script_behavior_rules.tf"),
				ExpectError: regexp.MustCompile(`listing script behavior rules`),
			}},
		},
		"missing required argument config_id": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSScriptBehaviorRules/script_behavior_rules_missing_config_id.tf"),
				ExpectError: regexp.MustCompile(`The argument "config_id" is required`),
			}},
		},
		"missing required argument security_policy_id": {
			steps: []resource.TestStep{{
				Config:      testutils.LoadFixtureString(t, "testdata/TestDSScriptBehaviorRules/script_behavior_rules_missing_security_policy_id.tf"),
				ExpectError: regexp.MustCompile(`The argument "security_policy_id" is required`),
			}},
		},
	}
}
