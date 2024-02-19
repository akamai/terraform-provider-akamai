package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestDataPropertyRulesBuilder(t *testing.T) {
	t.Run("valid rule with 3 children - v2023-01-05", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/rules_v2023_01_05.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2023-05-30", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/rules_v2023_05_30.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-05-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/default_v2023_05_30.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-05-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/content_compression_v2023_05_30.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2023-05-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/static_content_v2023_05_30.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2023-05-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/dynamic_content_v2023_05_30.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2023-09-20", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/rules_v2023_09_20.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-09-20"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/default_v2023_09_20.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-09-20"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/content_compression_v2023_09_20.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2023-09-20"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/static_content_v2023_09_20.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2023-09-20"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/dynamic_content_v2023_09_20.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2023-10-30", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/rules_v2023_10_30.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-10-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/default_v2023_10_30.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-10-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/content_compression_v2023_10_30.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2023-10-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/static_content_v2023_10_30.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2023-10-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/dynamic_content_v2023_10_30.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2024-01-09", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/rules_v2024_01_09.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2024-01-09"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/default_v2024_01_09.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2024-01-09"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/content_compression_v2024_01_09.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2024-01-09"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/static_content_v2024_01_09.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2024-01-09"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/dynamic_content_v2024_01_09.json")),
					),
				}},
			})
		})
	})
	t.Run("invalid rule with 3 children with different versions", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/rules_mixed_versions.tf"),
					ExpectError: regexp.MustCompile(`child rule is using different rule format \(rules_v2023_05_30\) than expected \(rules_v2023_01_05\)`),
				}},
			})
		})
	})
	t.Run("fails on rule with more than one behavior in one block", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/rules_error_too_many_elements.tf"),
					ExpectError: regexp.MustCompile(`expected 1 element\(s\), got 2`),
				}},
			})
		})
	})
	t.Run("fails on rule with is_secure outside default rule", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/rules_with_is_secure_outside_default.tf"),
					ExpectError: regexp.MustCompile(`cannot be used outside 'default' rule: is_secure`),
				}},
			})
		})
	})
	t.Run("fails on rule with variable outside default rule", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/rules_with_variable_outside_default.tf"),
					ExpectError: regexp.MustCompile(`cannot be used outside 'default' rule: variable`),
				}},
			})
		})
	})
	t.Run("valid rule with one child and some values are variables", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV5ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/rules_variables.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/default_variables.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/content_compression_variables.json")),
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
