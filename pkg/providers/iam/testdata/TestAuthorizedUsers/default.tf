provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_authorized_users" "test" {}
