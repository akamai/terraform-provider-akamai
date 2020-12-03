provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_eval_rule_condition_exception" "test" {
    config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
    rule_id = 12345
   
}


