provider "akamai" {
  edgerc = "~/.edgerc"
}




data "akamai_appsec_api_endpoints" "test" {
  config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
 // name = var.api_endpoint_name
}