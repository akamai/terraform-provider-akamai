provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_appsec_waf_attack_group_action" "test" {
    config_id = 43253
    version = 7
    policy_id = "AAAA_81230"
    group_id = "SQL"
    action = "alert"
}



