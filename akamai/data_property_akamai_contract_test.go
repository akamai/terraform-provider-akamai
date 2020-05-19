package akamai

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceContract_basic(t *testing.T) {
	dataSourceName := "data.akamai_contract.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceContract_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceContract_basic() string {
	return `
provider "akamai" {
  papi_section = "papi"
}

data "akamai_contract" "test" {
}

output "contractid" {
  value = "${data.akamai_contract.test.id}"
}
`
}

func testAccCheckAkamaiContractDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for Contract Delete skipped ")

	return nil
}
