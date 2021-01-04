provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_api_request_constraints" "api_request_constraints" {
  config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
  api_id = var.api_id
}
