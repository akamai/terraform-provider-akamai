provider "akamai" {
  edgerc = "~/.edgerc"
}


data "akamai_appsec_aag_rules" "test" {
    config_id = 43253
    version = 7
    policy_id = "AAAA_81230"
    group_id = "SQL"
}

output "configsedge_post_output_text" {
  value = data.akamai_appsec_aag_rules.test.output_text
}
