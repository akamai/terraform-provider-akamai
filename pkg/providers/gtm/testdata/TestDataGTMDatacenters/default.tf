provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_datacenters" "test" {
  domain = "gtm_terra_testdomain.akadns.net"
}
