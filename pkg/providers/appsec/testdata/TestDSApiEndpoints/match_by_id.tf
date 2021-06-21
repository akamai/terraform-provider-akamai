provider "akamai" {
  edgerc = "~/.edgerc"
}
data "akamai_appsec_api_endpoints" "test" {
  config_id = 43253
  // api_name = var.api_endpoint_name
}
