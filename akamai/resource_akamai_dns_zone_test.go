package akamai

import (
	"fmt"
	"log"
	"strings"
	"testing"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccAkamaiDNSPrimaryZoneConfig = fmt.Sprintf(`
provider "akamai" {
  dns_section = "dns"
}

locals {
  zone = "akavdev.net"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_dns_zone" "primary_test_zone" {
	contract = "${data.akamai_contract.contract.id}"
	zone = "primaryexampleterraform.io"
	type = "primary"
	comment =  "This is a test primary zone"
	group     = "${data.akamai_group.group.id}"
	sign_and_serve = false
}
`)

var testAccAkamaiDNSSecondaryZoneConfig = fmt.Sprintf(`
provider "akamai" {
  dns_section = "dns"
}
locals {
  zone = "akavdev.net"
}
data "akamai_contract" "contract" {
}
data "akamai_group" "group" {
}
resource "akamai_dns_zone" "test_secondary_zone" {
	contract = "${data.akamai_contract.contract.id}"
	zone = "secondaryexampleterraform.io"
	masters = ["1.2.3.4" , "1.2.3.5"]
	type = "secondary"
	comment =  "This is a secondary test zone"
	group     = "${data.akamai_group.group.id}"
	sign_and_serve = false
}
`)

var testAccAkamaiDNSPrimaryZoneConfigWithCounter = fmt.Sprintf(`
provider "akamai" {
  papi_section = "dns"
  dns_section = "dns"
}

data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}

resource "akamai_dns_zone" "primary_test_zone" {
	contract = "${data.akamai_contract.contract.id}"
	zone = "primaryexampleterraform.io"
	type = "primary"
	comment =  "This is a test primary zone"
	group     = "${data.akamai_group.group.id}"
	sign_and_serve = false
}
`)

func TestAccAkamaiDNSPrimaryZone_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiDNSPrimaryZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiDNSZoneExists,
				),
			},
		},
	})
}

func TestAccAkamaiDNSSecondaryZone_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiDNSSecondaryZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiDNSZoneExists,
				),
			},
		},
	})
}

func TestAccAkamaiDNSPrimaryZone_counter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiDNSZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAkamaiDNSPrimaryZoneConfigWithCounter,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAkamaiDNSZoneExists,
				),
			},
		},
	})
}

func testAccCheckAkamaiDNSZoneDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_dns_zone" {
			continue
		}

		hostname := strings.Split(rs.Primary.ID, "#")[1]
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

func testAccCheckAkamaiDNSZoneExists(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "akamai_dns_zone" {
			continue
		}

		hostname := strings.Split(rs.Primary.ID, "#")[1]
		_, err := dnsv2.GetZone(hostname)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestValidateZoneType(t *testing.T) {
	badValues := []string{"foo", "BAR"}
	goodValues := []string{"PRIMARY", "SECONDARY", "ALIAS", "primary", "secondary", "alias"}

	for _, bv := range badValues {
		_, err := validateZoneType(bv, "")
		if err == nil {
			t.Errorf("Value %v is invalid: %v", bv, err)
		}
	}

	for _, gv := range goodValues {
		_, err := validateZoneType(gv, "")
		if err != nil {
			t.Errorf("Value %v is invalid: %v", gv, err)
		}
	}
}
