provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_cloudlets_policy_activation" "test" {
  network = "staging"
}