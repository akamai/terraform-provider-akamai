provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_security_policy_protections" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
}

output "appsecwafmode" {
  value = data.akamai_appsec_security_policy_protections.test.output_text
}

