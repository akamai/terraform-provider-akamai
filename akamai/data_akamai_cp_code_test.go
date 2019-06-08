package akamai

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"testing"
)

func TestAccDataSourceCPCode_basic(t *testing.T) {
	dataSourceName := "data.akamai_cp_code.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiCPCodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCPCode_basic("www.example.org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "cpc_705684"),
				),
			},
		},
	})
}

func testAccDataSourceCPCode_basic(name string) string {
	return fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  dns_section = "global"
}

data  "akamai_cp_code" "test" {
    name = "%s"
    contract = "ctr_C-1FRYVV3"
    group = "grp_68817"

}

`, name)
}

func testAccCheckAkamaiCPCodeDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Akamai CP] Searching for CP code Delete skipped ")

	return nil
}
