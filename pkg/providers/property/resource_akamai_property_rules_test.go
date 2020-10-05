package property

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccAkamaiPropertyRulesConfig = `
provider "akamai" {
  papi_section = "papi"
  edgerc = "~/.edgerc"
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
`

var errorRegex, err = regexp.Compile("The akamai_property_rules resource has moved to a data source, please change 'resource \"akamai_property_rules\"' to 'data \"akamai_property_rules\"")

func TestAccAkamaiPropertyRules_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiPropertyRulesDestroy,
		Steps: []resource.TestStep{
			{
				ExpectError: errorRegex,
				Config:      testAccAkamaiPropertyRulesConfig,
			},
		},
	})
}

var testAccAkamaiPropertyRulesSiteshield = `
provider "akamai" {
  edgerc = "~/.edgerc"
}

output "json" {
	value = "${akamai_property_rules.rules.json}"
}

resource "akamai_property_rules" "rules" {
 	rules {
		behavior { 
			name = "siteShield" 
			option { 
				key = "ssmap" 
				value = "mapname.akamai.net"
			}
		}
	}
}
`

func TestAkamaiPropertyRules_siteshield(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectError: errorRegex,
				Destroy:     false,
				Config:      testAccAkamaiPropertyRulesSiteshield,
			},
		},
	})
}

var testAccAkamaiPropertyRulesCPCode = `
provider "akamai" {
	  edgerc = "~/.edgerc"
}

output "json" {
	value = "${akamai_property_rules.rules.json}"
}

resource "akamai_property_rules" "rules" {
 	rules {
		behavior { 
			name = "cpCode" 
			option { 
				key = "id" 
				value = "cpc_12345"
			}
		}
	}
}
`

func TestAkamaiPropertyRules_cpCode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectError: errorRegex,
				Destroy:     false,
				Config:      testAccAkamaiPropertyRulesCPCode,
			},
		},
	})
}

var testAccAkamaiPropertyRulesIsSecure = `
provider "akamai" {
  edgerc = "~/.edgerc"
}

output "json" {
	value = "${akamai_property_rules.rules.json}"
}

resource "akamai_property_rules" "rules" {
 	rules {
		is_secure = %s
	}
}
`

func TestAkamaiPropertyRules_isSecureTrue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectError: errorRegex,
				Destroy:     false,
				Config:      fmt.Sprintf(testAccAkamaiPropertyRulesIsSecure, "true"),
			},
		},
	})
}

func TestAkamaiPropertyRules_isSecureFalse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectError: errorRegex,
				Destroy:     false,
				Config:      fmt.Sprintf(testAccAkamaiPropertyRulesIsSecure, "false"),
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
