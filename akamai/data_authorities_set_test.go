package akamai

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceAuthoritiesSet_basic(t *testing.T) {
	dataSourceName := "data.akamai_authorities_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAuthoritiesSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAuthoritiesSet_basic("C-1FRYVV3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "C-1FRYVV3"),
				),
			},
		},
	})
}

func testAccDataSourceAuthoritiesSet_basic(name string) string {
	return fmt.Sprintf(`
provider "akamai" {
  dns_section = "dns"
}

data "akamai_authorities_set" "test" {
	contract = "%s"
  }
  
  
  output "authorities" {
	value = "${join(",", data.akamai_authorities_set.test.authorities)}"
  }

`, name)
}

func testAccCheckAuthoritiesSetDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for AuthoritiesSet Delete skipped ")

	return nil
}
