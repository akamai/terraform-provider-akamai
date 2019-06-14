package akamai

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
	json := `{"accountId":"","contractId":"","groupId":"","propertyId":"","propertyVersion":0,"etag":"","ruleFormat":"","rules":{"name":"default","behaviors":[{"name":"caching","options":{"behavior":"MAX_AGE","mustRevalidate":false,"ttl":"1d"}},{"name":"cpCode","options":{"id":"cp-code-id"}},{"name":"origin","options":{"cacheKeyHostname":"ORIGIN_HOSTNAME","compress":1,"enableTrueClientIp":0,"forwardHostHeader":"REQUEST_HOST_HEADER","hostname":"exampleterraform.io","httpPort":80,"httpsPort":443,"originCertificate":"","originSni":1,"originType":"CUSTOMER","ports":"","verificationMode":"PLATFORM_SETTINGS"}}],"options":{}}}`
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiPropertyVariablesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiPropertyVariablesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"akamai_property_variables.test", "json", json,
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
