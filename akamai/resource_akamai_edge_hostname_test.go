package akamai

import (
	"fmt"
	"log"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/papi-v1"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	//"strings"
	"testing"
)

var testAccAkamaiSecureEdgeHostNameConfig = fmt.Sprintf(`
provider "akamai" {
  papi_section = "papi"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_edge_hostname" "test" {
    product = "prd_SPM"
    contract = "${data.akamai_contract.contract.id}"
    group = "${data.akamai_group.group.id}"
    edge_hostname =  "terraform-test.exampleterraform.io.edgesuite.net"
    ipv6 = true
}
`)

func TestAccAkamaiSecureEdgeHostName_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiSecureEdgeHostNameDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiSecureEdgeHostNameConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiSecureEdgeHostNameExists,
				),
			},
		},
	})
}

func testAccCheckAkamaiSecureEdgeHostNameDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_edge_hostname" {
			continue
		}
		log.Printf("[DEBUG] [Akamai SecureEdgeHostName] Delete for edgehostname not supported  [%v]", rs.Primary.ID)
	}
	return nil
}

func testAccCheckAkamaiSecureEdgeHostNameExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_edge_hostname" {
			continue
		}
		log.Printf("[DEBUG] [Akamai SecureEdgeHostName] Searching for edgehostname [%v]", rs.Primary.ID)
		hostname := rs.Primary.Attributes["edge_hostname"]

		groups := papi.NewGroups()
		e := groups.GetGroups()
		if e != nil {
			return e
		}

		group, e := groups.FindGroup(rs.Primary.Attributes["group"])
		if e != nil {
			return e
		}

		log.Println("[DEBUG] Figuring out edgehostnames GROUP = ", group)

		contracts := papi.NewContracts()

		e = contracts.GetContracts()
		if e != nil {
			return e
		}

		contract, e := contracts.FindContract(rs.Primary.Attributes["contract"])
		if e != nil {
			return e
		}

		log.Println("[DEBUG] Figuring out edgehostnames CONTRACT = ", contract)

		property := papi.NewProperty(papi.NewProperties())
		property.Group = group
		property.Contract = contract

		log.Println("[DEBUG] Figuring out edgehostnames ", rs.Primary.ID)
		edgeHostnames := papi.NewEdgeHostnames()
		log.Println("[DEBUG] NewEdgeHostnames empty struct  ", edgeHostnames.ContractID)
		err := edgeHostnames.GetEdgeHostnames(property.Contract, property.Group, "")
		if err != nil {
			return err
		}
		log.Println("[DEBUG] Edgehostnames exist in contract ")

		log.Println("[DEBUG] Edgehostnames Default host ", edgeHostnames.EdgeHostnames.Items[0])
		defaultEdgeHostname := edgeHostnames.EdgeHostnames.Items[0]

		log.Println("[DEBUG] Default Edgehostnames  ", defaultEdgeHostname, hostname)
		for _, eHn := range edgeHostnames.EdgeHostnames.Items {
			log.Println("[DEBUG] Edgehostname SEARCH  ", eHn.EdgeHostnameDomain)
			if eHn.EdgeHostnameDomain == hostname {
				log.Println("[DEBUG] Edgehostname FOUND  ", eHn.EdgeHostnameID)
				return nil
			}
		}
		return fmt.Errorf("error looking up Edge Hostname for %s", hostname)
	}
	return nil
}
