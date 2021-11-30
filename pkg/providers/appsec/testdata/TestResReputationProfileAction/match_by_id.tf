provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_reputation_profile_action" "test" {
  config_id             = 43253
  security_policy_id    = "AAAA_81230"
  reputation_profile_id = 1685099
  action                = "none"
}


