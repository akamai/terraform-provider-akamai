provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property_activation" "test" {
  property_id                    = "prp_test"
  network                        = "STAGING"
  contact                        = ["user@example.com"]
  version                        = 1
  auto_acknowledge_rule_warnings = true
}