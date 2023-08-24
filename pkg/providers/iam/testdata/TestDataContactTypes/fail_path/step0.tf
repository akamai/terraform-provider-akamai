provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_iam_contact_types" "test" {}
