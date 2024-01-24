provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_cloudlets_policy_activation" "test" {
  policy_id = 1
}