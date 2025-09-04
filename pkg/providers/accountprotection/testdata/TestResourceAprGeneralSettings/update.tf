provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_apr_general_settings" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  general_settings = jsonencode(
    {
      "testKey" : "updated_testValue3"
    }
  )
}