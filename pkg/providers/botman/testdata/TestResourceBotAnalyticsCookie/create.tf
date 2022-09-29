provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_bot_analytics_cookie" "test" {
  config_id            = 43253
  bot_analytics_cookie = <<-EOF
{
  "testKey": "testValue3"
}
EOF
}