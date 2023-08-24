provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_iam_blocked_user_properties" "test" {
  identity_id        = "test_identity_id"
  group_id           = 12345
  blocked_properties = [1, 2, 3, 4, 5]
}
