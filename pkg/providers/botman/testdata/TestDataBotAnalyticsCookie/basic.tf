provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_botman_bot_analytics_cookie" "test" {
  config_id = 43253
}