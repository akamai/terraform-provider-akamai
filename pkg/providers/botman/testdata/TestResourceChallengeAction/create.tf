provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_challenge_action" "test" {
  config_id = 43253
  challenge_action = jsonencode(
    {
      "testKey" : "testValue3"
    }
  )
}