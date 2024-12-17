provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_datacenter" "test" {
  domain = "gtm_terra_testdomain.akadns.net"
}
