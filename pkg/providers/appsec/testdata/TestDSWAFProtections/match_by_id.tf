provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_waf_protection" "test" {
  config_id          = 43253
  version            = 7
  security_policy_id = "AAAA_81230"
}

output "appsecwafmode" {
  value = data.akamai_appsec_waf_protection.test.output_text
}

