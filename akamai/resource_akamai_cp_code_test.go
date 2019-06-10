package akamai

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAkamaiCPCode_basic(t *testing.T) {
	dataSourceName := "akamai_cp_code.cp_code"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiCPCodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiCpCodeConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccAkamaiCpCodeConfig() string {
	return `
provider "akamai" {
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
`
}

func testAccCheckAkamaiCPCodeDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai CP] CP code Destroy skipped ")
	return nil
}
