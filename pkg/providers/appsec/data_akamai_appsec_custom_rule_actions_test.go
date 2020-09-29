package appsec

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiCustomRuleActions_data_basic(t *testing.T) {
	dataSourceName := "data.akamai_appsec_custom_rule_actions.appseccustomruleactions"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiCustomRuleActionsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiCustomRuleActionsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiCustomRuleActionsConfig() string {
	return `
provider "akamai" {
  appsec_section = "default"
}


data "akamai_appsec_custom_rule_actions" "appsecreatecustomruleactions" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
}

output "customruleactions" {
  value = data.akamai_appsec_custom_rule_actions.appsecreatecustomruleactions.output_text
}


`
}

func testAccCheckAkamaiCustomRuleActionsDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai CustomRuleActions] CustomRuleActions Destroy skipped ")
	return nil
}
