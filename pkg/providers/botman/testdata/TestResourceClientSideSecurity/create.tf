provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_client_side_security" "test" {
  config_id            = 43253
  client_side_security = <<-EOF
{
  "testKey": "testValue3"
}
EOF
}