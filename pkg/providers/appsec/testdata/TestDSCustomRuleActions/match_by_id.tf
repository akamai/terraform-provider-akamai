provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_custom_rule_actions" "test" {
    config_id = 43253
    version = 7
    policy_id = "AAAA_81230"
}

output "customruleactions" {
  value = data.akamai_appsec_custom_rule_actions.test.output_text
}

