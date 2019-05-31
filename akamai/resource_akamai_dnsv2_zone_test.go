package akamai

import (
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"strings"
	"testing"
)

var testAccAkamaiDNSv2ZoneConfig = fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  dns_section = "dns"
}

resource "akamai_dns_zone" "test_zone" {
    contract = "C-1FRYVV3"
    zone = "akavaiodeveloper.net"
    masters = ["1.2.3.4" , "1.2.3.5"]
    type = "PRIMARY"
    comment =  "This is a test zone"
    group     = "64867"
    signandserve = false
}

`)

var testAccAkamaiDNSv2ZoneConfigWithCounter = fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  dns_section = "dns"
}

locals {
  zone = "akavaiodeveloper.net"
}

resource "akamai_dns_zone" "test_zone" {
    contract = "C-1FRYVV3"
    zone = "${local.zone}"
    masters = ["1.2.3.4" , "1.2.3.5"]
    type = "PRIMARY"
    comment =  "This is a test zone"
    group     = "64867"
    signandserve = false
}

`)

func TestAccAkamaiDNSv2Zone_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiDNSv2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiDNSv2ZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiDNSv2ZoneExists,
				),
			},
		},
	})
}

func TestAccAkamaiDNSv2Zone_counter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiDNSv2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiDNSv2ZoneConfigWithCounter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiDNSv2ZoneExists,
				),
			},
		},
	})
}

func testAccCheckAkamaiDNSv2ZoneDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_dns_zone" {
			continue
		}

		hostname := strings.Split(rs.Primary.ID, "-")[2]
		zone, err := dnsv2.GetZone(hostname)
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] [Akamai DNSv2] Searching for zone [%v]", zone)
		/*
			if len(zone.Zone.A) > 0 ||
				len(zone.Zone.Aaaa) > 0 ||
				len(zone.Zone.Afsdb) > 0 ||
				len(zone.Zone.Cname) > 0 ||
				len(zone.Zone.Dnskey) > 0 ||
				len(zone.Zone.Ds) > 0 ||
				len(zone.Zone.Hinfo) > 0 ||
				len(zone.Zone.Loc) > 0 ||
				len(zone.Zone.Mx) > 0 ||
				len(zone.Zone.Naptr) > 0 ||
				len(zone.Zone.Nsec3) > 0 ||
				len(zone.Zone.Nsec3param) > 0 ||
				len(zone.Zone.Ptr) > 0 ||
				len(zone.Zone.Rp) > 0 ||
				len(zone.Zone.Rrsig) > 0 ||
				len(zone.Zone.Spf) > 0 ||
				len(zone.Zone.Srv) > 0 ||
				len(zone.Zone.Sshfp) > 0 ||
				len(zone.Zone.Txt) > 0 {
				// These never get deleted
				// len(zone.Zone.Ns) > 0 ||
				// len(zone.Zone.Soa) > 0 ||
				return fmt.Errorf("zone was not deleted %s", hostname)
			}*/
	}
	return nil
}

func testAccCheckAkamaiDNSv2ZoneExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_dns_zone" {
			continue
		}

		hostname := strings.Split(rs.Primary.ID, "-")[2]
		_, err := dnsv2.GetZone(hostname)
		if err != nil {
			return err
		}
	}
	return nil
}
