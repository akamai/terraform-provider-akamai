package property

import
(
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccAkamaiPropertyVariablesConfig = fmt.Sprintf(`
provider "akamai" {
  papi_section = "papi"
}

resource "akamai_property_variables" "test" {
	variables {
	  	variable {
			name = "PMUSER_ORIGIN"
		  	value =  "origin.exampleterraform.io"
		  	description  =  "Example Origin"
		  	hidden =  true
		  	sensitive =  false
	  	}
	}
}
`)

func TestAccAkamaiPropertyVariables_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiPropertyVariablesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiPropertyVariablesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"akamai_property_variables.test", "json",
					),
				),
			},
		},
	})
}

func testAccCheckAkamaiPropertyVariablesDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_property_variables" {
			continue
		}

		hostname := rs.Primary.Attributes["variables.#"]
		log.Printf("[DEBUG] [Akamai PropertyVariables] Delete variable [%s]", hostname)

	}
	return nil
}
