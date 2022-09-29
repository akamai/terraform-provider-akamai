provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_transactional_endpoint" "test" {
  config_id              = 43253
  security_policy_id     = "AAAA_81230"
  operation_id           = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
  transactional_endpoint = <<-EOF
{
  "testKey": "updated_testValue3"
}
EOF
}