package akamai

import (
	"fmt"
	"log"
	//"strings"
	"testing"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var gtm_test_domain = "gtm_terra_testdomain.akadns.net"

//var gtm_test_contract = ""
//var gtm_test_group = ""

var testAccAkamaiGTMDomainConfig = fmt.Sprintf(`
provider "akamai" {
  gtm_section = "gtm"
}

locals {
  	domain = "%s"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_gtm_domain" "test_domain" {
        name = "${local.domain}"
        type = "weighted"
	contract = "${data.akamai_contract.contract.id}"
	comment =  "This is a test zone"
	group     = "${data.akamai_group.group.id}"
}
`, gtm_test_domain)

func TestAccAkamaiGTMDomain_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMDomainConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMDomainExists,
				),
			},
		},
	})
}

func testAccCheckAkamaiGTMDomainDestroy(s *terraform.State) error {

	// The API doesn't currently support Domain Delete
	return nil

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_domain" {
			continue
		}

		//hostname := strings.Split(rs.Primary.ID, "-")[5]
		domain, err := gtm.GetDomain(rs.Primary.ID)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] [Akamai GTMV1_3] deleting domain [%v]", domain)
		_, err = domain.Delete()
		if err != nil {
			if _, ok := err.(gtm.CommonError); ok {
				if err.(gtm.CommonError).GetItem("apiErrorMessage") == "DELETE method is not supported for this resource. Supported methods are: [GET, PUT]" {
					// can't delete Domain unless admin. Might be hole in this logic, but ....
					return nil
				}
			}
			return fmt.Errorf("domain %s was not deleted. Error:  %s", rs.Primary.ID, err.Error())
		}
	}
	return nil
}

func testAccCheckAkamaiGTMDomainExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_domain" {
			continue
		}
		_, err := gtm.GetDomain(rs.Primary.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
