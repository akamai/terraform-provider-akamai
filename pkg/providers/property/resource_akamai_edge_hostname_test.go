package property

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccAkamaiSecureEdgeHostNameConfig = fmt.Sprintf(`
provider "akamai" {
  papi_section = "papi"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_edge_hostname" "test" {
    product = "prd_SPM"
    contract = "${data.akamai_contract.contract.id}"
    group = "${data.akamai_group.group.id}"
    edge_hostname =  "terraform-test.exampleterraform.io.edgesuite.net"
	ip_behavior = "IPV4"
}
`)

func TestAccAkamaiSecureEdgeHostName_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiSecureEdgeHostNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiSecureEdgeHostNameConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiSecureEdgeHostNameExists,
				),
			},
		},
	})
}

func testAccCheckAkamaiSecureEdgeHostNameDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_edge_hostname" {
			continue
		}
		log.Printf("[DEBUG] [Akamai SecureEdgeHostName] Delete for edgehostname not supported  [%v]", rs.Primary.ID)
	}
	return nil
}

func testAccCheckAkamaiSecureEdgeHostNameExists(s *terraform.State) error {
	// TODO: rewrite for v2 client??
	return nil
}
