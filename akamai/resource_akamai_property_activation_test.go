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
	property = "${akamai_property.property.id}"
	version = "${akamai_property.property.version}"
	network = "STAGING"
	activate = true
	contact = ["dshafik@akamai.com"]
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_cp_code" "cp_code" {
	name = "terraform-testing"
	contract = "${data.akamai_contract.contract.id}"
	group = "${data.akamai_group.group.id}"
	product = "prd_SPM"
}

resource "akamai_edge_hostname" "test" {
    product = "prd_SPM"
    contract = "${data.akamai_contract.contract.id}"
    group = "${data.akamai_group.group.id}"
    edge_hostname =  "terraform-test.example.org.edgesuite.net"
    ipv6 = true
}

resource "akamai_property" "property" {
  name = "terraform-test"

  contact = ["user@example.org"]

  product = "prd_SPM"
  cp_code = "${akamai_cp_code.cp_code.id}"
  contract = "${data.akamai_contract.contract.id}"
  group = "${data.akamai_group.group.id}"

  hostnames = "${merge(akamai_edge_hostname.test.hostnames)}"
  
  rule_format = "v2016-11-15"
  
  rules = "${akamai_property_rules.rules.json}"
}

resource "akamai_property_rules" "rules" {
 	rules {
		behavior {
			name =  "origin"
        	option { 
       			key =  "cacheKeyHostname"
            	value = "ORIGIN_HOSTNAME"
        	}
			option { 
    			key =  "compress"
     			value = "true"
     		}
    		option { 
    			key =  "enableTrueClientIp"
     			value = "false"
     		}
    		option { 
    			key =  "forwardHostHeader"
     			value = "REQUEST_HOST_HEADER"
     		}
    		option { 
    			key =  "hostname"
     			value = "example.org"
     		}
    		option { 
    			key =  "httpPort"
     			value = "80"
     		}
    		option { 
    			key =  "httpsPort"
     			value = "443"
     		}
    		option { 
    			key =  "originSni"
     			value = "true"
     		}
    		option { 
    			key =  "originType"
     			value = "CUSTOMER"
     		}
    		option { 
    			key =  "verificationMode"
     			value = "PLATFORM_SETTINGS"
     		}
    		option { 
    			key =  "originCertificate"
     			value = ""
     		}
    		option { 
    			key =  "ports"
     			value = ""
     		}
      	}
		behavior {
			name =  "cpCode"
			option {
				key =  "id"
				value = "${akamai_cp_code.cp_code.id}"
			}
			option {
				key =  "name"
				value = "${akamai_cp_code.cp_code.name}"
			}
		}
		behavior {
			name =  "caching"
			option {
				key =  "behavior"
				value = "MAX_AGE"
			}
			option {
                key =  "mustRevalidate"
                value = "false"
			}
            option {
                key =  "ttl"
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
		propertyID := rs.Primary.Attributes["property"]

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
