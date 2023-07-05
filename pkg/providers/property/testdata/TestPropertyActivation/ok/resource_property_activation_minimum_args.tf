provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property_activation" "test" {
  property_id = "test"
  contact     = ["user@example.com"]
  version     = 1
}