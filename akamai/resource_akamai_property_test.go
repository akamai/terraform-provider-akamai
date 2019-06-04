package akamai

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

var testAccAkamaiPropertyConfig = fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "global"
}

resource "akamai_property" "akamai_developer" {
  name = "akamaideveloper.com"

  contact = ["dshafik@akamai.com"]


  product = "prd_SPM"
  cp_code = "409449"
  contract = "ctr_C-1FRYVV3"
  group = "grp_68817"
  

  rule_format = "v2016-11-15"
  
  rules = "${akamai_property_rules.akamai_developer.json}"

  origin {
    is_secure = false
    hostname = "akamaideveloper.net"
    forward_hostname = "ORIGIN_HOSTNAME"
  }

  
}

resource "akamai_property_rules" "akamai_developer" {

  rules {
    behavior {
      name = "downstreamCache"
      option {
        key = "behavior"
        value = "TUNNEL_ORIGIN"
      }
    }

    rule {
      name = "Uncacheable Responses"
      comment = "Cache me outside"
      criteria {
        name = "cacheability"
        option {
          key = "matchOperator"
          value = "IS_NOT"
        }
        option {
          key = "value"
          value = "CACHEABLE"
        }
      }
      behavior {
        name = "downstreamCache"
        option {
          key = "behavior"
          value = "TUNNEL_ORIGIN"
        }
      }
      rule {
        name = "Uncacheable Responses"
        comment = "Child rule"
        criteria {
          name = "cacheability"
          option {
            key = "matchOperator"
            value = "IS_NOT"
          }
          option {
            key = "value"
            value = "CACHEABLE"
          }
        }
        behavior {
          name = "downstreamCache"
          option {
            key = "behavior"
            value = "TUNNEL_ORIGIN"
          }
        }
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
