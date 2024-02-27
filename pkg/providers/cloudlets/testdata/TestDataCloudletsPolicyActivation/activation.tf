provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_policy_activation" "test" {
  policy_id = 1
  network   = "staging"
}