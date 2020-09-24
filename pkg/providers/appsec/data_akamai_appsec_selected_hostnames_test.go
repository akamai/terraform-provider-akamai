package appsec

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiSelectedHostnames_data_basic(t *testing.T) {
	dataSourceName := "data.akamai_appsec_selected_hostnames.appsecselectedhostnames"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiSelectedHostnamesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiSelectedHostnamesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiSelectedHostnamesConfig() string {
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

data "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version  
}


`
}

func testAccCheckAkamaiSelectedHostnamesDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai SelectedHostnames] SelectedHostnames Destroy skipped ")
	return nil
}
