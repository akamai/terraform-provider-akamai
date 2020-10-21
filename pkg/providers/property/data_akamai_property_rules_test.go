package property

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccAkamaiDataPropertyRulesConfig = `
provider "akamai" {
  papi_section = "papi"
}

output "json" {
	value = "${data.akamai_property_rules.rules.json}"
}

data "akamai_property_rules" "rules" {
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

func TestAccAkamaiDataPropertyRules_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiDataPropertyRulesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.akamai_property_rules.rules", "json",
					),
				),
			},
		},
	})
}
