package akamai

import (
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAkamaiPropertyConfig = fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "global"
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

func TestAccAkamaiProperty_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiPropertyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiPropertyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiPropertyExists,
				),
			},
		},
	})
}

func testAccCheckAkamaiPropertyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_property" {
			continue
		}

		property := papi.NewProperty(papi.NewProperties())
		property.PropertyID = rs.Primary.ID
		e := property.GetProperty()
		if e != nil {
			ee, ok := e.(client.APIError)
			if ok && ee.Status == 403 {
				return nil
			}
			return e
		}
	}
	return nil
}

func testAccCheckAkamaiPropertyExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_property" {
			continue
		}

		property := papi.NewProperty(papi.NewProperties())
		property.PropertyID = rs.Primary.ID
		e := property.GetProperty()
		if e != nil {
			return e
		}
	}
	return nil
}
