provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_activation" "test" {
  property_id                    = "prp_test"
  contact                        = ["user@example.com"]
  version                        = 1
  note                           = "property activation note for importing"
  auto_acknowledge_rule_warnings = false
}