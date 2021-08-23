provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_rate_protection" "test" {
    config_id = 43253
    security_policy_id = "AAAA_81230"
    enabled = true
}



