package akamai

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"testing"
)

func TestAccDataSourceGroup_basic(t *testing.T) {
	dataSourceName := "data.akamai_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroup_basic("Davey Shafik"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "grp_68817"),
				),
			},
		},
	})
}

func testAccDataSourceGroup_basic(name string) string {
	return fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  dns_section = "papi"
}

data "akamai_group" "test" {
	name = "%s"
}

output "groupid" {
value = "${data.akamai_group.test.id}"
}

`, name)
}

func testAccCheckAkamaiGroupDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for Group Delete skipped ")

	return nil
}
