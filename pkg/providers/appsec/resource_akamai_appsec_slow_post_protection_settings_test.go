package appsec

import (
	"strconv"
	"testing"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiSlowPostProtectionSettings_basic(t *testing.T) {

	dataSourceName := "akamai_appsec_slow_post_protection_settings.appsecslowpostprotectionsettings"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAkamaiSlowPostProtectionSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiSlowPostProtectionSettingsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiSlowPostProtectionSettingsExists,
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


resource "akamai_appsec_slow_post_protection_settings" "appsecslowpostprotectionsettings" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    slow_rate_action = "alert"                        
    slow_rate_threshold_rate = 10
    slow_rate_threshold_period = 60
    duration_threshold_timeout = 20
}

`
}

func testAccCheckAkamaiSlowPostProtectionSettingsExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_slow_post_protection_settings" {
			continue
		}
		//rname := rs.Primary.ID
		configid, _ := strconv.Atoi(rs.Primary.Attributes["config_id"])
		version, _ := strconv.Atoi(rs.Primary.Attributes["version"])
		policyid := rs.Primary.Attributes["policy_id"]
		ccresp := appsec.NewSlowPostProtectionSettingsResponse()
		err := ccresp.GetSlowPostProtectionSettings(configid, version, policyid, "TEST")

		if err != nil {
			return err
		}
	}

	return nil
}
