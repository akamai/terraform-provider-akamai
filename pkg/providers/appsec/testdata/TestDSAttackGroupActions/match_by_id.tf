provider "akamai" {
  edgerc = "~/.edgerc"
  //version = "1.0.1"
}

data "akamai_appsec_attack_group_actions" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  attack_group       = "SQL"
}



