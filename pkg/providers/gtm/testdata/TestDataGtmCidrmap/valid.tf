provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_gtm_cidrmap" "gtm_cidrmap" {
  domain   = "test.cidrmap.domain.net"
  map_name = "mapTest"
}
