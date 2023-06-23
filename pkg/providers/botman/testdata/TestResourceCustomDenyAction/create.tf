provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_custom_deny_action" "test" {
  config_id = 43253
  custom_deny_action = jsonencode(
    {
      "testKey" : "testValue3"
    }
  )
}