provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_user" "test" {
  ui_identity_id = "asd-12345"
}
