provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_advanced_settings_request_body" "policy" {
  config_id          = 43253
  security_policy_id = "test_policy"
}
