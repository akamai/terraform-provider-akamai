provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_iam_cidr_blocks" "test" {
  cidr_blocks = [
    {
      cidr_block = "128.2.2.5/28"
      enabled    = false
    },
    {
      cidr_block = "128.2.2.6/28"
      enabled    = true
    },
    {
      cidr_block = "128.2.2.7/28"
      enabled    = false
      comments   = "test12345"
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