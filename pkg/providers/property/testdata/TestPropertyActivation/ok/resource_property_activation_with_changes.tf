provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_activation" "test" {
  property_id                    = "test"
  contact                        = ["user@example1.com"]
  version                        = 1
  auto_acknowledge_rule_warnings = false
  note                           = "property activation note for creating"
}
