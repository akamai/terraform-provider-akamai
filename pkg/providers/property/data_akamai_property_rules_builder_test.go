package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestDataPropertyRulesBuilder(t *testing.T) {
	t.Run("rule empty options - v2024-01-09", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/rules_v2024_01_09_with_empty_options.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2024-01-09"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/default_v2024_01_09_with_empty_options.json")),
					),
				}},
			})
		})
	})
	t.Run("invalid rule with 3 children with different versions", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/rules_mixed_versions.tf"),
					ExpectError: regexp.MustCompile(`child rule is using different rule format \(rules_v2023_05_30\) than expected \(rules_v2023_01_05\)`),
				}},
			})
		})
	})
	t.Run("fails on rule with more than one behavior in one block", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/rules_error_too_many_elements.tf"),
					ExpectError: regexp.MustCompile(`expected 1 element\(s\), got 2`),
				}},
			})
		})
	})
	t.Run("fails on rule with is_secure outside default rule", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/rules_with_is_secure_outside_default.tf"),
					ExpectError: regexp.MustCompile(`cannot be used outside 'default' rule: is_secure`),
				}},
			})
		})
	})
	t.Run("fails on rule with variable outside default rule", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/rules_with_variable_outside_default.tf"),
					ExpectError: regexp.MustCompile(`cannot be used outside 'default' rule: variable`),
				}},
			})
		})
	})
	t.Run("valid rule with one child and some values are variables", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/rules_variables.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/default_variables.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/content_compression_variables.json")),
					),
				}},
			})
		})
	})
}

func testCheckResourceAttrJSON(name, key, value string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		attrs := s.RootModule().Resources[name].Primary.Attributes
		diff := cmp.Diff(value, attrs[key])
		if diff != "" {
			return fmt.Errorf("got from: %s. Diff between 'got' and 'expected' \n%s", name, diff)
		}
		return nil
	}
}
