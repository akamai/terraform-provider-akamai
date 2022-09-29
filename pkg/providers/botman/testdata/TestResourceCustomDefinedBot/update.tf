provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_custom_defined_bot" "test" {
  config_id          = 43253
  custom_defined_bot = <<-EOF
{
  "testKey": "updated_testValue3"
}
EOF
}