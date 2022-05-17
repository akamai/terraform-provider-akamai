provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_security_policy" "test" {
  config_id              = 43253
  security_policy_name   = "PLE Cloned Test for Launchpad 15"
  security_policy_prefix = "PLE"
}

