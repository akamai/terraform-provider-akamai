provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_appsec_reputation_profile_actions" "test" {
  config_id             = 43253
  security_policy_id    = "AAAA_81230"
  reputation_profile_id = 321456
}

