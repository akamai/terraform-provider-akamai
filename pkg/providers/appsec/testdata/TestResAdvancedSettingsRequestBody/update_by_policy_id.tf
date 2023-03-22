provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_advanced_settings_request_body" "policy" {
  config_id                     = 43253
  security_policy_id            = "test_policy"
  request_body_inspection_limit = 16
}