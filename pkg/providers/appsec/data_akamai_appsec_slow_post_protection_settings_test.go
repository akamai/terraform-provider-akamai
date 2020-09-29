package appsec

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiSlowPostProtectionSettings_data_basic(t *testing.T) {
	dataSourceName := "data.akamai_appsec_slow_post_protection_settings.appsecslowpostprotectionsettings"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiSlowPostProtectionSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiSlowPostProtectionSettingsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiSlowPostProtectionSettingsConfig() string {
	return `
provider "akamai" {
  appsec_section = "default"
}
data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
  
}



output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}


data "akamai_appsec_slow_post" "appsecslowpostprotectionsettings" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
}

output "configsedge_post_output_text" {
  value = data.akamai_appsec_slow_post.appsecslowpostprotectionsettings.output_text
}

`
}

func testAccCheckAkamaiSlowPostProtectionSettingsDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai SlowPostProtectionSettings] SlowPostProtectionSettings Destroy skipped ")
	return nil
}
