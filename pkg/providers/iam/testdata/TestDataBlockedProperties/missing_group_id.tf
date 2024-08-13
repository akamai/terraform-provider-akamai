provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_blocked_properties" "test" {
  contract_id    = "ctr_C-123"
  ui_identity_id = "user123"
}
