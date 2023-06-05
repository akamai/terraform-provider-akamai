package property

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResPropertyVariables(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{{
			Config:      loadFixtureString("testdata/TestResPropertyVariables/schema_version1_always_fails.tf"),
			ExpectError: regexp.MustCompile(`resource "akamai_property_variables" is no longer supported`),
		}},
	})
}
