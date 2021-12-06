provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_eval_rule_actions" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
}

output "configsedge_post_output_text" {
  value = data.akamai_appsec_eval_rule_actions.test.output_text
}
