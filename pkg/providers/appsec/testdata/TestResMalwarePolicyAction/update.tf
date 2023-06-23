provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_malware_policy_action" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  malware_policy_id  = 135355
  action             = "alert"
  unscanned_action   = "deny"
}

