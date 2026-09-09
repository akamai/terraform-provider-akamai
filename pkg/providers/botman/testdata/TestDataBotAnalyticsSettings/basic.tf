provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_botman_bot_analytics_settings" "test" {
  config_id = 43253
}
