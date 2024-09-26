provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_cidr_block" "test" {
  cidr_block = "128.5.6.99/24"
  comments   = "test-updated"
  enabled    = false
}