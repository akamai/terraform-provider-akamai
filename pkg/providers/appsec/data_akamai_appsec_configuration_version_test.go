package appsec

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiConfigurationVersion_data_basic(t *testing.T) {
	dataSourceName := "data.akamai_appsec_configuration_version.appsecconfigurationversion"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiConfigurationVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiConfigurationVersionConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiConfigurationVersionConfig() string {
	return `
provider "akamai" {
  appsec_section = "default"
}


data "akamai_appsec_configuration_version" "appsecconfigurationversion" {
    name = "Akamai Tools"
   }

output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}

output "configsedgelatestversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}

output "configsedgeconfiglist" {
  value = data.akamai_appsec_configuration.appsecconfigedge.output_text
}

output "configsedgeconfigversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.version
}

`
}

func testAccCheckAkamaiConfigurationVersionDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai ConfigurationVersion] ConfigurationVersion Destroy skipped ")
	return nil
}
