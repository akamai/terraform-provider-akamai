provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_allowed_apis" "test" {
  username             = "test"
  client_type          = "CLIENT"
  allow_account_switch = true
}
