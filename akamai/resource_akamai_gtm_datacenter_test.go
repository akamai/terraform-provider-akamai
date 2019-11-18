package akamai

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	gtmv1_3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAkamaiGTMDatacenterConfig = fmt.Sprintf(`
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
        //contract = "${local.contract}"
	contract = "${data.akamai_contract.contract.id}"
	comment =  "This is a test zone"
	//group     = "${local.group}"
	group  = "${data.akamai_group.group.id}"
	wait_on_complete = true
}

resource "akamai_gtm_datacenter" "test_datacenter" {
    domain = "${akamai_gtm_domain.test_domain.name}"
    nickname = "test_datacenter"
    wait_on_complete = true
    default_load_object = [{
        load_object = "test"
        load_object_port = 80
        load_servers = ["1.2.3.4", "1.2.3.5"]
    }]
    depends_on = [
         "akamai_gtm_domain.test_domain"
    ]
}
`, gtm_test_domain)

var testAccAkamaiGTMDatacenterUpdateConfig = fmt.Sprintf(`
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
        wait_on_complete = true
}

resource "akamai_gtm_datacenter" "test_datacenter" {
    domain = "${akamai_gtm_domain.test_domain.name}"
    nickname = "test_datacenter_updated"
    wait_on_complete = true
    default_load_object = [{
        load_object = "test"
        load_object_port = 80
        load_servers = ["1.2.3.4", "1.2.3.5"]
    }]  
    depends_on = [
         "akamai_gtm_domain.test_domain"
    ]    
}   
`, gtm_test_domain)

func TestAccAkamaiGTMDatacenter_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMDatacenterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMDatacenterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMDatacenterExists,
					resource.TestCheckResourceAttr("akamai_gtm_datacenter.test_datacenter", "nickname", "test_datacenter"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAkamaiGTMDatacenter_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMDatacenterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMDatacenterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMDatacenterExists,
					resource.TestCheckResourceAttr("akamai_gtm_datacenter.test_datacenter", "nickname", "test_datacenter"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAkamaiGTMDatacenterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMDatacenterExists,
					resource.TestCheckResourceAttr("akamai_gtm_datacenter.test_datacenter", "nickname", "test_datacenter_updated"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAkamaiGTMDatacenterDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_datacenter" {
			continue
		}

		dcid, dom, err := parseIntID(rs.Primary.ID)
		dc, err := gtmv1_3.GetDatacenter(dcid, dom)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] [Akamai GTMV1_3] Deleting test datacenter [%v]", dcid)
		_, err = dc.Delete(dom)
		if err != nil {
			return fmt.Errorf("datacenter was not deleted %s. Error: %s", rs.Primary.ID, err.Error())
		}
	}
	return nil
}

func parseIntID(id string) (int, string, error) {
	idComp := strings.Split(id, ":")
	if len(idComp) < 2 {
		return 0, "", errors.New("Invalid Datacenter ID")
	}
	dcid, err := strconv.Atoi(idComp[1])
	if err != nil {
		return 0, "", err
	}
	return dcid, idComp[0], nil

}

func testAccCheckAkamaiGTMDatacenterExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_datacenter" {
			continue
		}

		dcid, dom, err := parseIntID(rs.Primary.ID)
		_, err = gtmv1_3.GetDatacenter(dcid, dom)
		if err != nil {
			return err
		}
	}
	return nil
}
