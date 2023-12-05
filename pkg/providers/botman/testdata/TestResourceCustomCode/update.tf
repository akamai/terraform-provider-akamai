provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_botman_custom_code" "test" {
  config_id = 43253
  custom_code = jsonencode(
    {
      "testKey" : "updated_testValue3"
    }
  )
}