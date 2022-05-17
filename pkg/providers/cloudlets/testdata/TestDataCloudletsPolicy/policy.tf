provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_policy" "test" {
  policy_id = 1234
}