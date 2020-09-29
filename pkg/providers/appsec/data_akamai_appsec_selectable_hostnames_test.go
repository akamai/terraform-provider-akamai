package appsec

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiSelectableHostnames_data_basic(t *testing.T) {
	dataSourceName := "data.akamai_appsec_selectable_hostnames.appsecselectablehostnames"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiSelectableHostnamesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiSelectableHostnamesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiSelectableHostnamesConfig() string {
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

data "akamai_appsec_selectable_hostnames" "appsecselectablehostnames" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version   
}

output "selectablehostnames" {
  value = data.akamai_appsec_selectable_hostnames.appsecselectablehostnames.hostnames
}

`
}

func testAccCheckAkamaiSelectableHostnamesDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai SelectableHostnames] SelectableHostnames Destroy skipped ")
	return nil
}
