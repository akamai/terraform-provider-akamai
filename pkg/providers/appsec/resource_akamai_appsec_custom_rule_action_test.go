package appsec

import (
	"strconv"
	"testing"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiCustomRuleAction_basic(t *testing.T) {

	dataSourceName := "akamai_appsec_custom_rule_action.appseccustomruleaction"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAkamaiCustomRuleActionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiCustomRuleActionConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiCustomRuleActionExists,
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiCustomRuleActionConfig() string {
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


resource "akamai_appsec_custom_rule_action" "appsecreatecustomruleaction" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "f1rQ_106946"
    rule_id = 321456
    custom_rule_action = "alert"
}

output "customruleaction" {
  value = akamai_appsec_custom_rule_action.appsecreatecustomruleaction.rule_id
}


`
}

func testAccCheckAkamaiCustomRuleActionExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_custom_rule_action" {
			continue
		}
		//rname := rs.Primary.ID
		configid, _ := strconv.Atoi(rs.Primary.Attributes["config_id"])
		version, _ := strconv.Atoi(rs.Primary.Attributes["version"])
		policyid := rs.Primary.Attributes["policy_id"]
		ccresp := appsec.NewCustomRuleActionResponse()
		err := ccresp.GetCustomRuleAction(configid, version, policyid, "TEST")

		if err != nil {
			return err
		}
	}

	return nil
}
