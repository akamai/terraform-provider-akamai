provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_activation" "test" {
  property_id = "test"
  contact     = ["user@example.com"]
  version     = 1
}