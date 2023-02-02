provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_activation" "test" {
  property_id = "test"
  version     = 1
  network     = "PRODUCTION"
}