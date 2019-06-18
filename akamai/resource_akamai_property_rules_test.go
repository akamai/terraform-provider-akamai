package akamai

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAkamaiPropertyRulesConfig = fmt.Sprintf(`
provider "akamai" {
  papi_section = "papi"
}

output "json" {
	value = "${akamai_property_rules.rules.json}"
}

resource "akamai_property_rules" "rules" {
 	rules {
		behavior {
			name = "origin"
        	option { 
       			key ="cacheKeyHostname"
            	value = "ORIGIN_HOSTNAME"
        	}
			option { 
    			key ="compress"
     			value = true
     		}
    		option { 
    			key ="enableTrueClientIp"
     			value = false
     		}
    		option { 
    			key ="forwardHostHeader"
     			value = "REQUEST_HOST_HEADER"
     		}
    		option { 
    			key ="hostname"
     			value = "exampleterraform.io"
     		}
    		option { 
    			key ="httpPort"
     			value = 80
     		}
    		option { 
    			key ="httpsPort"
     			value = 443
     		}
    		option { 
    			key ="originSni"
     			value = true
     		}
    		option { 
    			key ="originType"
     			value = "CUSTOMER"
     		}
    		option { 
    			key ="verificationMode"
     			value = "PLATFORM_SETTINGS"
     		}
    		option { 
    			key ="originCertificate"
     			value = ""
     		}
    		option { 
    			key ="ports"
     			value = ""
     		}
      	}
		behavior {
			name ="cpCode"
			option {
				key ="id"
				value = "cp-code-id"
			}
		}
		behavior {
			name ="caching"
			option {
				key ="behavior"
				value = "MAX_AGE"
			}
			option {
                key ="mustRevalidate"
                value = "false"
			}
            option {
                key ="ttl"
                value = "1d"
            }
		}
    }
}
`)

func TestAccAkamaiPropertyRules_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiPropertyRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiPropertyRulesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"akamai_property_rules.rules", "json",
					),
				),
			},
		},
	})
}

func testAccCheckAkamaiPropertyRulesDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_property_rules" {
			continue
		}

		rules := rs.Primary.Attributes["rules.#"]
		log.Printf("[DEBUG] [Akamai PropertyRules] Delete Rules [%s]", rules)

	}
	return nil
}
