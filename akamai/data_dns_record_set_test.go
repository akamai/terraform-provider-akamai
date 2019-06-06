package akamai

import (
	"fmt"
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
				Config: testAccDataSourceDNSRecordSet_basic("akavdev.net"),
				Check: resource.ComposeTestCheckFunc(
					//resource.TestCheckOutput("test_addrs", "4 example.com."),
					//resource.TestCheckResourceAttr(dataSourceName, "rdata.%", "4 example.com."),
					resource.TestCheckResourceAttr(dataSourceName, "host", "akavdev.net"),
				),
			},
		},
	})
}

func testAccDataSourceDNSRecordSet_basic(name string) string {
	return fmt.Sprintf(`
provider "akamai" {
  edgerc = "~/.edgerc"
  dns_section = "dns"
}

locals {
  zone = "akavdev.net"
}

resource "akamai_dns_record" "test" {
	zone = "${local.zone}"
	name = "%s"
	recordtype =  "A"
	active = true
	ttl = 300
	target = ["10.0.0.2","10.0.0.3"]
}


data "akamai_dns_record_set" "test" {
	  zone = "${local.zone}"
		host = "akavdev.net"
		record_type = "A"
	}

	output "test_addrs" {
		value = "${join(",", data.akamai_dns_record_set.test.rdata)}"
	}
`, name)
}
