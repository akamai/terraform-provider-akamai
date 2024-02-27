provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_cloudlets_policy_activation" "test" {
  policy_id = 1234
  network   = "production"
  version   = 1
}

output "status" {
  value = akamai_cloudlets_policy_activation.test.status
}