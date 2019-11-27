package akamai

import (
	"fmt"
	"log"
	"testing"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_3"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAkamaiGTMGeoMapConfig = fmt.Sprintf(`
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
	wait_on_complete = true
}

resource "akamai_gtm_datacenter" "test_datacenter" {
    domain = "${akamai_gtm_domain.test_domain.name}"
    nickname = "test_geo_datacenter"
    wait_on_complete = true
    virtual = true
    default_load_object = [{
        load_object = "test"
        load_object_port = 80
        load_servers = ["1.2.3.4", "1.2.3.5"]
    }]
    depends_on = [
         "akamai_gtm_domain.test_domain"
    ]
}

resource "akamai_gtm_geomap" "test_geo" {
    domain = "${akamai_gtm_domain.test_domain.name}"
    name = "test_geomap"
    default_datacenter = [{
        datacenter_id = 5400
        nickname = "All Others"
        }]
    assignments = [{
        datacenter_id = "${akamai_gtm_datacenter.test_datacenter.datacenter_id}"
        nickname = "${akamai_gtm_datacenter.test_datacenter.nickname}"
       	countries = ["US"] 
        }]
    wait_on_complete = true
    depends_on = [
        "akamai_gtm_domain.test_domain",
        "akamai_gtm_datacenter.test_datacenter"
    ]
}`, gtm_test_domain)

var testAccAkamaiGTMGeoMapUpdateConfig = fmt.Sprintf(`
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
        wait_on_complete = true
}

resource "akamai_gtm_datacenter" "test_datacenter" {
    domain = "${akamai_gtm_domain.test_domain.name}"
    nickname = "test_geo_datacenter"
    wait_on_complete = true
    virtual = true
    default_load_object = [{
        load_object = "test"
        load_object_port = 80
        load_servers = ["1.2.3.4", "1.2.3.5"]
    }]  
    depends_on = [
         "akamai_gtm_domain.test_domain"
    ]    
}  

resource "akamai_gtm_geomap" "test_geo" {
    domain = "${akamai_gtm_domain.test_domain.name}"
    name = "test_geomap"
    default_datacenter = [{
        datacenter_id = 5400
        nickname = "All Others"
        }]
    assignments = [{
        datacenter_id = "${akamai_gtm_datacenter.test_datacenter.datacenter_id}"
        nickname = "${akamai_gtm_datacenter.test_datacenter.nickname}"
        countries = ["US"]
        }]
    wait_on_complete = false
    depends_on = [
        "akamai_gtm_domain.test_domain",
        "akamai_gtm_datacenter.test_datacenter"
    ]
 
}`, gtm_test_domain)

func TestAccAkamaiGTMGeoMap_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTF(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMGeoMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMGeoMapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMGeoMapExists,
					resource.TestCheckResourceAttr("akamai_gtm_geomap.test_geo", "wait_on_complete", "true"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAkamaiGTMGeoMap_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckTF(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMGeoMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMGeoMapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMGeoMapExists,
					resource.TestCheckResourceAttr("akamai_gtm_geomap.test_geo", "wait_on_complete", "true"),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAkamaiGTMGeoMapUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMGeoMapExists,
					resource.TestCheckResourceAttr("akamai_gtm_geomap.test_geo", "wait_on_complete", "false"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAkamaiGTMGeoMapDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_geomap" {
			continue
		}

		geoName, dom, err := parseStringID(rs.Primary.ID)
		geo, err := gtm.GetGeoMap(geoName, dom)
		if geo == nil {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] [Akamai GTMV1_3] Deleting test geomap [%v]", geoName)
		_, err = geo.Delete(dom)
		if err != nil {
			return fmt.Errorf("geomap was not deleted %s. Error: %s", rs.Primary.ID, err.Error())
		}
	}
	return nil
}

func testAccCheckAkamaiGTMGeoMapExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_geomap" {
			continue
		}

		geoName, dom, err := parseStringID(rs.Primary.ID)
		_, err = gtm.GetGeoMap(geoName, dom)
		if err != nil {
			return err
		}
	}
	return nil
}
