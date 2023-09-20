package cps

import (
	"testing"

	"github.com/akamai/terraform-provider-akamai/v5/pkg/common/testutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataWarnings(t *testing.T) {
	t.Run("run warning datasource", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProviders,
			IsUnitTest:               true,
			Steps: []resource.TestStep{
				{
					Config: testutils.LoadFixtureString(t, "testdata/TestDataWarnings/warnings.tf"),
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
