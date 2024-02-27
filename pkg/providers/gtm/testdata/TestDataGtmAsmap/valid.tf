provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_asmap" "my_gtm_asmap" {
  domain   = "test.domain.net"
  map_name = "map1"
}
