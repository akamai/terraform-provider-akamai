provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_appsec_eval_rule_action" "test" {
    config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
  eval_rule_id = 699989
  rule_action = "none"
}

