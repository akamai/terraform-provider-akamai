provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_bot_category_exception" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  bot_category_exception = jsonencode(
    {
      "testKey" : "testValue3"
    }
  )
}