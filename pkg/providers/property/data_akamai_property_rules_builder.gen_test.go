package property

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v6/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Do not modify those tests manually.
func TestDataPropertyRulesBuilderGenerated(t *testing.T) {
	t.Run("valid rule with 3 children - v2023-01-05", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_01_05/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_01_05/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_01_05/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_01_05/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_01_05/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2023-05-30", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_05_30/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-05-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_05_30/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-05-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_05_30/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2023-05-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_05_30/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2023-05-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_05_30/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2023-09-20", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_09_20/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-09-20"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_09_20/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-09-20"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_09_20/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2023-09-20"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_09_20/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2023-09-20"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_09_20/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2023-10-30", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_10_30/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-10-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_10_30/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-10-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_10_30/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2023-10-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_10_30/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2023-10-30"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2023_10_30/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2024-01-09", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_01_09/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2024-01-09"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_01_09/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2024-01-09"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_01_09/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2024-01-09"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_01_09/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2024-01-09"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_01_09/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2024-02-12", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_02_12/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2024-02-12"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_02_12/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2024-02-12"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_02_12/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2024-02-12"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_02_12/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2024-02-12"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_02_12/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2024-05-31", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_05_31/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2024-05-31"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_05_31/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2024-05-31"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_05_31/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2024-05-31"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_05_31/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2024-05-31"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_05_31/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2024-08-13", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_08_13/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2024-08-13"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_08_13/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2024-08-13"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_08_13/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2024-08-13"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_08_13/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2024-08-13"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_08_13/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2024-10-21", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_10_21/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2024-10-21"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_10_21/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2024-10-21"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_10_21/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2024-10-21"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_10_21/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2024-10-21"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2024_10_21/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("valid rule with 3 children - v2025-01-13", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testutils.NewProtoV6ProviderFactory(NewSubprovider()),
				Steps: []resource.TestStep{{
					Config: testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2025_01_13/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2025-01-13"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2025_01_13/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2025-01-13"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2025_01_13/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2025-01-13"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2025_01_13/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2025-01-13"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							testutils.LoadFixtureString(t, "testdata/TestDSPropertyRulesBuilder/ruleformat/v2025_01_13/dynamic_content.json")),
					),
				}},
			})
		})
	})
}
