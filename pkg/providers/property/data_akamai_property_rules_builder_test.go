package property

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestDataPropertyRulesBuilder(t *testing.T) {
	t.Run("valid rule with 3 children", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSPropertyRulesBuilder/rules.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.default",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.default",
							"json",
							loadFixtureString("testdata/TestDSPropertyRulesBuilder/default.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.content_compression",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.content_compression",
							"json",
							loadFixtureString("testdata/TestDSPropertyRulesBuilder/content_compression.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.static_content",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.static_content",
							"json",
							loadFixtureString("testdata/TestDSPropertyRulesBuilder/static_content.json")),

						resource.TestCheckResourceAttr("data.akamai_property_rules_builder.dynamic_content",
							"rule_format",
							"v2023-01-05"),
						testCheckResourceAttrJSON("data.akamai_property_rules_builder.dynamic_content",
							"json",
							loadFixtureString("testdata/TestDSPropertyRulesBuilder/dynamic_content.json")),
					),
				}},
			})
		})
	})
	t.Run("fails on rule with more than one behavior in one block", func(t *testing.T) {
		useClient(nil, nil, func() {
			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestDSPropertyRulesBuilder/rules_error_too_many_elements.tf"),
					ExpectError: regexp.MustCompile(`expected 1 element\(s\), got 2`),
				}},
			})
		})
	})
}

func testCheckResourceAttrJSONFromFixture(name, key, fixturePath string) func(s *terraform.State) error {
	value := loadFixtureString(fixturePath)
	checkFn := testCheckResourceAttrJSON(name, key, value)

	return func(s *terraform.State) error {
		if err := checkFn(s); err != nil {
			return fmt.Errorf("expected from: %s, %w", fixturePath, err)
		}
		return nil
	}
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
