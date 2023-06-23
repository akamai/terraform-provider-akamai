provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_default_datacenter" "test" {
  domain     = "testdomain.net"
  datacenter = 5400
}

output "datacenter_id" {
  value = data.akamai_gtm_default_datacenter.test.datacenter_id
}
