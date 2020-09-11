package dns

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"strings"
	"testing"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

func TestCheckDNSv2Zone(t *testing.T) {
	tests := map[string]struct {
		init      func(*mocked)
		withError bool
	}{
		"type SECONDARY, validation ok": {
			init: func(m *mocked) {
				m.On("GetOk", "zone").Return("test", true).Once()
				m.On("GetOk", "type").Return("SECONDARY", true).Once()
				m.On("GetOk", "masters").Return(schema.NewSet(func(i interface{}) int {
					return 1
				}, []interface{}{"1"}), true).Once()
				m.On("GetOk", "target").Return("", true).Once()
				m.On("GetOk", "tsig_key").Return([]interface{}{"1"}, true).Once()
				m.On("GetOk", "sign_and_serve").Return(false, true).Once()
			},
			withError: false,
		},
		"type ALIAS, validation ok": {
			init: func(m *mocked) {
				m.On("GetOk", "zone").Return("test", true).Once()
				m.On("GetOk", "type").Return("ALIAS", true).Once()
				m.On("GetOk", "masters").Return(&schema.Set{}, true)
				m.On("GetOk", "target").Return("test-target", true).Once()
				m.On("GetOk", "tsig_key").Return([]interface{}{}, true).Once()
				m.On("GetOk", "sign_and_serve").Return(false, true).Once()
			},
			withError: false,
		},
		"different type, validation ok": {
			init: func(m *mocked) {
				m.On("GetOk", "zone").Return("test", true).Once()
				m.On("GetOk", "type").Return("SOME_TYPE", true).Once()
				m.On("GetOk", "masters").Return(&schema.Set{}, true)
				m.On("GetOk", "target").Return("", true).Once()
				m.On("GetOk", "tsig_key").Return([]interface{}{}, true).Once()
				m.On("GetOk", "sign_and_serve").Return(false, true).Once()
			},
			withError: false,
		},
		"type SECONDARY, masters is empty": {
			init: func(m *mocked) {
				m.On("GetOk", "zone").Return("test", true).Once()
				m.On("GetOk", "type").Return("SECONDARY", true).Once()
				m.On("GetOk", "masters").Return(schema.NewSet(func(i interface{}) int {
					return 1
				}, []interface{}{"1"}), true).Once()
				m.On("GetOk", "target").Return("test-target", true).Once()
				m.On("GetOk", "tsig_key").Return([]interface{}{"1"}, true).Once()
				m.On("GetOk", "sign_and_serve").Return(false, true).Once()
			},
			withError: true,
		},
		"type SECONDARY, target is not empty": {
			init: func(m *mocked) {
				m.On("GetOk", "zone").Return("test", true).Once()
				m.On("GetOk", "type").Return("SECONDARY", true).Once()
				m.On("GetOk", "masters").Return(&schema.Set{}, true).Once()
				m.On("GetOk", "target").Return("", true).Once()
				m.On("GetOk", "tsig_key").Return([]interface{}{"1"}, true).Once()
				m.On("GetOk", "sign_and_serve").Return(false, true).Once()
			},
			withError: true,
		},
		"type ALIAS, masters is not empty": {
			init: func(m *mocked) {
				m.On("GetOk", "zone").Return("test", true).Once()
				m.On("GetOk", "type").Return("ALIAS", true).Once()
				m.On("GetOk", "masters").Return(schema.NewSet(func(i interface{}) int {
					return 1
				}, []interface{}{"1"}), true).Once()
				m.On("GetOk", "target").Return("test-target", true).Once()
				m.On("GetOk", "tsig_key").Return([]interface{}{}, true).Once()
				m.On("GetOk", "sign_and_serve").Return(false, true).Once()
			},
			withError: true,
		},
		"type ALIAS, target is empty": {
			init: func(m *mocked) {
				m.On("GetOk", "zone").Return("test", true).Once()
				m.On("GetOk", "type").Return("ALIAS", true).Once()
				m.On("GetOk", "masters").Return(&schema.Set{}, true).Once()
				m.On("GetOk", "target").Return("", true).Once()
				m.On("GetOk", "tsig_key").Return([]interface{}{}, true).Once()
				m.On("GetOk", "sign_and_serve").Return(false, true).Once()
			},
			withError: true,
		},
		"type ALIAS, sign and serve is on": {
			init: func(m *mocked) {
				m.On("GetOk", "zone").Return("test", true).Once()
				m.On("GetOk", "type").Return("ALIAS", true).Once()
				m.On("GetOk", "masters").Return(&schema.Set{}, true).Once()
				m.On("GetOk", "target").Return("test-target", true).Once()
				m.On("GetOk", "tsig_key").Return([]interface{}{}, true).Once()
				m.On("GetOk", "sign_and_serve").Return(true, true).Once()
			},
			withError: true,
		},
		"type ALIAS, tsig is not empty": {
			init: func(m *mocked) {
				m.On("GetOk", "zone").Return("test", true).Once()
				m.On("GetOk", "type").Return("ALIAS", true).Once()
				m.On("GetOk", "masters").Return(&schema.Set{}, true).Once()
				m.On("GetOk", "target").Return("test-target", true).Once()
				m.On("GetOk", "tsig_key").Return([]interface{}{"1"}, true).Once()
				m.On("GetOk", "sign_and_serve").Return(false, true).Once()
			},
			withError: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mocked{}
			test.init(m)
			err := checkDNSv2Zone(m)
			m.AssertExpectations(t)
			if test.withError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

type mocked struct {
	mock.Mock
}

func (m *mocked) GetOk(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}
