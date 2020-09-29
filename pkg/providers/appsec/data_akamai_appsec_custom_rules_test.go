package appsec

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiCustomRules_data_basic(t *testing.T) {
	dataSourceName := "data.akamai_appsec_custom_rules.appseccustomrules"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiCustomRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiCustomRulesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiCustomRulesConfig() string {
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


data "akamai_appsec_custom_rules" "appseccustomrule" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
}

output "appseccustomrules" {
  value = data.akamai_appsec_custom_rules.appseccustomrule.output_text
}

`
}

func testAccCheckAkamaiCustomRulesDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai CustomRules] CustomRules Destroy skipped ")
	return nil
}
