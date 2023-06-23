provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_policy" "test" {
  policy_id = 1234
  version   = 3
}