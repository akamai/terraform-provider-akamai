provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_cidr_blocks" "test" {
  cidr_blocks = [
    {
      cidr_block = "128.5.6.5/24"
      comments   = "test"
      enabled    = true
    },
    {
      cidr_block = "128.5.6.6/24"
      enabled    = false
    },
    {
      cidr_block = "128.5.6.7/24"
      enabled    = true
      comments   = "test1234"
    },
    {
      cidr_block = "128.5.6.8/24"
      enabled    = false
      comments   = "abcd12345"
    },
    {
      cidr_block = "128.5.6.9/24"
      enabled    = true
    },
  ]
}