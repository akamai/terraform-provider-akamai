provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_rule_upgrade" "test" {
 config_id = 43253
  version = 7
 security_policy_id = "AAAA_81230"
}

