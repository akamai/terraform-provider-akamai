provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_security_policy_protections" "test" {
  config_id                        = 12345
  security_policy_id               = ""
  apply_api_constraints            = true
  apply_application_layer_controls = true
  apply_botman_controls            = false
  apply_malware_controls           = true
  apply_network_layer_controls     = false
  apply_rate_controls              = true
  apply_reputation_controls        = false
  apply_slow_post_controls         = true
  apply_url_protection_controls    = true
}
