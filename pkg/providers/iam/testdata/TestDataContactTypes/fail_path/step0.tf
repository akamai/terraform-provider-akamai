provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_contact_types" "test" {}
