provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_activation" "test" {
  property_id                    = "test"
  contact                        = ["user2@example.com"]
  version                        = 2
  auto_acknowledge_rule_warnings = true
  note                           = "property activation note for updating"
}