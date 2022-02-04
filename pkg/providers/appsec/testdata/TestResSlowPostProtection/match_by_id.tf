provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

resource "akamai_appsec_slowpost_protection" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  enabled            = false
}

