provider "akamai" {
  edgerc = "~/.edgerc"
}



resource "akamai_appsec_slowpost_protection" "test" {
    config_id = 43253
    version = 7
    security_policy_id = "AAAA_81230"
    enabled = false
}


