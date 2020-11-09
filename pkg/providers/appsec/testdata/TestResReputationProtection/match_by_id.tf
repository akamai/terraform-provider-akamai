provider "akamai" {
  edgerc = "~/.edgerc"
}



resource "akamai_appsec_reputation_protection" "test" {
    config_id = 43253
    version = 7
    policy_id = "AAAA_81230"
    enabled = false
}

output "appsecrateprotection" {
  value = akamai_appsec_reputation_protection.test.output_text
}

