package akamai

import (
	"fmt"
	gtmv1_3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"regexp"
	"testing"
)

var testAccAkamaiGTMResourceConfig = fmt.Sprintf(`
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
/*
resource "akamai_gtm_domain" "test_domain" {
        name = "${local.domain}"
        type = "weighted"
        //contract = "${local.contract}"
	contract = "${data.akamai_contract.contract.id}"
	comment =  "This is a test domain"
	//group     = "${local.group}"
	group  = "${data.akamai_group.group.id}"
	wait_on_complete = true
}
*/
resource "akamai_gtm_resource" "test_resource" {
    	//domain = "${akamai_gtm_domain.test_domain.name}"
	domain = "${local.domain}"
    	name = "test_resource_1"
    	aggregation_type = "latest"
    	type = "XML load object via HTTP"
    	load_imbalance_percentage = 50
    	wait_on_complete = true
	/*
    	depends_on = [
         	"akamai_gtm_domain.test_domain"
    	]
	*/
}
`, gtm_test_domain)

var testAccAkamaiGTMResourceUpdateConfig = fmt.Sprintf(`
provider "akamai" {
  gtm_section = "gtm"
} 

locals {
        domain = "%s"
}       
/*
data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_gtm_domain" "test_domain" {
        name = "${local.domain}"
        type = "weighted"
        contract = "${data.akamai_contract.contract.id}"
        comment =  "This is a test zone"
        group  = "${data.akamai_group.group.id}"
        wait_on_complete = true
}
*/
resource "akamai_gtm_resource" "test_resource" {
        domain = "${local.domain}" // "${akamai_gtm_domain.test_domain.name}"
        name = "test_resource_1"
        aggregation_type = "latest"
        type = "XML load object via HTTP"
        load_imbalance_percentage = 70
        wait_on_complete = true
        depends_on = [
                "akamai_gtm_domain.test_domain"
        ]
}
`, gtm_test_domain)

func TestAccAkamaiGTMResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAkamaiGTMDomainConfig,
				ExpectError: regexp.MustCompile(fmt.Sprintf("Domain Validation Error proposed domain name \"%s\" conflicts with existing domain", gtm_test_domain)),
			},
			{
				Config: testAccAkamaiGTMResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMResourceExists,
					resource.TestCheckResourceAttr("akamai_gtm_resource.test_resource", "load_imbalance_percentage", "50"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAkamaiGTMResource_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAkamaiGTMDomainConfig,
				ExpectError: regexp.MustCompile(fmt.Sprintf("Domain Validation Error proposed domain name \"%s\" conflicts with existing domain", gtm_test_domain)),
			},
			{
				Config: testAccAkamaiGTMResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMResourceExists,
					resource.TestCheckResourceAttr("akamai_gtm_resource.test_resource", "load_imbalance_percentage", "50"),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAkamaiGTMResourceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMResourceExists,
					resource.TestCheckResourceAttr("akamai_gtm_resource.test_resource", "load_imbalance_percentage", "70"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAkamaiGTMResourceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_resource" {
			continue
		}

		rname, dom, err := parseStringID(rs.Primary.ID)
		rsrc, err := gtmv1_3.GetResource(rname, dom)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] [Akamai GTMV1_3] Deleting test resource [%v]", rname)
		_, err = rsrc.Delete(dom)
		if err != nil {
			return fmt.Errorf("resource was not deleted %s. Error: %s", rs.Primary.ID, err.Error())
		}
	}
	return nil
}

func testAccCheckAkamaiGTMResourceExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_resource" {
			continue
		}
		rname, dom, err := parseStringID(rs.Primary.ID)
		_, err = gtmv1_3.GetResource(rname, dom)
		if err != nil {
			return err
		}
	}
	return nil
}
