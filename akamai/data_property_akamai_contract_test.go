package akamai

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"testing"
)

func TestAccDataSourceContract_basic(t *testing.T) {
	dataSourceName := "data.akamai_contract.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceContract_basic("Davey Shafik"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "ctr_C-1FRYVV3"),
				),
			},
		},
	})
}

func testAccDataSourceContract_basic(name string) string {
	return fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  dns_section = "papi"
}

data "akamai_contract" "test" {
    name = "%s"
}

output "contractid" {
  value = "${data.akamai_contract.test.id}"
}

`, name)
}

func testAccCheckAkamaiContractDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for Contract Delete skipped ")

	return nil
}
