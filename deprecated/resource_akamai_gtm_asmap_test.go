package deprecated

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccAkamaiGTMAsMapConfig = fmt.Sprintf(`
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

resource "akamai_gtm_datacenter" "test_as_datacenter" {
    domain = akamai_gtm_domain.test_domain.name
    nickname = "test_as_datacenter"
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

resource "akamai_gtm_asmap" "test_as" {
    domain = akamai_gtm_domain.test_domain.name
    name = "test_asmap"
    default_datacenter {
        datacenter_id = data.akamai_gtm_default_datacenter.default_datacenter.datacenter_id
        nickname = data.akamai_gtm_default_datacenter.default_datacenter.nickname
    }
    assignment {
        datacenter_id = akamai_gtm_datacenter.test_as_datacenter.datacenter_id
        nickname = akamai_gtm_datacenter.test_as_datacenter.nickname
        as_numbers = [17334]
    }
    wait_on_complete = false
    depends_on = [
        akamai_gtm_domain.test_domain,
        akamai_gtm_datacenter.test_as_datacenter
    ]
}`, gtm_test_domain)

var testAccAkamaiGTMAsMapUpdateConfig = fmt.Sprintf(`
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

resource "akamai_gtm_datacenter" "test_as_datacenter" {
    domain = akamai_gtm_domain.test_domain.name
    nickname = "test_as_datacenter"
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

resource "akamai_gtm_asmap" "test_as" {
    domain = akamai_gtm_domain.test_domain.name
    name = "test_asmap"
    default_datacenter {
        datacenter_id = data.akamai_gtm_default_datacenter.default_datacenter.datacenter_id
        nickname = data.akamai_gtm_default_datacenter.default_datacenter.nickname
    }
    assignment {
        datacenter_id = akamai_gtm_datacenter.test_as_datacenter.datacenter_id
        nickname = akamai_gtm_datacenter.test_as_datacenter.nickname
        as_numbers = [17334]
    }
    wait_on_complete = false
    depends_on = [
        akamai_gtm_domain.test_domain,
        akamai_gtm_datacenter.test_as_datacenter
    ]
 
}`, gtm_test_domain)

var asMap *gtm.AsMap

func TestAccAkamaiGTMAsMap_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMAsMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMAsMapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMAsMapExists,
					testAccCheckNumbersValues,
					resource.TestCheckResourceAttr("akamai_gtm_asmap.test_as", "wait_on_complete", "false"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAkamaiGTMAsMap_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckAS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiGTMAsMapDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiGTMAsMapConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMAsMapExists,
					testAccCheckNumbersValues,
					resource.TestCheckResourceAttr("akamai_gtm_asmap.test_as", "wait_on_complete", "false"),
				),
				//ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccAkamaiGTMAsMapUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiGTMAsMapExists,
					testAccCheckNumbersValues,
					resource.TestCheckResourceAttr("akamai_gtm_asmap.test_as", "wait_on_complete", "false"),
				),
				//ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccPreCheckAS(t *testing.T) {

	testAccPreCheckTF(t)
	testCheckDeleteAsMap("test_asmap", gtm_test_domain)
	testAccDeleteDatacenterByNickname("test_as_datacenter", gtm_test_domain)

}

func testCheckDeleteAsMap(asName string, dom string) error {

	as, err := gtm.GetAsMap(asName, dom)
	if as == nil {
		return nil
	}
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] [Akamai GTMv1] Deleting test asmap [%v]", asName)
	_, err = as.Delete(dom)
	if err != nil {
		return fmt.Errorf("asmap was not deleted %s. Error: %s", asName, err.Error())
	}
	return nil

}

func testAccCheckAkamaiGTMAsMapDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_asmap" {
			continue
		}

		asName, dom, _ := parseStringID(rs.Primary.ID)
		if err := testCheckDeleteAsMap(asName, dom); err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckNumbersValues(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_asmap" {
			continue
		}
		if asMap == nil {
			return fmt.Errorf("asmap was not found for as_Numbers check")
		}
		log.Printf("[DEBUG] [Akamai GTMV1_3] ASMAP Validating as_numbers")
		// Walk thru all attributes
		mapAttribs := rs.Primary.Attributes
		assignEntries, err := strconv.Atoi(mapAttribs["assignment.#"])
		if err != nil {
			return fmt.Errorf("Assignments attribute was not found")
		}
		// Construct a list to compare
		assignMap := make(map[int][]int)
		for i := 0; i < assignEntries; i++ {
			iString := strconv.Itoa(i)
			assignBaseIndex := "assignments." + iString + "."
			dcid, _ := strconv.Atoi(mapAttribs[assignBaseIndex+"datacenter_id"])
			numbersEntries, _ := strconv.Atoi(mapAttribs[assignBaseIndex+"as_numbers.#"])
			numbersMap := []int{}
			numbersBaseIndex := assignBaseIndex + "as_numbers."
			for j := 0; j < numbersEntries; j++ {
				jString := strconv.Itoa(j)
				numEntry, _ := strconv.Atoi(mapAttribs[numbersBaseIndex+jString])
				numbersMap = append(numbersMap, numEntry)
			}
			assignMap[dcid] = numbersMap
		}
		for id, entry := range assignMap {
			for _, rAssignment := range asMap.Assignments {
				if id != rAssignment.DatacenterId {
					continue
				}
				compares := 0
				for _, n := range entry {
					for _, rasn := range rAssignment.AsNumbers {
						if rasn == int64(n) {
							compares += 1
							continue
						}
					}
				}
				if compares != len(entry) {
					return fmt.Errorf("Assignments numbers mismatch")
				}
				log.Printf("[DEBUG] [Akamai GTMV1_3] ASMAP assignment numbers DC match [%v]", id)
			}
		}
		return nil // only one
	}
	return fmt.Errorf("AsMap resource not found in state")

}

func testAccCheckAkamaiGTMAsMapExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_gtm_asmap" {
			continue
		}

		asName, dom, err := parseStringID(rs.Primary.ID)
		asMap, err = gtm.GetAsMap(asName, dom)
		if err != nil {
			return err
		}
	}
	return nil
}
