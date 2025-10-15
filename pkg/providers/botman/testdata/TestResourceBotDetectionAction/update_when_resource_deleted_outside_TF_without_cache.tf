provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_bot_detection_action" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  detection_id       = "cc9c3f89-e179-4892-89cf-d5e623ba9dc7"
  bot_detection_action = jsonencode(
    {
      "testKey" : "updated_testValue4"
    }
  )
}