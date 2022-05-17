provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property_activation" "test" {
  property = "test"
  contact  = ["user@example.com"]
  version  = 1
}