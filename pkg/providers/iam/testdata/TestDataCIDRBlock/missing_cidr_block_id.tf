provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_iam_cidr_block" "test" {
}
