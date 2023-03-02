provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_appsec_advanced_settings_attack_payload_logging" "policy" {
  config_id          = 43253
  security_policy_id = "test_policy"
}
