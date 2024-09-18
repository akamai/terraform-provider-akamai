provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_cidr_blocks" "test" {
  cidr_blocks = [
    {
      comments = "test"
      enabled  = true
    }
  ]
}