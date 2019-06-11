package akamai

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceDNSRecordSet_basic(t *testing.T) {
	dataSourceName := "data.akamai_dns_record_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAkamaiDNSv2RecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSRecordSet_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "host", "example.org"),
				),
			},
		},
	})
}

func testAccDataSourceDNSRecordSet_basic() string {
	return `provider "akamai" {
  dns_section = "dns"
}

resource "akamai_dns_record" "test" {
	zone = "example.org"
	name = "example.org"
	recordtype =  "A"
	active = true
	ttl = 300
	target = ["10.0.0.2","10.0.0.3"]
}


data "akamai_dns_record_set" "test" {
	zone = "example.org"
	host = "example.org"
	record_type = "A"
}

output "test_addrs" {
	value = "${join(",", data.akamai_dns_record_set.test.rdata)}"
}
`
}
