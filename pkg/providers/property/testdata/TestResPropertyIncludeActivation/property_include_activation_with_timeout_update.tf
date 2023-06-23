provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_include_activation" "activation" {
  include_id    = "12345"
  contract_id   = "test_contract"
  group_id      = "test_group"
  version       = 3
  network       = "STAGING"
  notify_emails = ["jbond@example.com"]
  note          = "test activation"
  timeouts {
    default = "2h2m"
  }
}