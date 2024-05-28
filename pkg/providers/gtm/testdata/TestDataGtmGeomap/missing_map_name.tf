provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_geomap" "testmap" {
  domain = "test.geomap.domain.net"
}