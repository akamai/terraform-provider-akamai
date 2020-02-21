package akamai

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceProperty_basic(t *testing.T) {
	dataSourceName := "data.akamai_property.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDataSourcePropertyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProperty_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", "test"),
					resource.TestCheckResourceAttr(dataSourceName, "version", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules"),
				),
			},
		},
	})
}

func testAccDataSourceProperty_basic() string {
	return `
provider "akamai" {
}

data "akamai_property" "test" {
	name = "test"
	version = 1
}
`
}

func testAccCheckDataSourcePropertyDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for Property Delete skipped ")

	return nil
}
