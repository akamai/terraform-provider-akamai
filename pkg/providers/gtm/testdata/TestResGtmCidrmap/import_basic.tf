provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_cidrmap" "test" {
  domain = "gtm_terra_testdomain.akadns.net"
  name   = "tfexample_cidrmap_1"
}