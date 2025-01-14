provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_resources" "my_gtm_resources" {
  domain = "gtm_terra_testdomain.akadns.net"
}