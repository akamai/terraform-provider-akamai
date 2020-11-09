provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_reputation_protections" "test" {
    config_id = 43253
    version = 7
    policy_id = "AAAA_81230"
}

output "appsecwafmode" {
  value = data.akamai_appsec_reputation_protections.test.output_text
}

