provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_cidrmap" "gtm_cidrmap" {
  domain = "gtm_terra_testdomain.akadns.net"
}