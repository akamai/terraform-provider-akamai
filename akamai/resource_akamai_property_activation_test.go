package akamai

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAkamaiPropertyActivationConfig = fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "papi"
}

resource "akamai_property_activation" "property_activation" {
	property = "${akamai_property.akamai_developer.id}"
	version = "${akamai_property.akamai_developer.version}"
	network = "STAGING"
	activate = true
	contact = ["user@example.org"]
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_cp_code" "cp_code" {
	name = "terraform-testing"
	contract = "${akamai_contract.contract.id}"
	group = "${akamai_group.group.id}"
	product = "prd_SPM"
}

resource "random_pet" "property_name" {
}

resource "akamai_property" "akamai_developer" {
  name = "${random_pet.property_name.id}"

  contact = ["user@example.org"]

  product = "prd_SPM"
  cp_code = "${akamai_cp_code.cp_code.id}"
  contract = "${data.akamai_contract.contract.id}"
  group = "${data.akamai_group.group.id}"
  
  rule_format = "v2016-11-15"
  
  rules = "${akamai_property_rules.rules.json}"
}

resource "akamai_property_rules" "rules" {
 	rules {
		behavior {
			name = "origin"
        	option { 
       			name = "cacheKeyHostname"
            	value = "ORIGIN_HOSTNAME"
        	}
			option { 
    			name = "compress"
     			value = true
     		}
    		option { 
    			name = "enableTrueClientIp"
     			value = false
     		}
    		option { 
    			name = "forwardHostHeader"
     			value = "REQUEST_HOST_HEADER"
     		}
    		option { 
    			name = "hostname"
     			value = "example.org"
     		}
    		option { 
    			name = "httpPort"
     			value = 80
     		}
    		option { 
    			name = "httpsPort"
     			value = 443
     		}
    		option { 
    			name = "originSni"
     			value = true
     		}
    		option { 
    			name = "originType"
     			value = "CUSTOMER"
     		}
    		option { 
    			name = "verificationMode"
     			value = "PLATFORM_SETTINGS"
     		}
    		option { 
    			name = "originCertificate"
     			value = ""
     		}
    		option { 
    			name = "ports"
     			value = ""
     		}
      	}
		behavior {
			name = "cpCode"
			option {
				name = "id"
				value = "${akamai_cp_code.cp_code.id}"
			}
			option {
				name = "name"
				value = "${akamai_cp_code.cp_code.name}"
			}
		}
		behavior {
			name = "caching"
			option {
				name = "behavior"
				value = "MAX_AGE"
			}
			option {
                name = "mustRevalidate"
                value = "false"
			}
            option {
                name = "ttl"
                value = "1d"
            }
		}
    }
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
