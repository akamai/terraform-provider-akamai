provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

data "akamai_appsec_reputation_profile_analysis" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
}

