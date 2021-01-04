provider "akamai" {
  edgerc = "~/.edgerc"
}



data "akamai_appsec_reputation_profile_actions" "test" {
    config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
    reputation_profile_id = 321456
}


