package akamai

import (
	"fmt"
	"log"
	"testing"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAkamaiGTMCidrMapConfig = fmt.Sprintf(`
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

data "akamai_gtm_default_datacenter" "default_datacenter" {
    domain = akamai_gtm_domain.test_domain.name
    datacenter = 5400
}

resource "akamai_gtm_domain" "test_domain" {
    name = local.domain
    type = "weighted"
    contract = data.akamai_contract.contract.id
    comment =  "This is a test zone"
    group  = data.akamai_group.group.id
    load_imbalance_percentage = 10
    wait_on_complete = false
}

resource "akamai_gtm_datacenter" "test_cidr_datacenter" {
    domain = akamai_gtm_domain.test_domain.name
    nickname = "test_cidr_datacenter"
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

resource "akamai_gtm_cidrmap" "test_cidr" {
    domain = akamai_gtm_domain.test_domain.name
    name = "test_cidrmap"
    default_datacenter {
        datacenter_id = data.akamai_gtm_default_datacenter.default_datacenter.datacenter_id
        nickname = data.akamai_gtm_default_datacenter.default_datacenter.nickname
    }
    assignment {
        datacenter_id = akamai_gtm_datacenter.test_cidr_datacenter.datacenter_id
        nickname = akamai_gtm_datacenter.test_cidr_datacenter.nickname
        blocks = ["1.2.3.9/24"]
    }
    wait_on_complete = false
    depends_on = [
        akamai_gtm_domain.test_domain,
        akamai_gtm_datacenter.test_cidr_datacenter
    ]
}`, gtm_test_domain)

var testAccAkamaiGTMCidrMapUpdateConfig = fmt.Sprintf(`
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

data "akamai_gtm_default_datacenter" "default_datacenter" {
    domain = akamai_gtm_domain.test_domain.name
    datacenter = 5400
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

resource "akamai_gtm_datacenter" "test_cidr_datacenter" {
    domain = akamai_gtm_domain.test_domain.name
    nickname = "test_cidr_datacenter"
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

resource "akamai_gtm_cidrmap" "test_cidr" {
    domain = akamai_gtm_domain.test_domain.name
    name = "test_cidrmap"
    default_datacenter {
        datacenter_id = data.akamai_gtm_default_datacenter.default_datacenter.datacenter_id
        nickname = data.akamai_gtm_default_datacenter.default_datacenter.nickname
    }
    assignment {
        datacenter_id = akamai_gtm_datacenter.test_cidr_datacenter.datacenter_id
        nickname = akamai_gtm_datacenter.test_cidr_datacenter.nickname
        blocks = ["1.2.3.9/24"]
    }
    wait_on_complete = false
    depends_on = [
        akamai_gtm_domain.test_domain,
        akamai_gtm_datacenter.test_cidr_datacenter
    ]
 
}`, gtm_test_domain)

func TestAccAkamaiGTMCidrMap_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCidr(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMCidrMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMCidrMapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMCidrMapExists,
					resource.TestCheckResourceAttr("akamai_gtm_cidrmap.test_cidr", "wait_on_complete", "false"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAkamaiGTMCidrMap_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCidr(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMCidrMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMCidrMapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMCidrMapExists,
					resource.TestCheckResourceAttr("akamai_gtm_cidrmap.test_cidr", "wait_on_complete", "false"),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAkamaiGTMCidrMapUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMCidrMapExists,
					resource.TestCheckResourceAttr("akamai_gtm_cidrmap.test_cidr", "wait_on_complete", "false"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccPreCheckCidr(t *testing.T) {

	testAccPreCheckTF(t)
	testCheckDeleteCidrMap("test_cidrmap", gtm_test_domain)
	testAccDeleteDatacenterByNickname("test_cidr_datacenter", gtm_test_domain)

}

func testCheckDeleteCidrMap(cidrName string, dom string) error {

	cidr, err := gtm.GetCidrMap(cidrName, dom)
	if cidr == nil {
		return nil
	}
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Deleting test cidrmap [%v]", cidrName)
	_, err = cidr.Delete(dom)
	if err != nil {
		return fmt.Errorf("cidrmap was not deleted %s. Error: %s", cidrName, err.Error())
	}
	return nil

}

func testAccCheckAkamaiGTMCidrMapDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_cidrmap" {
			continue
		}

		cidrName, dom, _ := parseStringID(rs.Primary.ID)
		if err := testCheckDeleteCidrMap(cidrName, dom); err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckAkamaiGTMCidrMapExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_cidrmap" {
			continue
		}

		cidrName, dom, err := parseStringID(rs.Primary.ID)
		_, err = gtm.GetCidrMap(cidrName, dom)
		if err != nil {
			return err
		}
	}
	return nil
}
