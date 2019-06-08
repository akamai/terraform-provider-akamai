package akamai

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"testing"
)

var testAccAkamaiPropertyVariablesConfig = fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "global"
}

resource "akamai_property_variables" "test" {
	variables {
	  variable {
		  name = "PMUSER_DRUPAL_ORIGIN"
		  value =  "master-7twti-3svdm.us.platform.sh"
		  description  =  "Platform.sh Drupal URL"
		  hidden =  true
		  sensitive =  false
	  }
	  variable {
		  name =  "PMUSER_MIDDLEMAN_ORIGIN"
		  value = "dac-static.akamaideveloper.com"
		  description = "Heroku MiddleMan URL"
		  hidden =  true
		  sensitive = false
	  }
	  variable {
		  name = "PMUSER_PAPI_JSON_ORIGIN"
		  value =  "protected-sands-33849.herokuapp.com"
		  description = "Heroku"
		  hidden =  true
		  sensitive =  false
	  }
	  variable {
		  name = "PMUSER_TEST1"
		  value =  "prot.herokuapp.com"
		  description = "Heroku_test"
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
					resource.TestCheckResourceAttr(
						"akamai_property_variables.test", "variables.#", "1",
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