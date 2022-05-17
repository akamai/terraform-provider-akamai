provider "akamai" {
  edgerc = "../../test/edgerc"
}


data "akamai_property" "prop" {
  name    = "property_name"
  version = 2
}
