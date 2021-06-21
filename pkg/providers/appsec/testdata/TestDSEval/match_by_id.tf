provider "akamai" {
  edgerc = "~/.edgerc"
}



data "akamai_appsec_eval" "test" {
    config_id = 43253
    security_policy_id = "AAAA_81230"
}


