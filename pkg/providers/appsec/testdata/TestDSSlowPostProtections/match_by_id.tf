provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_slowpost_protections" "test" {
  config_id          = 43253
  version            = 7
  security_policy_id = "AAAA_81230"
}

