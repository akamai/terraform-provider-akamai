provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_gtm_geomap" "testmap" {
  map_name = "mapTest"
}