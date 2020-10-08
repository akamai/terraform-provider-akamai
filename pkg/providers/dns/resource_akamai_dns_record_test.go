package dns

import (
	"context"
	"fmt"
	"log"

	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

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
		//hostname := strings.Split(rs.Primary.ID, "#")[2]
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

func TestDiffQuotedDNSRecord(t *testing.T) {
	tests := map[string]struct {
		givenOldTargetList []string
		givenNewTargetList []string
		givenOld           string
		givenNew           string
		givenRecordType    string
		expected           bool
	}{
		"target lists have different len": {
			givenOldTargetList: []string{"192.168.0.0", "192.168.0.0"},
			givenNewTargetList: []string{"192.168.0.0"},
			expected:           false,
		},
		"old and new values not empty, baseVal does not match any entry from list": {
			givenOldTargetList: []string{"record A", "record B"},
			givenNewTargetList: []string{"record C", "record D"},
			givenOld:           "old \\\" val",
			givenNew:           "value",
			expected:           false,
		},
		"old and new values not empty, baseVal matches (no quotations)": {
			givenOldTargetList: []string{"record A", "record B"},
			givenNewTargetList: []string{"this is record C", "record D"},
			givenOld:           "this is record C",
			givenNew:           "value",
			expected:           true,
		},
		"old and new values not empty, baseVal matches entry from list": {
			givenOldTargetList: []string{"record A", "record B"},
			givenNewTargetList: []string{"\"this \"is\" record C\"", "record D"},
			givenOld:           "\\\"this \\\"is\\\" record C\\\"",
			givenNew:           "value",
			expected:           true,
		},
		"old and new values not empty, baseVal does not match entry from list (invalid quotations on given)": {
			givenOldTargetList: []string{"record A", "record B"},
			givenNewTargetList: []string{"\"this \"is\" record C\"", "record D"},
			givenOld:           "\"this is record C\"",
			givenNew:           "value",
			expected:           false,
		},
		"old and new values not empty, baseVal does not match entry from list (invalid quotations in target)": {
			givenOldTargetList: []string{"record A", "record B"},
			givenNewTargetList: []string{"\\\"this is record C\\\"", "record D"},
			givenOld:           "\"this is record C\"",
			givenNew:           "value",
			expected:           false,
		},
		"old value is empty, baseVal matches (no quotations)": {
			givenOldTargetList: []string{"this is record A", "record B"},
			givenNewTargetList: []string{"record C", "record D"},
			givenOld:           "",
			givenNew:           "this is record A",
			expected:           true,
		},
		"old value is empty, baseVal matches (given value in single quotes)": {
			givenOldTargetList: []string{"this is record A", "record B"},
			givenNewTargetList: []string{"record C", "record D"},
			givenOld:           "",
			givenNew:           "\"this is record A\"",
			expected:           true,
		},
		"old value is empty, baseVal matches (target contains escaped quotation)": {
			givenOldTargetList: []string{"\"this \\\"is\\\" record A\"", "record B"},
			givenNewTargetList: []string{"record C", "record D"},
			givenOld:           "",
			givenNew:           "\"this \"is\" record A\"",
			expected:           true,
		},
		"old value is empty, baseVal does not match (invalid quotations)": {
			givenOldTargetList: []string{"\"this \"is\" record A\"", "record B"},
			givenNewTargetList: []string{"record C", "record D"},
			givenOld:           "",
			givenNew:           "\"this is record A\"",
			expected:           false,
		},
		"new value is empty, baseVal does not match any entry from list": {
			givenOldTargetList: []string{"record A", "record B"},
			givenNewTargetList: []string{"record C", "record D"},
			givenOld:           "old \\\" val",
			givenNew:           "",
			expected:           false,
		},
		"new value is empty, baseVal matches (no quotations)": {
			givenOldTargetList: []string{"record A", "record B"},
			givenNewTargetList: []string{"this is record C", "record D"},
			givenOld:           "this is record C",
			givenNew:           "",
			expected:           true,
		},
		"new value is empty, baseVal matches entry from list": {
			givenOldTargetList: []string{"record A", "record B"},
			givenNewTargetList: []string{"\"this \"is\" record C\"", "record D"},
			givenOld:           "\\\"this \\\"is\\\" record C\\\"",
			givenNew:           "",
			expected:           true,
		},
		"new value is empty, baseVal does not match entry from list (invalid quotations on given)": {
			givenOldTargetList: []string{"record A", "record B"},
			givenNewTargetList: []string{"\"this \"is\" record C\"", "record D"},
			givenOld:           "\"this is record C\"",
			givenNew:           "",
			expected:           false,
		},
		"new value is empty, baseVal does not match entry from list (invalid quotations in target)": {
			givenOldTargetList: []string{"record A", "record B"},
			givenNewTargetList: []string{"\\\"this is record C\\\"", "record D"},
			givenOld:           "\"this is record C\"",
			givenNew:           "",
			expected:           false,
		},
		"old and new values not empty, record type is AAAA, IPs match (both are short IPv6)": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8::2:2", "2001:db8::2:3"},
			givenOld:           `\"2001:db8::2:2\""`,
			givenNew:           "value",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"old and new values not empty, record type is AAAA, IPs match (target is full IPv6)": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8:0:0:0:0:2:2", "2001:db8:0:0:0:0:2:3"},
			givenOld:           `\"2001:db8::2:2\""`,
			givenNew:           "value",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"old and new values not empty, record type is AAAA, IPs match (both are full IPv6)": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8:0:0:0:0:2:2", "2001:db8:0:0:0:0:2:3"},
			givenOld:           `\"2001:db8:0:0:0:0:2:2\""`,
			givenNew:           "value",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"old and new values not empty, record type is AAAA, IPs match (given is full IPv6)": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8::2:1", "2001:db8::2:2"},
			givenOld:           `\"2001:db8:0:0:0:0:2:2\""`,
			givenNew:           "value",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"old and new values not empty, record type is AAAA, no match found": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8::2:2", "2001:db8::2:3"},
			givenOld:           `\2001:db8::2:5\""`,
			givenNew:           "value",
			givenRecordType:    "AAAA",
			expected:           false,
		},
		"old value is empty, record type is AAAA, IPs match (both are short IPv6)": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8::2:2", "2001:db8::2:3"},
			givenOld:           "",
			givenNew:           "\"2001:db8::2:1\"",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"old value is empty, record type is AAAA, IPs do not match": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8::2:2", "2001:db8::2:3"},
			givenOld:           "",
			givenNew:           "\"2001:db8::2:2\"",
			givenRecordType:    "AAAA",
			expected:           false,
		},
		"old value is empty, record type is AAAA, IPs match (target is full IPv6)": {
			givenOldTargetList: []string{"2001:db8:0:0:0:0:2:2", "2001:db8:0:0:0:0:2:3"},
			givenNewTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenOld:           "",
			givenNew:           "\"2001:db8::2:2\"",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"old value is empty, record type is AAAA, IPs match (both are full IPv6)": {
			givenOldTargetList: []string{"2001:db8:0:0:0:0:2:2", "2001:db8:0:0:0:0:2:3"},
			givenNewTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenOld:           "",
			givenNew:           "\"2001:db8:0:0:0:0:2:2\"",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"old value is empty, record type is AAAA, IPs match (given is full IPv6)": {
			givenOldTargetList: []string{"2001:db8::2:1", "2001:db8::2:2"},
			givenNewTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenOld:           "",
			givenNew:           "\"2001:db8:0:0:0:0:2:2\"",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"new value is empty, record type is AAAA, IPs match": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8::2:2", "2001:db8::2:3"},
			givenOld:           "",
			givenNew:           "\"2001:db8::2:1\"",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"new  value is empty, record type is AAAA, IPs do not match": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8::2:2", "2001:db8::2:3"},
			givenOld:           "\"2001:db8::2:5\"",
			givenNew:           "",
			givenRecordType:    "AAAA",
			expected:           false,
		},
		"new  value is empty, record type is AAAA, IPs match (both are full IPv6)": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8:0:0:0:0:2:2", "2001:db8:0:0:0:0:2:3"},
			givenOld:           `\"2001:db8:0:0:0:0:2:2\""`,
			givenNew:           "value",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"new  value is empty, record type is AAAA, IPs match (given is full IPv6)": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8::2:1", "2001:db8::2:2"},
			givenOld:           `\"2001:db8:0:0:0:0:2:2\""`,
			givenNew:           "value",
			givenRecordType:    "AAAA",
			expected:           true,
		},
		"new  value is empty, record type is AAAA, no match found": {
			givenOldTargetList: []string{"2001:db8::2:0", "2001:db8::2:1"},
			givenNewTargetList: []string{"2001:db8::2:2", "2001:db8::2:3"},
			givenOld:           `\2001:db8::2:5\""`,
			givenNew:           "value",
			givenRecordType:    "AAAA",
			expected:           false,
		},
		"old and new values not empty, record type is CAA, values match (no quotations)": {
			givenOldTargetList: []string{"example1.com.  CAA 0 issue akamairecord1.com", "example2.com.  CAA 0 issue \"akamairecord2.com\""},
			givenNewTargetList: []string{"example3.com.  CAA 0 issue \"akamairecord3.com\"", "example4.com.  CAA 0 issue \"akamairecord4.com\""},
			givenOld:           "example3.com.  CAA 0 issue akamairecord3.com",
			givenNew:           "value",
			givenRecordType:    "CAA",
			expected:           true,
		},
		"old and new values not empty, record type is CAA, values match": {
			givenOldTargetList: []string{"\"example1.com.\"  \"CAA\"\" 0 issue \"akamairecord1.com\"", "example2.com.  CAA 0 issue \"akamairecord2.com\""},
			givenNewTargetList: []string{"example3.com.  CAA 0 issue \"akamairecord3.com\"", "example4.com.  CAA 0 issue \"akamairecord4.com\""},
			givenOld:           "\\\"example3.com.  CAA \"0\" \"issue\" \\\"akamairecord3.com\\\"\\\"",
			givenNew:           "value",
			givenRecordType:    "CAA",
			expected:           true,
		},
		"old and new values not empty, record type is CAA, no match (invalid quotations)": {
			givenOldTargetList: []string{"example1.com.  CAA 0 issue akamairecord1.com", "example2.com.  CAA 0 issue \"akamairecord2.com\""},
			givenNewTargetList: []string{"example3.com.  CAA 0 issue akamairecord3.com", "example4.com.  CAA 0 issue \"akamairecord4.com\""},
			givenOld:           "example3.com.  \\\\\"CAA\" 0 issue akamairecord3.com",
			givenNew:           "value",
			givenRecordType:    "CAA",
			expected:           false,
		},
		"old value is empty, record type is CAA, values match (no quotations)": {
			givenOldTargetList: []string{"example1.com.  CAA 0 issue akamairecord1.com", "example2.com.  CAA 0 issue \"akamairecord2.com\""},
			givenNewTargetList: []string{"example3.com.  CAA 0 issue \"akamairecord3.com\"", "example4.com.  CAA 0 issue \"akamairecord4.com\""},
			givenOld:           "",
			givenNew:           "example1.com.  CAA 0 issue akamairecord1.com",
			givenRecordType:    "CAA",
			expected:           true,
		},
		"old value is empty, record type is CAA, values match": {
			givenOldTargetList: []string{"\"example1.com.\"  \"CAA\"\" 0 issue \"akamairecord1.com\"", "example2.com.  CAA 0 issue \"akamairecord2.com\""},
			givenNewTargetList: []string{"example3.com.  CAA 0 issue \"akamairecord3.com\"", "example4.com.  CAA 0 issue \"akamairecord4.com\""},
			givenOld:           "",
			givenNew:           "\"example1.com.  CAA \"0\" \"issue\" \"akamairecord1.com\"\"",
			givenRecordType:    "CAA",
			expected:           true,
		},
		"old value is empty, record type is CAA, no match (invalid quotations)": {
			givenOldTargetList: []string{"example1.com.  CAA 0 issue akamairecord1.com", "example2.com.  CAA 0 issue \"akamairecord2.com\""},
			givenNewTargetList: []string{"example3.com.  CAA 0 issue akamairecord3.com", "example4.com.  CAA 0 issue \"akamairecord4.com\""},
			givenOld:           "",
			givenNew:           "\"example1.com.  CAA \\\"0\" \"issue\" \"akamairecord1.com\"\"",
			givenRecordType:    "CAA",
			expected:           false,
		},
		"new value is empty, record type is CAA, values match (no quotations)": {
			givenOldTargetList: []string{"example1.com.  CAA 0 issue akamairecord1.com", "example2.com.  CAA 0 issue \"akamairecord2.com\""},
			givenNewTargetList: []string{"example3.com.  CAA 0 issue \"akamairecord3.com\"", "example4.com.  CAA 0 issue \"akamairecord4.com\""},
			givenOld:           "example3.com.  CAA 0 issue akamairecord3.com",
			givenNew:           "",
			givenRecordType:    "CAA",
			expected:           true,
		},
		"new value is empty, record type is CAA, values match": {
			givenOldTargetList: []string{"\"example1.com.\"  \"CAA\"\" 0 issue \"akamairecord1.com\"", "example2.com.  CAA 0 issue \"akamairecord2.com\""},
			givenNewTargetList: []string{"example3.com.  CAA 0 issue \"akamairecord3.com\"", "example4.com.  CAA 0 issue \"akamairecord4.com\""},
			givenOld:           "\\\"example3.com.  CAA \"0\" \"issue\" \\\"akamairecord3.com\\\"\\\"",
			givenNew:           "",
			givenRecordType:    "CAA",
			expected:           true,
		},
		"new value is empty, record type is CAA, no match (invalid quotations)": {
			givenOldTargetList: []string{"example1.com.  CAA 0 issue akamairecord1.com", "example2.com.  CAA 0 issue \"akamairecord2.com\""},
			givenNewTargetList: []string{"example3.com.  CAA 0 issue akamairecord3.com", "example4.com.  CAA 0 issue \"akamairecord4.com\""},
			givenOld:           "example3.com.  \\\\\"CAA\" 0 issue akamairecord3.com",
			givenNew:           "",
			givenRecordType:    "CAA",
			expected:           false,
		},
		"both empty": {
			givenOldTargetList: []string{},
			givenNewTargetList: []string{},
			givenOld:           "aaa",
			givenNew:           "aaa",
			givenRecordType:    "",
			expected:           false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := diffQuotedDNSRecord(
				test.givenOldTargetList,
				test.givenNewTargetList,
				test.givenOld,
				test.givenNew,
				test.givenRecordType,
				akamai.Log())
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestPadCoordinates(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected string
	}{
		"too short string passed": {
			given:    "asdasd",
			expected: "",
		},
		"valid string passed": {
			given:    "A B C D E F G H 1.234 2.345 3.456 4.567",
			expected: "A B C D E F G H 1.23m 2.35m 3.46m 4.57m",
		},
		"invalid float values": {
			given:    "A B C D E F G H W X Y Z",
			expected: "A B C D E F G H 0.00m 0.00m 0.00m 0.00m",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := padCoordinates(test.given, akamai.Log())
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestBuildRecordsList(t *testing.T) {
	tests := map[string]struct {
		givenRecordType string
		givenTarget     []interface{}
		expected        []string
		withError       bool
	}{
		"target record is of type from outside of the map": {
			givenRecordType: "ABC",
			givenTarget:     []interface{}{"192.168.0.1"},
			expected:        []string{"192.168.0.1"},
		},
		"target contains AAAA record": {
			givenRecordType: "AAAA",
			givenTarget:     []interface{}{"192.168.0.1"},
			expected:        []string{"0000:0000:0000:0000:0000:ffff:c0a8:0001"},
		},
		"target contains LOC record": {
			givenRecordType: "LOC",
			givenTarget:     []interface{}{"A B C D E F G H 1.234 2.345 3.456 4.567"},
			expected:        []string{"A B C D E F G H 1.23m 2.35m 3.46m 4.57m"},
		},
		"target contains SFP record without quotes": {
			givenRecordType: "SPF",
			givenTarget:     []interface{}{"v=spf1 ip4:1.2. 3.4 ip4:2.3. 4.5"},
			expected:        []string{"\"v=spf1 ip4:1.2. 3.4 ip4:2.3. 4.5\""},
		},
		"target contains SFP record with quotes": {
			givenRecordType: "SPF",
			givenTarget:     []interface{}{"\"v=spf1 ip4:1.2. 3.4 ip4:2.3. 4.5\""},
			expected:        []string{"\"v=spf1 ip4:1.2. 3.4 ip4:2.3. 4.5\""},
		},
		"target contains TXT record without quotes": {
			givenRecordType: "SPF",
			givenTarget:     []interface{}{"ABC123"},
			expected:        []string{"\"ABC123\""},
		},
		"target contains TXT record with embeddedescaped characters": {
			givenRecordType: "TXT",
			givenTarget:     []interface{}{"A\\\\\\\"BC123"},
			expected:        []string{"\"A\\\"BC123\""},
		},
		"target contains CAA record without quotes": {
			givenRecordType: "CAA",
			givenTarget:     []interface{}{"0 abc test.akamai.com\"\"\""},
			expected:        []string{"0 abc \"test.akamai.com\""},
		},
		"target contains invalid CAA record": {
			givenRecordType: "CAA",
			givenTarget:     []interface{}{"0 abc"},
			withError:       true,
		},
		"target contains different record type, ends with '.'": {
			givenRecordType: "CNAME",
			givenTarget:     []interface{}{"ABC."},
			expected:        []string{"ABC."},
		},
		"target contains different record type, ends without '.'": {
			givenRecordType: "CNAME",
			givenTarget:     []interface{}{"ABC"},
			expected:        []string{"ABC."},
		},
		"target contains invalid type": {
			givenRecordType: "CNAME",
			givenTarget:     []interface{}{1},
			withError:       true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := buildRecordsList(test.givenTarget, test.givenRecordType, akamai.Log())
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestNewRecordCreate(t *testing.T) {
	tests := map[string]struct {
		givenRecordType string
		givenData       *schema.ResourceData
		givenTarget     []interface{}
		givenHost       string
		givenTTL        int
		expected        dnsv2.RecordBody
		withError       bool
	}{
		"record type AFSDB, with dot": {
			givenRecordType: "AFSDB",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"subtype": 234,
			}),
			givenTarget: []interface{}{"AAA."},
			givenHost:   "test-host",
			givenTTL:    123,
			expected:    dnsv2.RecordBody{Name: "test-host", RecordType: "AFSDB", TTL: 123, Target: []string{"234 AAA."}},
		},

		"record type AFSDB, without dot": {
			givenRecordType: "AFSDB",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"subtype": 234,
			}),
			givenTarget: []interface{}{"AAA"},
			givenHost:   "test-host",
			givenTTL:    123,
			expected:    dnsv2.RecordBody{Name: "test-host", RecordType: "AFSDB", TTL: 123, Target: []string{"234 AAA."}},
		},
		"record type DNSKEY": {
			givenRecordType: "DNSKEY",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"flags":     1,
				"protocol":  2,
				"algorithm": 3,
				"key":       "test-key",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "DNSKEY", TTL: 123, Target: []string{"1 2 3 test-key"}},
		},
		"record type DS": {
			givenRecordType: "DS",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"digest_type": 1,
				"keytag":      2,
				"algorithm":   3,
				"digest":      "test-digest",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "DS", TTL: 123, Target: []string{"2 3 1 test-digest"}},
		},
		"record type HINFO": {
			givenRecordType: "HINFO",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"hardware": "\\\"test-hardware",
				"software": "\\\"test-software",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "HINFO", TTL: 123, Target: []string{"\"test-hardware\" \"test-software\""}},
		},
		"record type LOC": {
			givenRecordType: "LOC",
			givenData:       schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{}),
			givenTarget:     []interface{}{"2", "1"},
			givenHost:       "test-host",
			givenTTL:        123,
			expected:        dnsv2.RecordBody{Name: "test-host", RecordType: "LOC", TTL: 123, Target: []string{"1", "2"}},
		},
		"record type NAPTR": {
			givenRecordType: "NAPTR",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"flagsnaptr":  1,
				"order":       2,
				"preference":  3,
				"regexp":      "test-regexp",
				"replacement": "test-replacement",
				"service":     "test-service",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "NAPTR", TTL: 123, Target: []string{"2 3 \"1\" \"test-service\" \"test-regexp\" test-replacement"}},
		},
		"record type NSEC3": {
			givenRecordType: "NSEC3",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"flags":                  1,
				"algorithm":              2,
				"iterations":             3,
				"next_hashed_owner_name": "test-next_hashed_owner_name",
				"salt":                   "test-salt",
				"type_bitmaps":           "test-type_bitmaps",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "NSEC3", TTL: 123, Target: []string{"2 1 3 test-salt test-next_hashed_owner_name test-type_bitmaps"}},
		},
		"record type NSEC3PARAM": {
			givenRecordType: "NSEC3PARAM",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"flags":      1,
				"algorithm":  2,
				"iterations": 3,
				"salt":       "test-salt",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "NSEC3PARAM", TTL: 123, Target: []string{"2 1 3 test-salt"}},
		},
		"record type RP": {
			givenRecordType: "RP",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"mailbox": "test-mailbox",
				"txt":     "test-txt",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "RP", TTL: 123, Target: []string{"test-mailbox. test-txt."}},
		},
		"record type RRSIG": {
			givenRecordType: "RRSIG",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"expiration":   "test-expiration",
				"inception":    "test-inception",
				"original_ttl": 1,
				"algorithm":    2,
				"labels":       3,
				"keytag":       4,
				"signature":    "test-signature",
				"signer":       "test-signer",
				"type_covered": "test-type_covered",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "RRSIG", TTL: 123, Target: []string{"test-type_covered 2 3 1 test-expiration test-inception 4 test-signer test-signature"}},
		},
		"record type SRV": {
			givenRecordType: "SRV",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"priority": 1,
				"weight":   2,
				"port":     3,
			}),
			givenTarget: []interface{}{"B.", "A"},
			givenHost:   "test-host",
			givenTTL:    123,
			expected:    dnsv2.RecordBody{Name: "test-host", RecordType: "SRV", TTL: 123, Target: []string{"1 2 3 A.", "1 2 3 B."}},
		},
		"record type SSHFP": {
			givenRecordType: "SSHFP",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"algorithm":        1,
				"fingerprint_type": 2,
				"fingerprint":      "test-fingerprint",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "SSHFP", TTL: 123, Target: []string{"1 2 test-fingerprint"}},
		},
		"record type SOA": {
			givenRecordType: "SOA",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"name_server":   "test-name_server",
				"email_address": "test-email_address",
				"serial":        1,
				"refresh":       2,
				"retry":         3,
				"expiry":        4,
				"nxdomain_ttl":  5,
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "SOA", TTL: 123, Target: []string{"test-name_server test-email_address. 1 2 3 4 5"}},
		},
		"record type AKAMAITLC": {
			givenRecordType: "AKAMAITLC",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"dns_name":    "test-dns_name",
				"answer_type": "test-answer_type",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "AKAMAITLC", TTL: 123, Target: []string{"test-answer_type test-dns_name"}},
		},
		"record type CERT without type_menmonic": {
			givenRecordType: "CERT",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"type_mnemonic": "",
				"type_value":    1,
				"keytag":        2,
				"algorithm":     3,
				"certificate":   "test-certificate",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "CERT", TTL: 123, Target: []string{"1 2 3 test-certificate"}},
		},
		"record type CERT with type_mnemonic": {
			givenRecordType: "CERT",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"type_mnemonic": "test-type_mnemonic",
				"type_value":    1,
				"keytag":        2,
				"algorithm":     3,
				"certificate":   "test-certificate",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "CERT", TTL: 123, Target: []string{"test-type_mnemonic 2 3 test-certificate"}},
		},
		"record type TLSA": {
			givenRecordType: "TLSA",
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"usage":       1,
				"selector":    2,
				"match_type":  3,
				"certificate": "test-certificate",
			}),
			givenHost: "test-host",
			givenTTL:  123,
			expected:  dnsv2.RecordBody{Name: "test-host", RecordType: "TLSA", TTL: 123, Target: []string{"1 2 3 test-certificate"}},
		},
		"different record type": {
			givenRecordType: "INVALID",
			givenData:       schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{}),
			withError:       true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := newRecordCreate(
				context.Background(),
				nil,
				test.givenData,
				test.givenRecordType,
				test.givenTarget,
				test.givenHost,
				test.givenTTL,
				akamai.Log(),
			)
			if test.withError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expected, res)
		})
	}
}

