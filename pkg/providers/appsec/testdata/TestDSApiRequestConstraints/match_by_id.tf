provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_api_request_constraints" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  api_id             = 1
}
