provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_apr_user_risk_response_strategy" "test" {
  config_id = 43253
}