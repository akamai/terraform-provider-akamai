package gtm

import (
	"fmt"
	"log"
	"testing"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var gtmTestDomain = "gtm_terra_testdomain.akadns.net"

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
        name = local.domain
        type = "weighted"
	contract = data.akamai_contract.contract.id
	comment =  "Test"
	group     = data.akamai_group.group.id
	load_imbalance_percentage = 10
}
`, gtmTestDomain)

var testAccAkamaiGTMDomainUpdateConfig = fmt.Sprintf(`
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
        name = local.domain
        type = "weighted"
        contract = data.akamai_contract.contract.id
        comment =  "Test update"
        group     = data.akamai_group.group.id
        load_imbalance_percentage = 10
}
`, gtmTestDomain)

func TestAccAkamaiGTMADomain_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTF(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMDomainConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMDomainExists,
					resource.TestCheckResourceAttr("akamai_gtm_domain.test_domain", "comment", "Test"),
				),
			},
		},
	})
}

func TestAccAkamaiGTMADomain_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTF(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMDomainConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMDomainExists,
					resource.TestCheckResourceAttr("akamai_gtm_domain.test_domain", "comment", "Test"),
				),
			},
			{
				Config: testAccAkamaiGTMDomainUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMDomainExists,
					resource.TestCheckResourceAttr("akamai_gtm_domain.test_domain", "comment", "Test update"),
				),
			},
		},
	})
}

func testAccCheckAkamaiGTMDomainDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_domain" {
			continue
		}

		domain, err := gtm.GetDomain(rs.Primary.ID)
		if domain == nil {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] [Akamai GTMV1] deleting domain [%v]", domain)
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

// Sets a Hack flag so cn work with existing Domains (only Admin can Delete)
func testAccPreCheckTF(_ *testing.T) {

	// by definition, we are running acceptance tests. ;-)
	log.Printf("[DEBUG] [Akamai GTMV1] Setting HashiAcc true")
	HashiAcc = true

}
