provider "akamai" {
  edgerc = "~/.edgerc"
}



resource "akamai_appsec_rate_policy_action" "test" {
    config_id = 43253
    version = 15
    policy_id = "AAAA_81230"
    rate_policy_id = 135355
    ipv4_action = "none"
    ipv6_action = "none"
}


