provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_bot_analytics_settings" "test" {
  config_id = 43253
  bot_analytics_settings = jsonencode(
    {
      "testKey" : "testValue3"
    }
  )
}

