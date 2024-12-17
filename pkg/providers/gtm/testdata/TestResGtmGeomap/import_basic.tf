provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_gtm_geomap" "test" {
  domain = "gtm_terra_testdomain.akadns.net"
  name   = "tfexample_geomap_1"
  default_datacenter {
    datacenter_id = 5400
    nickname      = "default datacenter"
  }
}
