provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_blocked_properties" "test" {
  group_id       = 1
  ui_identity_id = "user123"
}
