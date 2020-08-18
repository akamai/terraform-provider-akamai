package property

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					resource.TestCheckResourceAttr(dataSourceName, "name", "terraform-test-datasource"),
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
		papi_section = "papi"
	  }

data "akamai_property" "test" {
	name = "terraform-test-datasource"
	version = 1
}
`
}

func testAccCheckDataSourcePropertyDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for Property Delete skipped ")

	return nil
}
