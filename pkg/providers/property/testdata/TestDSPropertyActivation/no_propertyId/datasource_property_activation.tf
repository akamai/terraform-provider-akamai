provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_activation" "test" {
  version = 1
}