provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_policy_api_endpoints" "api_endpoint" {
  config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
  name = "TestEndpoint"
}

