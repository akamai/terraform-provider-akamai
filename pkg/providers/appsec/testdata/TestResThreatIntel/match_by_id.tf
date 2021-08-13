provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_threat_intel" "test" {
  config_id           = 43253
  security_policy_id  = "AAAA_81230"
  threat_intel        = "off"
}
