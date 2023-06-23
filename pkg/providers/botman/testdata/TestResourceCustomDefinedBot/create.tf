provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_custom_defined_bot" "test" {
  config_id = 43253
  custom_defined_bot = jsonencode(
    {
      "testKey" : "testValue3"
    }
  )
}