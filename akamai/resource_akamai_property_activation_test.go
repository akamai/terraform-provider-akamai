package akamai

import (
    "fmt"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
    "log"
    "strings"
	"testing"
)

var testAccAkamaiPropertyActivationConfig = fmt.Sprintf(`
provider "akamai" {
    edgerc = "~/.edgerc"
    papi_section = "papi_section"
}

variable "activate" {
   default = true
}

resource "akamai_property_activation" "dshafik_sandbox" {

       name = "akavadeveloper.com"
       contact = ["martin@akava.io"]
       hostname = ["akavadeveloper.com"]
       contract = "${data.akamai_contract.our_contract.id}"
       group =  "${data.akamai_group.our_group.id}"
       network = "STAGING"
       activate = "${var.activate}"
}

data "akamai_group" "our_group" {
   name = "Davey Shafik"
}

output "groupid" {
 value = "${data.akamai_group.our_group.id}"
}


data "akamai_contract" "our_contract" {
   name = "Davey Shafik"
}

output "contractid" {
 value = "${data.akamai_contract.our_contract.id}"
}

`)

func TestAccAkamaiPropertyActivation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiPropertyActivationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiPropertyActivationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiPropertyActivationExists,
				),
			},
		},
	})
}

func testAccCheckAkamaiPropertyActivationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_property_activation" {
			continue
		}
				
		log.Printf("[DEBUG] [Akamai PropertyActivation] Activation Delete")
				
	}
	return nil
}

func testAccCheckAkamaiPropertyActivationExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_property_activation" {
			continue
		}

		
	propertyID := rs.Primary.ID

	property := papi.NewProperty(papi.NewProperties())
	property.PropertyID = propertyID
	property.Contract = &papi.Contract{ContractID: rs.Primary.Attributes["contract"]}
	property.Group = &papi.Group{GroupID: rs.Primary.Attributes["group"]}

	e := property.GetProperty()
	if e != nil {
		return e
	}

	log.Printf("[DEBUG] GET ACTIVATION PROPERTY %v", property)

	activations, e := property.GetActivations()
	if e != nil {
		return e
	}
	log.Printf("[DEBUG] ACTIVATION activations %v", activations)
	activation, e := activations.GetLatestActivation(papi.NetworkValue(strings.ToUpper(rs.Primary.Attributes["network"])), papi.StatusActive)
	log.Printf("[DEBUG] ACTIVATION activations get latest %v", activations)
    log.Printf("[DEBUG] ACTIVATION activation get latest %v", activation)
	}
	return nil
}
