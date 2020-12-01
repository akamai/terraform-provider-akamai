provider "akamai" {
  edgerc = "~/.edgerc"
}



resource "akamai_appsec_security_policy_protections" "test" {
    config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
   apply_application_layer_controls = false
  apply_network_layer_controls = false
  apply_rate_controls = false
  apply_reputation_controls = false
  apply_botman_controls = false
  apply_api_constraints = false
  apply_slow_post_controls = false
}

output "appsecwafprotection" {
  value = akamai_appsec_security_policy_protections.test.output_text
}

