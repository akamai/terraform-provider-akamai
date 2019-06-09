package akamai

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceCPCode_basic(t *testing.T) {
	dataSourceName := "data.akamai_cp_code.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourceCPCodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCPCode_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceCPCode_basic() string {
	return `
provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "papi"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_cp_code" "cp_code" {
	name = "terraform-testing"
	contract = "${data.akamai_contract.contract.id}"
	group = "${data.akamai_group.group.id}"
	product = "prd_SPM"
}

data  "akamai_cp_code" "test" {
	depends_on = ["akamai_cp_code.cp_code"]
    name = "terraform-testing"
    contract = "${akamai_contract.contract.id}"
    group = "${akamai_contract.group.id}"
}
`
}

func testAccCheckDataSourceCPCodeDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai CP] CP code Destroy skipped ")
	return nil
}
