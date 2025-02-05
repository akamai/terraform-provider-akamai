provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_geomap" "testmap" {
  domain   = "gtm_terra_testdomain.akadns.net"
  map_name = "tfexample_geomap_1"
}
