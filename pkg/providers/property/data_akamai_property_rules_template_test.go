package property

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestDataAkamaiPropertyRulesRead(t *testing.T) {
	t.Run("valid nested template with vars", func(t *testing.T) {
		client := mockpapi{}
		useClient(&client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSRulesTemplate/template.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_rules_template.test", "json", loadFixtureString("testdata/TestDSRulesTemplate/rules.json")),
						),
					},
				},
			})
		})
	})

	t.Run("recursive templates", func(t *testing.T) {
		client := mockpapi{}
		useClient(&client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config:      loadFixtureString("testdata/TestDSRulesTemplate/template-circular.tf"),
						ExpectError: regexp.MustCompile("exceeded maximum template depth"),
					},
				},
			})
		})
	})
}
