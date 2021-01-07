provider "akamai" {
  edgerc = "~/.edgerc"
}



data "akamai_appsec_reputation_analysis" "reputation_analysis" {
  config_id = 43253
  version = 7
  security_policy_id = "AAAA_81230"
}