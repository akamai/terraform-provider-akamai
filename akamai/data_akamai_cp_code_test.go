package akamai

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
  papi_section = "papi"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

data  "akamai_cp_code" "test" {
    name = "terraform-testing"
    contract = "${data.akamai_contract.contract.id}"
    group = "${data.akamai_group.group.id}"
}
`
}

func testAccCheckDataSourceCPCodeDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai CP] CP code Destroy skipped ")
	return nil
}
