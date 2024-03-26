provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_include_activation" "activation" {
  include_id    = "12345"
  contract_id   = "test_contract"
  group_id      = "test_group"
  version       = 3
  network       = "PRODUCTION"
  notify_emails = ["jbond@example.com"]
  note          = "test activation"
  compliance_record {
    noncompliance_reason_other {}
  }
  auto_acknowledge_rule_warnings = true
}