func TestValidateRecord(t *testing.T) {
	tests := map[string]struct {
		givenData *schema.ResourceData
		withError bool
	}{
		"type A, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "A",
				"name":       "test-name",
				"ttl":        1,
				"target":     []interface{}{"ABC"},
			}),
			withError: false,
		},
		"type A, missing host": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "A",
				"ttl":        1,
				"target":     []interface{}{"ABC"},
			}),
			withError: true,
		},
		"type A, missing ttl": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "A",
				"name":       "test-name",
				"target":     []interface{}{"ABC"},
			}),
			withError: true,
		},
		"type A, missing target": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "A",
				"name":       "test-name",
				"ttl":        1,
			}),
			withError: true,
		},
		"type AFSDB, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "AFSDB",
				"subtype":    1,
				"target":     []interface{}{"ABC"},
			}),
			withError: false,
		},
		"type AFSDB, missing subtype": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "AFSDB",
				"target":     []interface{}{"ABC"},
			}),
			withError: true,
		},
		"type DNSKEY, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "DNSKEY",
				"flags":      0,
				"protocol":   1,
				"algorithm":  5,
				"key":        "test-key",
				"ttl":        1,
			}),
			withError: false,
		},
		"type DNSKEY, invalid ttl": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "DNSKEY",
				"flags":      0,
				"protocol":   1,
				"algorithm":  5,
				"key":        "test-key",
				"ttl":        0,
			}),
			withError: true,
		},
		"type DNSKEY, invalid protocol": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "DNSKEY",
				"flags":      0,
				"protocol":   0,
				"algorithm":  5,
				"key":        "test-key",
				"ttl":        1,
			}),
			withError: true,
		},
		"type DNSKEY, invalid flag": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "DNSKEY",
				"flags":      5,
				"protocol":   1,
				"algorithm":  5,
				"key":        "test-key",
				"ttl":        1,
			}),
			withError: true,
		},
		"type DNSKEY, invalid algorithm": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "DNSKEY",
				"flags":      0,
				"protocol":   1,
				"algorithm":  10,
				"key":        "test-key",
				"ttl":        1,
			}),
			withError: true,
		},
		"type DNSKEY, empty key": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "DNSKEY",
				"flags":      0,
				"protocol":   1,
				"algorithm":  5,
				"key":        "",
				"ttl":        1,
			}),
			withError: true,
		},
		"type DS, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "DS",
				"digest_type": 1,
				"keytag":      1,
				"algorithm":   1,
				"digest":      "test-digest",
			}),
			withError: false,
		},
		"type DS, empty digest type": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "DS",
				"keytag":     1,
				"algorithm":  1,
				"digest":     "test-digest",
			}),
			withError: true,
		},
		"type DS, empty keytag": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "DS",
				"digest_type": 1,
				"algorithm":   1,
				"digest":      "test-digest",
			}),
			withError: true,
		},
		"type DS, empty algorithm": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "DS",
				"digest_type": 1,
				"keytag":      1,
				"digest":      "test-digest",
			}),
			withError: true,
		},
		"type DS, empty digest": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "DS",
				"digest_type": 1,
				"keytag":      1,
				"algorithm":   1,
			}),
			withError: true,
		},
		"type HINFO, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "HINFO",
				"hardware":   "test-hardware",
				"software":   "test-software",
			}),
			withError: false,
		},
		"type HINFO, missing hardware": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "HINFO",
				"software":   "test-software",
			}),
			withError: true,
		},
		"type HINFO, missing software": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "HINFO",
				"hardware":   "test-hardware",
			}),
			withError: true,
		},
		"type MX, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "MX",
				"name":       "test-name",
				"ttl":        1,
				"target":     []interface{}{"ABC"},
				"priority":   5,
			}),
			withError: false,
		},
		"type MX, invalid priority": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "MX",
				"name":       "test-name",
				"ttl":        1,
				"target":     []interface{}{"ABC"},
				"priority":   -1,
			}),
			withError: true,
		},
		"type MX, invalid targets": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "MX",
				"name":       "test-name",
				"ttl":        1,
				"priority":   1,
			}),
			withError: true,
		},
		"type MX, invalid base": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "MX",
				"ttl":        1,
				"target":     []interface{}{"ABC"},
				"priority":   1,
			}),
			withError: true,
		},
		"type NAPTR, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "NAPTR",
				"name":        "test-name",
				"ttl":         1,
				"flagsnaptr":  "test-flagsnaptr",
				"order":       1,
				"preference":  1,
				"regexp":      "test-regexp",
				"replacement": "test-replacement",
				"service":     "test-service",
			}),
			withError: false,
		},
		"type NAPTR, invalid base": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "NAPTR",
				"ttl":         1,
				"flagsnaptr":  "test-flagsnaptr",
				"order":       1,
				"preference":  1,
				"regexp":      "test-regexp",
				"replacement": "test-replacement",
				"service":     "test-service",
			}),
			withError: true,
		},
		"type NAPTR, missing flagsnaptr": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "NAPTR",
				"name":        "test-name",
				"ttl":         1,
				"order":       1,
				"preference":  1,
				"regexp":      "test-regexp",
				"replacement": "test-replacement",
				"service":     "test-service",
			}),
			withError: true,
		},
		"type NAPTR, invalid order": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "NAPTR",
				"name":        "test-name",
				"ttl":         1,
				"flagsnaptr":  "test-flagsnaptr",
				"order":       -1,
				"preference":  1,
				"regexp":      "test-regexp",
				"replacement": "test-replacement",
				"service":     "test-service",
			}),
			withError: true,
		},
		"type NAPTR, missing preference": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "NAPTR",
				"name":        "test-name",
				"ttl":         1,
				"flagsnaptr":  "test-flagsnaptr",
				"order":       1,
				"preference":  0,
				"regexp":      "test-regexp",
				"replacement": "test-replacement",
				"service":     "test-service",
			}),
			withError: true,
		},
		"type NAPTR, missing regexp": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "NAPTR",
				"name":        "test-name",
				"ttl":         1,
				"flagsnaptr":  "test-flagsnaptr",
				"order":       1,
				"preference":  1,
				"replacement": "test-replacement",
				"service":     "test-service",
			}),
			withError: true,
		},
		"type NAPTR, missing replacement": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "NAPTR",
				"name":       "test-name",
				"ttl":        1,
				"flagsnaptr": "test-flagsnaptr",
				"order":      1,
				"preference": 1,
				"regexp":     "test-regexp",
				"service":    "test-service",
			}),
			withError: true,
		},
		"type NAPTR, missing service": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "NAPTR",
				"name":        "test-name",
				"ttl":         1,
				"flagsnaptr":  "test-flagsnaptr",
				"order":       1,
				"preference":  1,
				"regexp":      "test-regexp",
				"replacement": "test-replacement",
			}),
			withError: true,
		},
		"type NSEC3, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":             "NSEC3",
				"name":                   "test-name",
				"ttl":                    1,
				"flags":                  1,
				"algorithm":              1,
				"iterations":             1,
				"next_hashed_owner_name": "test-next_hashed_owner_name",
				"salt":                   "test-salt",
				"type_bitmaps":           "test-type_bitmaps",
			}),
			withError: false,
		},
		"type NSEC3, invalid base": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":             "NSEC3",
				"ttl":                    1,
				"flags":                  1,
				"algorithm":              1,
				"iterations":             1,
				"next_hashed_owner_name": "test-next_hashed_owner_name",
				"salt":                   "test-salt",
				"type_bitmaps":           "test-type_bitmaps",
			}),
			withError: true,
		},
		"type NSEC3, invalid flags": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":             "NSEC3",
				"name":                   "test-name",
				"ttl":                    1,
				"flags":                  2,
				"algorithm":              1,
				"iterations":             1,
				"next_hashed_owner_name": "test-next_hashed_owner_name",
				"salt":                   "test-salt",
				"type_bitmaps":           "test-type_bitmaps",
			}),
			withError: true,
		},
		"type NSEC3, invalid algorithm": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":             "NSEC3",
				"name":                   "test-name",
				"ttl":                    1,
				"flags":                  1,
				"algorithm":              2,
				"iterations":             1,
				"next_hashed_owner_name": "test-next_hashed_owner_name",
				"salt":                   "test-salt",
				"type_bitmaps":           "test-type_bitmaps",
			}),
			withError: true,
		},
		"type NSEC3, missing iterations": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":             "NSEC3",
				"name":                   "test-name",
				"ttl":                    1,
				"flags":                  1,
				"algorithm":              1,
				"next_hashed_owner_name": "test-next_hashed_owner_name",
				"salt":                   "test-salt",
				"type_bitmaps":           "test-type_bitmaps",
			}),
			withError: true,
		},
		"type NSEC3, missing next_hashed_owner_name": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "NSEC3",
				"name":         "test-name",
				"ttl":          1,
				"flags":        1,
				"algorithm":    1,
				"iterations":   1,
				"salt":         "test-salt",
				"type_bitmaps": "test-type_bitmaps",
			}),
			withError: true,
		},
		"type NSEC3, missing salt": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":             "NSEC3",
				"name":                   "test-name",
				"ttl":                    1,
				"flags":                  1,
				"algorithm":              1,
				"iterations":             1,
				"next_hashed_owner_name": "test-next_hashed_owner_name",
				"type_bitmaps":           "test-type_bitmaps",
			}),
			withError: true,
		},
		"type NSEC3, missing type_bitmaps": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":             "NSEC3",
				"name":                   "test-name",
				"ttl":                    1,
				"flags":                  1,
				"algorithm":              1,
				"iterations":             1,
				"next_hashed_owner_name": "test-next_hashed_owner_name",
				"salt":                   "test-salt",
			}),
			withError: true,
		},
		"type NSEC3PARAM, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "NSEC3PARAM",
				"name":       "test-name",
				"ttl":        1,
				"flags":      1,
				"algorithm":  1,
				"iterations": 1,
				"salt":       "test-salt",
			}),
			withError: false,
		},
		"type NSEC3PARAM, invalid base": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "NSEC3PARAM",
				"ttl":        1,
				"flags":      1,
				"algorithm":  1,
				"iterations": 1,
				"salt":       "test-salt",
			}),
			withError: true,
		},
		"type NSEC3PARAM, invalid flags": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "NSEC3PARAM",
				"name":       "test-name",
				"ttl":        1,
				"flags":      2,
				"algorithm":  1,
				"iterations": 1,
				"salt":       "test-salt",
			}),
			withError: true,
		},
		"type NSEC3PARAM, invalid algorithm": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "NSEC3PARAM",
				"name":       "test-name",
				"ttl":        1,
				"flags":      1,
				"algorithm":  2,
				"iterations": 1,
				"salt":       "test-salt",
			}),
			withError: true,
		},
		"type NSEC3PARAM, missing iterations": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "NSEC3PARAM",
				"name":       "test-name",
				"ttl":        1,
				"flags":      1,
				"algorithm":  1,
				"salt":       "test-salt",
			}),
			withError: true,
		},
		"type NSEC3PARAM, missing salt": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "NSEC3PARAM",
				"name":       "test-name",
				"ttl":        1,
				"flags":      1,
				"algorithm":  1,
				"iterations": 1,
			}),
			withError: true,
		},
		"type RP, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "RP",
				"name":       "test-name",
				"ttl":        1,
				"mailbox":    "test-mailbox",
				"txt":        "test-txt",
			}),
			withError: false,
		},
		"type RP, invalid base": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "RP",
				"ttl":        1,
				"mailbox":    "test-mailbox",
				"txt":        "test-txt",
			}),
			withError: true,
		},
		"type RP, missing mailbox": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "RP",
				"name":       "test-name",
				"ttl":        1,
				"txt":        "test-txt",
			}),
			withError: true,
		},
		"type RP, missing txt": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "RP",
				"name":       "test-name",
				"ttl":        1,
				"mailbox":    "test-mailbox",
			}),
			withError: true,
		},
		"type RRSIG, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"name":         "test-name",
				"ttl":          1,
				"expiration":   "test-expiration",
				"inception":    "test-inception",
				"original_ttl": 1,
				"algorithm":    1,
				"labels":       1,
				"keytag":       1,
				"signature":    "test",
				"signer":       "test",
				"type_covered": "test",
			}),
			withError: false,
		},
		"type RRSIG, invalid base": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"ttl":          1,
				"expiration":   "test-expiration",
				"inception":    "test-inception",
				"original_ttl": 1,
				"algorithm":    1,
				"labels":       1,
				"keytag":       1,
				"signature":    "test",
				"signer":       "test",
				"type_covered": "test",
			}),
			withError: true,
		},
		"type RRSIG, missing inception": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"name":         "test-name",
				"ttl":          1,
				"expiration":   "test-expiration",
				"original_ttl": 1,
				"algorithm":    1,
				"labels":       1,
				"keytag":       1,
				"signature":    "test",
				"signer":       "test",
				"type_covered": "test",
			}),
			withError: true,
		},
		"type RRSIG, missing expiration": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"name":         "test-name",
				"ttl":          1,
				"inception":    "test-inception",
				"original_ttl": 1,
				"algorithm":    1,
				"labels":       1,
				"keytag":       1,
				"signature":    "test",
				"signer":       "test",
				"type_covered": "test",
			}),
			withError: true,
		},
		"type RRSIG, missing original ttl": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"name":         "test-name",
				"ttl":          1,
				"expiration":   "test-expiration",
				"inception":    "test-inception",
				"algorithm":    1,
				"labels":       1,
				"keytag":       1,
				"signature":    "test",
				"signer":       "test",
				"type_covered": "test",
			}),
			withError: true,
		},
		"type RRSIG, missing algorithm": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"name":         "test-name",
				"ttl":          1,
				"expiration":   "test-expiration",
				"inception":    "test-inception",
				"original_ttl": 1,
				"labels":       1,
				"keytag":       1,
				"signature":    "test",
				"signer":       "test",
				"type_covered": "test",
			}),
			withError: true,
		},
		"type RRSIG, missing labels": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"name":         "test-name",
				"ttl":          1,
				"expiration":   "test-expiration",
				"inception":    "test-inception",
				"original_ttl": 1,
				"algorithm":    1,
				"keytag":       1,
				"signature":    "test",
				"signer":       "test",
				"type_covered": "test",
			}),
			withError: true,
		},
		"type RRSIG, missing keytag": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"name":         "test-name",
				"ttl":          1,
				"expiration":   "test-expiration",
				"inception":    "test-inception",
				"original_ttl": 1,
				"algorithm":    1,
				"labels":       1,
				"signature":    "test",
				"signer":       "test",
				"type_covered": "test",
			}),
			withError: true,
		},
		"type RRSIG, missing signature": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"name":         "test-name",
				"ttl":          1,
				"expiration":   "test-expiration",
				"inception":    "test-inception",
				"original_ttl": 1,
				"algorithm":    1,
				"labels":       1,
				"keytag":       1,
				"signer":       "test",
				"type_covered": "test",
			}),
			withError: true,
		},
		"type RRSIG, missing signer": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"name":         "test-name",
				"ttl":          1,
				"expiration":   "test-expiration",
				"inception":    "test-inception",
				"original_ttl": 1,
				"algorithm":    1,
				"labels":       1,
				"keytag":       1,
				"signature":    "test",
				"type_covered": "test",
			}),
			withError: true,
		},
		"type RRSIG, missing type covered": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "RRSIG",
				"name":         "test-name",
				"ttl":          1,
				"expiration":   "test-expiration",
				"inception":    "test-inception",
				"original_ttl": 1,
				"algorithm":    1,
				"labels":       1,
				"keytag":       1,
				"signature":    "test",
				"signer":       "test",
			}),
			withError: true,
		},
		"type SSHFP, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":       "SSHFP",
				"name":             "test-name",
				"ttl":              1,
				"algorithm":        1,
				"fingerprint_type": 1,
				"fingerprint":      "test",
			}),
			withError: false,
		},
		"type SSHFP, invalid base": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":       "SSHFP",
				"ttl":              1,
				"algorithm":        1,
				"fingerprint_type": 1,
				"fingerprint":      "test",
			}),
			withError: true,
		},
		"type SSHFP, missing algorithm": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":       "SSHFP",
				"name":             "test-name",
				"ttl":              1,
				"fingerprint_type": 1,
				"fingerprint":      "test",
			}),
			withError: true,
		},
		"type SSHFP, missing fingerprint_type": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "SSHFP",
				"name":        "test-name",
				"ttl":         1,
				"algorithm":   1,
				"fingerprint": "test",
			}),
			withError: true,
		},
		"type SSHFP, invalid fingerprint": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":       "SSHFP",
				"name":             "test-name",
				"ttl":              1,
				"algorithm":        1,
				"fingerprint_type": 1,
				"fingerprint":      "null",
			}),
			withError: true,
		},
		"type SRV, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "SRV",
				"target":     []interface{}{"ABC"},
				"name":       "test-name",
				"ttl":        1,
				"priority":   1,
				"weight":     1,
				"port":       1,
			}),
			withError: false,
		},
		"type SRV, invalid base": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "SRV",
				"target":     []interface{}{"ABC"},
				"ttl":        1,
				"priority":   1,
				"weight":     1,
				"port":       1,
			}),
			withError: true,
		},
		"type SRV, invalid target": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "SRV",
				"name":       "test-name",
				"ttl":        1,
				"priority":   1,
				"weight":     1,
				"port":       1,
			}),
			withError: true,
		},
		"type SRV,invalid priority ": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "SRV",
				"target":     []interface{}{"ABC"},
				"name":       "test-name",
				"ttl":        1,
				"priority":   -1,
				"weight":     1,
				"port":       1,
			}),
			withError: true,
		},
		"type SRV, invalid weight": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "SRV",
				"target":     []interface{}{"ABC"},
				"name":       "test-name",
				"ttl":        1,
				"priority":   1,
				"weight":     -1,
				"port":       1,
			}),
			withError: true,
		},
		"type SRV, missing port": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "SRV",
				"target":     []interface{}{"ABC"},
				"name":       "test-name",
				"ttl":        1,
				"priority":   1,
				"weight":     1,
			}),
			withError: true,
		},
		"type SOA, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":    "SOA",
				"name_server":   "test",
				"email_address": "test",
				"refresh":       1,
				"retry":         1,
				"expiry":        1,
				"nxdomain_ttl":  1,
			}),
			withError: false,
		},
		"type SOA, missing nameserver": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":    "SOA",
				"email_address": "test",
				"refresh":       1,
				"retry":         1,
				"expiry":        1,
				"nxdomain_ttl":  1,
			}),
			withError: true,
		},
		"type SOA, missing email_address": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":   "SOA",
				"name_server":  "test",
				"refresh":      1,
				"retry":        1,
				"expiry":       1,
				"nxdomain_ttl": 1,
			}),
			withError: true,
		},
		"type SOA, missing refresh": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":    "SOA",
				"name_server":   "test",
				"email_address": "test",
				"retry":         1,
				"expiry":        1,
				"nxdomain_ttl":  1,
			}),
			withError: true,
		},
		"type SOA, missing retry": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":    "SOA",
				"name_server":   "test",
				"email_address": "test",
				"refresh":       1,
				"expiry":        1,
				"nxdomain_ttl":  1,
			}),
			withError: true,
		},
		"type SOA, missing expiry": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":    "SOA",
				"name_server":   "test",
				"email_address": "test",
				"refresh":       1,
				"retry":         1,
				"nxdomain_ttl":  1,
			}),
			withError: true,
		},
		"type SOA, missing nxdomain": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":    "SOA",
				"name_server":   "test",
				"email_address": "test",
				"refresh":       1,
				"retry":         1,
				"expiry":        1,
			}),
			withError: true,
		},
		"type AKAMAITLC": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "AKAMAITLC",
			}),
			withError: true,
		},
		"type CAA, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "CAA",
				"target":     []interface{}{"0 issue akamai.com"},
				"name":       "test-name",
				"ttl":        1,
			}),
			withError: false,
		},
		"type CAA, invalid target len": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "CAA",
				"target":     []interface{}{"0 akamai.com"},
				"name":       "test-name",
				"ttl":        1,
			}),
			withError: true,
		},
		"type CAA, invalid flag": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "CAA",
				"target":     []interface{}{"500 issue akamai.com"},
				"name":       "test-name",
				"ttl":        1,
			}),
			withError: true,
		},
		"type CAA, invalid tag": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "CAA",
				"target":     []interface{}{"0 issue-* akamai.com"},
				"name":       "test-name",
				"ttl":        1,
			}),
			withError: true,
		},
		"type CERT, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":    "CERT",
				"name":          "test-name",
				"ttl":           1,
				"type_mnemonic": "test",
				"type_value":    0,
				"certificate":   "test",
			}),
			withError: false,
		},
		"type CERT, invalid base": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":    "CERT",
				"ttl":           1,
				"type_mnemonic": "test",
				"type_value":    0,
				"certificate":   "test",
			}),
			withError: true,
		},
		"type CERT, empty type mnemonic and type value": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "CERT",
				"name":        "test-name",
				"ttl":         1,
				"certificate": "test",
			}),
			withError: true,
		},
		"type CERT, both type mnemonic and type_value are not empty": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":    "CERT",
				"name":          "test-name",
				"ttl":           1,
				"type_mnemonic": "test",
				"type_value":    1,
				"certificate":   "test",
			}),
			withError: true,
		},
		"type CERT, missing certificate": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":    "CERT",
				"name":          "test-name",
				"ttl":           1,
				"type_mnemonic": "test",
				"type_value":    0,
			}),
			withError: true,
		},
		"type TLSA, valid": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "TLSA",
				"name":        "test-name",
				"ttl":         1,
				"usage":       1,
				"certificate": "test",
			}),
			withError: false,
		},
		"type TLSA, invalid base": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "TLSA",
				"ttl":         1,
				"usage":       1,
				"certificate": "test",
			}),
			withError: true,
		},
		"type TLSA, missing certificate": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype": "TLSA",
				"name":       "test-name",
				"ttl":        1,
				"usage":      1,
			}),
			withError: true,
		},
		"type TLSA, missing usage": {
			givenData: schema.TestResourceDataRaw(t, resourceDNSv2Record().Schema, map[string]interface{}{
				"recordtype":  "TLSA",
				"name":        "test-name",
				"ttl":         1,
				"certificate": "test",
			}),
			withError: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := validateRecord(test.givenData)
			if test.withError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
