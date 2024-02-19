provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_shared_policy" "test" {
  policy_id = 1
  version   = 2
}