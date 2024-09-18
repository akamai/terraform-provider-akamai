provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_cidr_blocks" "test" {
  cidr_blocks = [
    {
      cidr_block = "128.1.2.5/24"
      comments   = "test"
    }
  ]
}