provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_attack_group_action" "test" {
  config_id           = 43253
  security_policy_id  = "AAAA_81230"
  attack_group        = "SQL"
  attack_group_action = "none"
}


