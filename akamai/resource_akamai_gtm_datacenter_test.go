package akamai

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
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
        name = local.domain
        type = "weighted"
	contract = data.akamai_contract.contract.id
	comment =  "This is a test domain"
	group  = data.akamai_group.group.id
	load_imbalance_percentage = 10
	wait_on_complete = false
}

resource "akamai_gtm_datacenter" "test_datacenter" {
    domain = akamai_gtm_domain.test_domain.name
    nickname = "test_datacenter"
    continent = "EU"
    wait_on_complete = false
    default_load_object {
        load_object = "test"
        load_object_port = 80
        load_servers = ["1.2.3.4", "1.2.3.5"]
    }
    depends_on = [
         akamai_gtm_domain.test_domain
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
        name = local.domain
        type = "weighted"
        contract = data.akamai_contract.contract.id
        comment =  "This is a test domain"
        group  = "${data.akamai_group.group.id}"
        load_imbalance_percentage = 10
        wait_on_complete = false
}

resource "akamai_gtm_datacenter" "test_datacenter" {
    domain = akamai_gtm_domain.test_domain.name
    nickname = "test_datacenter"
    continent = "NA"
    wait_on_complete = false
    default_load_object {
        load_object = "test"
        load_object_port = 80
        load_servers = ["1.2.3.4", "1.2.3.5"]
    }  
    depends_on = [
         akamai_gtm_domain.test_domain
    ]    
}   
`, gtm_test_domain)

func TestAccAkamaiGTMDatacenter_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDC(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMDatacenterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMDatacenterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMDatacenterExists,
					resource.TestCheckResourceAttr("akamai_gtm_datacenter.test_datacenter", "continent", "EU"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAkamaiGTMDatacenter_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDC(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMDatacenterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMDatacenterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMDatacenterExists,
					resource.TestCheckResourceAttr("akamai_gtm_datacenter.test_datacenter", "continent", "EU"),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAkamaiGTMDatacenterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMDatacenterExists,
					resource.TestCheckResourceAttr("akamai_gtm_datacenter.test_datacenter", "continent", "NA"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccPreCheckDC(t *testing.T) {

	testAccPreCheckTF(t)
	testAccDeleteDatacenterByNickname("test_datacenter", gtm_test_domain)

}

func testAccCheckAkamaiGTMDatacenterDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_datacenter" {
			continue
		}

		dcid, dom, _ := parseIntID(rs.Primary.ID)
		if err := testAccDeleteDatacenter(dcid, dom); err != nil {
			return err
		}
	}
	return nil
}

func testAccDeleteDatacenterByNickname(nickname string, dom string) error {

	dcList, err := gtm.ListDatacenters(dom)
	if dcList == nil || err != nil {
		return err
	}
	for _, dc := range dcList {
		if dc.Nickname == nickname {
			_, err := dc.Delete(dom)
			return err
		}
	}
	return nil

}

func testAccDeleteDatacenter(dcid int, dom string) error {

	dc, err := gtm.GetDatacenter(dcid, dom)
	if dc == nil {
		return nil
	}
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Deleting test datacenter [%v]", dcid)
	_, err = dc.Delete(dom)
	if err != nil {
		return fmt.Errorf("datacenter was not deleted %d. Error: %s", dcid, err.Error())
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
		_, err = gtm.GetDatacenter(dcid, dom)
		if err != nil {
			return err
		}
	}
	return nil
}
