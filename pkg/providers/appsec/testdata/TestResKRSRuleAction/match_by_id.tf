provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_appsec_krs_rule_action" "test" {
    config_id = 43253
    version = 7
    policy_id = "AAAA_81230"
    rule_id = 699989
    krs_rule_action = "alert"
}

