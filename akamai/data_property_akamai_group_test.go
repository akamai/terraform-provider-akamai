package akamai

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceGroup_basic(t *testing.T) {
	dataSourceName := "data.akamai_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroup_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceGroup_basic() string {
	return `
provider "akamai" {
  papi_section = "papi"
}

data "akamai_group" "test" {
}

output "groupid" {
value = "${data.akamai_group.test.id}"
}
`
}

func testAccCheckAkamaiGroupDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for Group Delete skipped ")

	return nil
}
