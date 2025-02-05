provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_asmap" "test" {
  domain = "gtm_terra_testdomain.akadns.net"
  name   = "tfexample_as_1"
  default_datacenter {
    datacenter_id = 5400
    nickname      = "default datacenter"
  }
}
