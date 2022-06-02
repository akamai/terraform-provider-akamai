provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cloudlets_application_load_balancer_activation" "test" {
  origin_id = "org_1"
  network   = "staging"
  version   = 2
}

output "status" {
  value = akamai_cloudlets_application_load_balancer_activation.test.status
}