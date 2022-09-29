provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_transactional_endpoint_protection" "test" {
  config_id                         = 43253
  transactional_endpoint_protection = <<-EOF
{
  "testKey": "updated_testValue3"
}
EOF
}