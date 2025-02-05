provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_domain" "test" {
  name = "gtm_terra_testdomain.akadns.net"
  type = "weighted"
}
