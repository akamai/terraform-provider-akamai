provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property_activation" "test" {
  property_id                    = "test"
  contact                        = ["user@example.com"]
  version                        = 1
  auto_acknowledge_rule_warnings = true
  note                           = "property activation note for creating"
  compliance_record {
    noncompliance_reason_none {
      peer_reviewed_by = "user1@example.com"
      customer_email   = "user@example.com"
      unit_tested      = true
    }
  }
}
