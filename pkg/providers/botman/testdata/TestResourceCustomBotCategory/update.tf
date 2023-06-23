provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_custom_bot_category" "test" {
  config_id = 43253
  custom_bot_category = jsonencode(
    {
      "testKey" : "updated_testValue3"
    }
  )
}