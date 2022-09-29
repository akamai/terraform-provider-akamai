provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_serve_alternate_action" "test" {
  config_id              = 43253
  serve_alternate_action = <<-EOF
{
  "testKey": "updated_testValue3"
}
EOF
}