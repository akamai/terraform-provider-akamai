package appsec

import (
	"testing"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiConfigurationClone_basic(t *testing.T) {
	//dataSourceName := "akamai_appsec_configuration_clone.appsecconfigurationclone"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAkamaiConfigurationCloneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiConfigurationCloneConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiConfigurationCloneExists,
					//resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiConfigurationCloneConfig() string {
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


resource "akamai_appsec_configuration_version_clone" "appsecconfigurationclone" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    create_from_version = data.akamai_appsec_configuration.appsecconfigedge.latest_version 
    rule_update  = false
   }


`
}

func testAccCheckAkamaiConfigurationCloneExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_configuration_version_clone" {
			continue
		}
		//rname := rs.Primary.ID
		ccresp := appsec.NewConfigurationCloneResponse()
		err := ccresp.GetConfigurationClone("TEST")

		if err != nil {
			return err
		}
	}
	return nil
}
