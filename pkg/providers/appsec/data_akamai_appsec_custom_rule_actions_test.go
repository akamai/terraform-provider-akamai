package appsec

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiCustomRuleActions_basic(t *testing.T) {
	dataSourceName := "data.appsec_akamai_custom_rule_actions.akamaicustomruleactions"

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


resource "akamai_appsec_export_config" "appsecexport" {
    name = "Akamai Tools"
}


`
}

func testAccCheckAkamaiCustomRuleActionsDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai CustomRuleActions] CustomRuleActions Destroy skipped ")
	return nil
}
