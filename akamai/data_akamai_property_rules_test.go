package akamai

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
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

var testAccAkamaiDataPropertyRulesSiteshield = `
provider "akamai" {
	property {
		host = "test"
		access_token = "test"
		client_token = "test"
		client_secret = "test"
	}
}

output "json" {
	value = "${data.akamai_property_rules.rules.json}"
}

data "akamai_property_rules" "rules" {
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

func TestAkamaiDataPropertyRules_siteshield(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Destroy: false,
				Config:  testAccAkamaiDataPropertyRulesSiteshield,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "json", "{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\",\"behaviors\":[{\"name\":\"siteShield\",\"options\":{\"ssmap\":{\"value\":\"mapname.akamai.net\"}}}],\"options\":{}}}"),
					resource.TestCheckResourceAttrSet(
						"data.akamai_property_rules.rules", "json",
					),
				),
			},
		},
	})
}

var testAccAkamaiDataPropertyRulesCPCode = `
provider "akamai" {
	property {
		host = "test"
		access_token = "test"
		client_token = "test"
		client_secret = "test"
	}
}

output "json" {
	value = "${data.akamai_property_rules.rules.json}"
}

data "akamai_property_rules" "rules" {
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

func TestAkamaiDataPropertyRules_cpCode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Destroy: false,
				Config:  testAccAkamaiDataPropertyRulesCPCode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "json", "{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\",\"behaviors\":[{\"name\":\"cpCode\",\"options\":{\"value\":{\"id\":12345}}}],\"options\":{}}}"),
					resource.TestCheckResourceAttrSet(
						"data.akamai_property_rules.rules", "json",
					),
				),
			},
		},
	})
}

var testAccAkamaiDataPropertyRulesIsSecure = `
provider "akamai" {
	property {
		host = "test"
		access_token = "test"
		client_token = "test"
		client_secret = "test"
	}
}

output "json" {
	value = "${data.akamai_property_rules.rules.json}"
}

data "akamai_property_rules" "rules" {
 	rules {
		is_secure = %s
	}
}
`

func TestAkamaiDataPropertyRules_isSecureTrue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Destroy: false,
				Config:  fmt.Sprintf(testAccAkamaiDataPropertyRulesIsSecure, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "json", "{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\",\"options\":{\"is_secure\":true}}}"),
					resource.TestCheckResourceAttrSet(
						"data.akamai_property_rules.rules", "json",
					),
				),
			},
		},
	})
}

func TestAkamaiDataPropertyRules_isSecureFalse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Destroy: false,
				Config:  fmt.Sprintf(testAccAkamaiDataPropertyRulesIsSecure, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.akamai_property_rules.rules", "json", "{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\",\"options\":{}}}"),
					resource.TestCheckResourceAttrSet(
						"data.akamai_property_rules.rules", "json",
					),
				),
			},
		},
	})
}
