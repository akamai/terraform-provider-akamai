provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_apr_user_risk_response_strategy" "test" {
  config_id = 43253
  user_risk_response_strategy = jsonencode(
    {
      "testKey" : "testValue3"
    }
  )
}