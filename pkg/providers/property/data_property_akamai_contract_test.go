package property

import (
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceContract_basic(t *testing.T) {
	dataSourceName := "data.akamai_contract.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiContractDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceContractBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
				),
			},
		},
	})
}

func testAccDataSourceContractBasic() string {
	return `
provider "akamai" {
  papi_section = "papi"
}

data "akamai_contract" "test" {
}

output "contractid" {
  value = "${data.akamai_contract.test.id}"
}
`
}

func testAccCheckAkamaiContractDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] [Group] Searching for Contract Delete skipped ")

	return nil
}

func TestDataSourceSchema(t *testing.T) {

	t.Run("contract collides with contract ID", func(t *testing.T) {
		resource.UnitTest(t, resource.TestCase{
			Providers:  testAccProviders,
			IsUnitTest: true,
			Steps: []resource.TestStep{{
				Config:             testDataSourceSchema(),
				ExpectNonEmptyPlan: true,
				ExpectError:        regexp.MustCompile("only one of `group,group_id,group_name` can be specified"),
			}},
		})
	})
}

func testDataSourceSchema() string {
	return `
provider "akamai" {
  papi_section = "papi"
}

data "akamai_contract" "test" {
  group = "grp_test"
  group_id = "grp_test"
  group_name = "grp_test"
}
`
}
