provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_akamai_bot_category_action" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  category_id        = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
  akamai_bot_category_action = jsonencode(
    {
      "testKey" : "updated_testValue3"
    }
  )
}