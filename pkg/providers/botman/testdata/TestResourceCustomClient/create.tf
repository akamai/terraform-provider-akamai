provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_custom_client" "test" {
  config_id     = 43253
  custom_client = <<-EOF
{
  "testKey": "testValue3"
}
EOF
}