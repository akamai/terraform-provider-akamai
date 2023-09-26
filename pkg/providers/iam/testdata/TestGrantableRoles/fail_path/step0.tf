provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_iam_grantable_roles" "test" {}
