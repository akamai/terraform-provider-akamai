package akamai

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					resource.TestCheckResourceAttr(dataSourceName, "host", "exampleterraform.io"),
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
	zone = "exampleterraform.io"
	name = "exampleterraform.io"
	recordtype =  "A"
	active = true
	ttl = 300
	target = ["10.0.0.2","10.0.0.3"]
}


data "akamai_dns_record_set" "test" {
	zone = "exampleterraform.io"
	host = "exampleterraform.io"
	record_type = "A"
}

output "test_addrs" {
	value = "${join(",", data.akamai_dns_record_set.test.rdata)}"
}
`
}
