provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_cidr_block" "test" {
  enabled = true
}