provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property_activation" "test" {
  contact = ["user@example.com"]
  version = 1
}