provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_eval_protect_host" "protect_host" {
 config_id = 43253
    version = 7
  hostnames = ["example.com"]
}
resource "akamai_appsec_attack_group_action" "test" {
config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
  api_endpoint_id = data.akamai_appsec_api_endpoints.api_endpoint.id
  action = "alert"
}


