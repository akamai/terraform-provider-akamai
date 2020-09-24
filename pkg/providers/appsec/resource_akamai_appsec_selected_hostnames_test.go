package appsec

import (
	"strconv"
	"testing"

	appsec "github.com/akamai/AkamaiOPEN-edgegrid-golang/appsec-v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAkamaiSelectedHostnames_basic(t *testing.T) {

	dataSourceName := "akamai_appsec_selected_hostnames.appsecselectedhostnames"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAkamaiSelectedHostnamesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiSelectedHostnamesNConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiSelectedHostnamesExists,
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiSelectedHostnamesNConfig() string {
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

resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version 
    hostnames = ["rinaldi.sandbox.akamaideveloper.com","sujala.sandbox.akamaideveloper.com"]  
}


`
}

func testAccCheckAkamaiSelectedHostnamesExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_appsec_selected_hostnames" {
			continue
		}
		//rname := rs.Primary.ID
		configid, _ := strconv.Atoi(rs.Primary.Attributes["config_id"])
		version, _ := strconv.Atoi(rs.Primary.Attributes["version"])
		ccresp := appsec.NewSelectedHostnamesResponse()
		err := ccresp.GetSelectedHostnames(configid, version, "TEST")

		if err != nil {
			return err
		}
	}

	return nil
}
