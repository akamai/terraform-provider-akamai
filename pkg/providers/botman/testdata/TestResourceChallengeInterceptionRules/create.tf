provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_challenge_interception_rules" "test" {
  config_id = 43253
  challenge_interception_rules = jsonencode(
    {
      "testKey" : "testValue3"
    }
  )
}