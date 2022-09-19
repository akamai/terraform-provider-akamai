package cps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataWarnings(t *testing.T) {
	t.Run("run warning datasource", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			Providers:  testAccProviders,
			IsUnitTest: true,
			Steps: []resource.TestStep{
				{
					Config: loadFixtureString("testdata/TestDataWarnings/warnings.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_cps_warnings.test", "warnings.%", "112"),
						resource.TestCheckResourceAttr("data.akamai_cps_warnings.test", "warnings.CERTIFICATE_NULL_OR_EMPTY", "Null or empty [<certificateDescription>] Certificate."),
						resource.TestCheckNoResourceAttr("data.akamai_cps_warnings.test", "warnings.a"),
					),
				},
			},
		})
	})
}
