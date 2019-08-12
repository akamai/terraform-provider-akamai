package akamai

import (
	"github.com/hashicorp/terraform/helper/schema"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAkamaiPropertyConfig = `
provider "akamai" {
  papi_section = "papi"
}

resource "akamai_property" "property" {
  name = "terraform-test1"

  contact = ["user@exampleterraform.io"]

  product = "prd_SPM"
  cp_code = "${akamai_cp_code.cp_code.id}"
  contract = "${data.akamai_contract.contract.id}"
  group = "${data.akamai_group.group.id}"

  hostnames = {
	"terraform.example.org" = "${akamai_edge_hostname.test.edge_hostname}"
  }
  
  rule_format = "v2016-11-15"
  
  rules = "${akamai_property_rules.rules.json}"
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
    edge_hostname =  "terraform-test1.exampleterraform.io.edgesuite.net"
	ipv4 = true
	ipv6 = true
}

data "akamai_property_rules" "rules" {
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
     			value = "exampleterraform.io"
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
`

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

func TestAccAkamaiProperty_isSecureTrue(t *testing.T) {
	config := map[string]interface{}{
		"rules":     `{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\"}}`,
		"is_secure": true,
	}

	d := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).ResourcesMap["akamai_property"].Schema, config)
	rules := &papi.Rules{}
	unmarshalRulesFromJSON(d, rules)
	if rules.Rule.Options.IsSecure != true {
		t.Fail()
	}
}

func TestAccAkamaiProperty_isSecureFalse(t *testing.T) {
	config := map[string]interface{}{
		"rules":     `{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\"}}`,
		"is_secure": false,
	}

	d := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).ResourcesMap["akamai_property"].Schema, config)
	rules := &papi.Rules{}
	unmarshalRulesFromJSON(d, rules)
	if rules.Rule.Options.IsSecure == true {
		t.Fail()
	}
}

func TestAccAkamaiProperty_isSecureUnset(t *testing.T) {
	config := map[string]interface{}{
		"rules": `{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\"}}`,
	}

	d := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).ResourcesMap["akamai_property"].Schema, config)
	rules := &papi.Rules{}
	unmarshalRulesFromJSON(d, rules)
	if rules.Rule.Options.IsSecure == true {
		t.Fail()
	}
}

func TestAccAkamaiProperty_isSecureJsonTrue(t *testing.T) {
	config := map[string]interface{}{
		"rules": "{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\",\"options\":{\"is_secure\":true}}}",
	}

	d := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).ResourcesMap["akamai_property"].Schema, config)
	rules := &papi.Rules{}
	unmarshalRulesFromJSON(d, rules)
	if rules.Rule.Options.IsSecure != true {
		t.Fail()
	}
}

func TestAccAkamaiProperty_isSecureJsonFalse(t *testing.T) {
	config := map[string]interface{}{
		"rules": "{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\",\"options\":{\"is_secure\":false}}}",
	}

	d := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).ResourcesMap["akamai_property"].Schema, config)
	rules := &papi.Rules{}
	unmarshalRulesFromJSON(d, rules)
	if rules.Rule.Options.IsSecure == true {
		t.Fail()
	}
}

func TestAccAkamaiProperty_isSecureJsonTrueConfigFalse(t *testing.T) {
	config := map[string]interface{}{
		"rules":     "{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\",\"options\":{\"is_secure\":true}}}",
		"is_secure": false,
	}

	d := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).ResourcesMap["akamai_property"].Schema, config)
	rules := &papi.Rules{}
	unmarshalRulesFromJSON(d, rules)
	if rules.Rule.Options.IsSecure == true {
		t.Fail()
	}
}

func TestAccAkamaiProperty_isSecureJsonFalseConfigTrue(t *testing.T) {
	config := map[string]interface{}{
		"rules":     "{\"accountId\":\"\",\"contractId\":\"\",\"groupId\":\"\",\"propertyId\":\"\",\"propertyVersion\":0,\"etag\":\"\",\"ruleFormat\":\"\",\"rules\":{\"name\":\"default\",\"options\":{\"is_secure\":false}}}",
		"is_secure": true,
	}

	d := schema.TestResourceDataRaw(t, Provider().(*schema.Provider).ResourcesMap["akamai_property"].Schema, config)
	rules := &papi.Rules{}
	unmarshalRulesFromJSON(d, rules)
	if rules.Rule.Options.IsSecure != true {
		t.Fail()
	}
}

func testAccCheckAkamaiPropertyDestroy(s *terraform.State) error {
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
