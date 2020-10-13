provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_gtm_dafault_datacenter" "test" {
	domain = "testdomain.net"
	datacenter = 5400
}
  
output "datacenter_id" {
	value = "${join(",", data.akamai_gtm_default_datacenter.test.datacenter_id)}"
}
