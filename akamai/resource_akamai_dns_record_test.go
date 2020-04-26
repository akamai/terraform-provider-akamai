package akamai

import (
	"fmt"
	"log"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	//"strings"
	"testing"
)

var testAccAkamaiDNSv2RecordConfig = fmt.Sprintf(`
provider "akamai" {
  papi_section = "dns"
  dns_section = "dns"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_dns_zone" "test_zone" {
	contract = "${data.akamai_contract.contract.id}"
	zone = "exampleterraform.io"
	type = "primary"
	comment =  "This is a test zone"
	group     = "${data.akamai_group.group.id}"
	sign_and_serve = false
}

resource "akamai_dns_record" "a_record" {
	zone = "${akamai_dns_zone.test_zone.zone}"
	name = "exampleterraform.io"
	recordtype =  "A"
	active = true
	ttl = 300
	target = ["10.0.0.2","10.0.0.3"]
}
`)

var testAccAkamaiDNSv2RecordConfigWithCounter = fmt.Sprintf(`
provider "akamai" {
  papi_section = "dns"
  dns_section = "dns"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_dns_zone" "test_zone" {
	contract = "${data.akamai_contract.contract.id}"
	zone = "exampleterraform.io"
	type = "primary"
	comment =  "This is a test zone"
	group     = "${data.akamai_group.group.id}"
	sign_and_serve = false
}

resource "akamai_dns_record" "a_record" {
	count = 3
	zone = "${akamai_dns_zone.test_zone.zone}"
	name = "${count.index}.exampleterraform.io"
	recordtype =  "A"
	active = true
	ttl = 300
	target = ["10.0.0.2","10.0.0.3"]
}
`)

func TestAccAkamaiDNSv2Record_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiDNSv2RecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiDNSv2RecordConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiDNSv2RecordExists,
				),
			},
		},
	})
}

func TestAccAkamaiDNSv2Record_counter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiDNSv2RecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiDNSv2RecordConfigWithCounter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiDNSv2RecordExists,
				),
			},
		},
	})
}

func testAccCheckAkamaiDNSv2RecordDestroy(s *terraform.State) error {
	//conn := testAccProvider.Meta().(*Config)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_dns_record" {
			continue
		}

		log.Printf("[DEBUG] [Akamai DNSv2] Searching for zone [%v]", rs.Type)
		//request := &
		//hostname := strings.Split(rs.Primary.ID, "-")[2]
		//zone, err := dnsv2.GetZone(hostname)
		//if err != nil {
		//		return err
		//}
		//log.Printf("[DEBUG] [Akamai DNSv2] Searching for zone [%v]", zone)
		log.Printf("[DEBUG] [Akamai DNSv2] Searching for zone [%s]" + rs.Primary.Attributes["zone"])
		var zone string
		var host string
		var recordtype string

		zone = rs.Primary.Attributes["zone"]
		host = rs.Primary.Attributes["host"]
		recordtype = rs.Primary.Attributes["recordtype"]

		rdata, err := dnsv2.GetRdata(zone, host, recordtype)
		if err != nil {
			return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, err)
		}

		log.Printf("[DEBUG] [Akamai DNSv2] Searching for records [%v]", rdata)

	}
	return nil

}

func testAccCheckAkamaiDNSv2RecordExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_dns_record" {
			continue
		}

		var zone string
		var host string
		var recordtype string

		zone = rs.Primary.Attributes["zone"]
		host = rs.Primary.Attributes["host"]
		recordtype = rs.Primary.Attributes["recordtype"]

		rdata, err := dnsv2.GetRdata(zone, host, recordtype)
		if err != nil {
			return fmt.Errorf("error looking up "+recordtype+" records for %q: %s", host, err)
		}

		log.Printf("[DEBUG] [Akamai DNSv2] Searching for records [%v]", rdata)

	}
	return nil
}
