provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_cidr_block" "test" {
  cidr_block = "128.1.2.5/24"
}