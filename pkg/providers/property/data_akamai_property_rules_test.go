package property

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDSPropertyRules(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			Config:      loadFixtureString("testdata/TestDSPropertyRules/always_fails.tf"),
			ExpectError: regexp.MustCompile(`data "akamai_property_rules" is no longer supported`),
		}},
	})
}
