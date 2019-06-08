package akamai

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"testing"
)

var testAccAkamaiPropertyRulesConfig = fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  papi_section = "global"
}

resource "akamai_property_rules" "test" {
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

func TestAccAkamaiPropertyRules_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiPropertyRulesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiPropertyRulesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"akamai_property_rules.test", "rules.#", "1",
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