provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_bot_management_settings" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  bot_management_settings = jsonencode(
    {
      "testKey" : "updated_testValue3"
    }
  )
}