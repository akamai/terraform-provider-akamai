provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_rate_policy_actions" "test" {
    config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
    
}

output "ds_rate_policy_actions" {
  value = data.akamai_appsec_rate_policy_actions.test.output_text
}

