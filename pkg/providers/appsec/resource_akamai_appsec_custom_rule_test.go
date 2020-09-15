package appsec

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiCustomRule_basic(t *testing.T) {
	dataSourceName := "akamai_appsec_custom_rule.appseccustomrule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiCustomRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiCustomRuleConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func TestAccAkamaiCustomRule_update(t *testing.T) {
	dataSourceName := "akamai_appsec_custom_rule.appseccustomrule"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiCustomRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiCustomRuleConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					testAccCheckAkamaiCustomRuleExists,
					//resource.TestCheckResourceAttr("akamai_appsec_custom_rule.appseccustomrule", "load_imbalance_percentage", "50"),
				),
			},
			{
				Config: testAccAkamaiCustomRuleUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					testAccCheckAkamaiCustomRuleExists,
					//resource.TestCheckResourceAttr("akamai_appsec_custom_rule.appseccustomrule", "load_imbalance_percentage", "50"),
				),
			},
		},
	})
}

func testAccAkamaiCustomRuleConfig() string {
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


resource "akamai_appsec_custom_rule" "appseccustomrule" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    custom_rules_json =  file("${path.module}/custom_rules.json")
}

`
}

func testAccAkamaiCustomRuleUpdateConfig() string {
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


resource "akamai_appsec_custom_rule" "appseccustomrule" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
}

`
}

func testCheckDeleteCustomRuleResource(s *terraform.State, rscName string) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_custom_rule" {
			continue
		}

		ccresp := appsec.NewCustomRuleResponse()

		configid, _ := strconv.Atoi(rs.Primary.Attributes["config_id"])
		ruleid, _ := strconv.Atoi(rs.Primary.ID)

		err := ccresp.GetCustomRule(configid, ruleid, "TEST")

		if err != nil {
			return err
		}
		log.Printf("[DEBUG] [Test] Deleting test resource [%v]", rscName)

		err = ccresp.DeleteCustomRule(configid, ruleid, "TEST")
		if err != nil {
			return fmt.Errorf("resource was not deleted %s. Error: %s", rscName, err.Error())
		}
	}

	return nil
}

func testAccCheckAkamaiCustomRuleDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_custom_rule" {
			continue
		}

		rscName := "akamai_appsec_custom_rule"
		if err := testCheckDeleteCustomRuleResource(s, rscName); err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckAkamaiCustomRuleExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_custom_rule" {
			continue
		}

		//rname := rs.Primary.ID

		ccresp := appsec.NewCustomRuleResponse()
		configid, _ := strconv.Atoi(rs.Primary.Attributes["config_id"])
		ruleid, _ := strconv.Atoi(rs.Primary.ID)
		err := ccresp.GetCustomRule(configid, ruleid, "TEST")

		if err != nil {
			return err
		}
	}

	return nil
}
