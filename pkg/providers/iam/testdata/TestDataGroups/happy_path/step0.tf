provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_iam_groups" "test" {}
