provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_activation" "test" {
  property_id = "prp_test"
  version     = 1
}