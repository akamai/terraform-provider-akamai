package deprecated

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceAuthoritiesSet_basic(t *testing.T) {
	dataSourceName := "data.akamai_authorities_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuthoritiesSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAuthoritiesSet_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceAuthoritiesSet_basic() string {
	return `
provider "akamai" {
  papi_section = "dns"
  dns_section = "dns"
}

data "akamai_contract" "contract" { }

data "akamai_authorities_set" "test" {
	contract = "${data.akamai_contract.contract.id}"
}
  
output "authorities" {
	value = "${join(",", data.akamai_authorities_set.test.authorities)}"
}`
}

func testAccCheckAuthoritiesSetDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for AuthoritiesSet Delete skipped ")

	return nil
}
