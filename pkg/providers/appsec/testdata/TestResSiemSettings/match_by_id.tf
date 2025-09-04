provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_siem_settings" "test" {
  config_id                       = 43253
  enable_siem                     = true
  enable_for_all_policies         = false
  enable_botman_siem              = true
  include_ja4_fingerprint_to_siem = true
  siem_id                         = 1
  security_policy_ids             = ["12345"]
}

