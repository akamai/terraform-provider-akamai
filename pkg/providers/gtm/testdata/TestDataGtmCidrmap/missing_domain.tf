provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_gtm_cidrmap" "gtm_cidrmap" {
  map_name = "mapTest"
}