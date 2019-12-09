package akamai

import (
	"fmt"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
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

resource "akamai_gtm_domain" "test_domain" {
        name = "${local.domain}"
        type = "weighted"
	contract = "${data.akamai_contract.contract.id}"
	comment =  "This is a test domain"
	group  = "${data.akamai_group.group.id}"
        load_imbalance_percentage = 10
	wait_on_complete = false
}

resource "akamai_gtm_resource" "test_resource" {
    	//domain = "${akamai_gtm_domain.test_domain.name}"
	domain = "${local.domain}"
    	name = "test_resource_1"
    	aggregation_type = "latest"
    	type = "XML load object via HTTP"
    	load_imbalance_percentage = 50
    	wait_on_complete = false
    	depends_on = [
         	"akamai_gtm_domain.test_domain"
    	]
}
`, gtm_test_domain)

var testAccAkamaiGTMResourceUpdateConfig = fmt.Sprintf(`
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
        group  = "${data.akamai_group.group.id}"
        load_imbalance_percentage = 10
        wait_on_complete = false
}

resource "akamai_gtm_resource" "test_resource" {
        domain = "${local.domain}" // "${akamai_gtm_domain.test_domain.name}"
        name = "test_resource_1"
        aggregation_type = "latest"
        type = "XML load object via HTTP"
        load_imbalance_percentage = 70
        wait_on_complete = false
        depends_on = [
                "akamai_gtm_domain.test_domain"
        ]
}
`, gtm_test_domain)

func TestAccAkamaiGTMResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTF(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMResourceDestroy,
		Steps: []resource.TestStep{
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
		PreCheck:     func() { testAccPreCheckTF(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMResourceDestroy,
		Steps: []resource.TestStep{
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

		rname, dom, _ := parseStringID(rs.Primary.ID)
		rsrc, err := gtm.GetResource(rname, dom)
		if rsrc == nil {
			log.Printf("[DEBUG] [Akamai GTM] : Resource %s not found. Ignoring error (Hashicorp).", rname)
			return nil
		}
		if err != nil {
			log.Printf("[INFO] [Akamai GTM] Resource Destroy: Error reading resource [%s]", err.Error())
			return err
		}
		log.Printf("[DEBUG] [Akamai GTMv1] Deleting test resource [%v]", rname)
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
		_, err = gtm.GetResource(rname, dom)
		if err != nil {
			return err
		}
	}
	return nil
}
