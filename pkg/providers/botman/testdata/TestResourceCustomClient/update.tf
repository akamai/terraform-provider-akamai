provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_custom_client" "test" {
  config_id = 43253
  custom_client = jsonencode(
    {
      "testKey" : "updated_testValue3"
    }
  )
}