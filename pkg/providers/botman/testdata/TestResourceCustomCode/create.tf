provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_custom_code" "test" {
  config_id = 43253
  custom_code = jsonencode(
    {
      "testKey" : "testValue3"
    }
  )
}