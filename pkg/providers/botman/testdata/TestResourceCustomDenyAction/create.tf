provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_custom_deny_action" "test" {
  config_id          = 43253
  custom_deny_action = <<-EOF
{
  "testKey": "testValue3"
}
EOF
}