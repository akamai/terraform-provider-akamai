provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_waf_attack_group_actions" "test" {
config_id = 43253
    version = 7
    policy_id = "AAAA_81230"
    group_id = "SQL"
}



