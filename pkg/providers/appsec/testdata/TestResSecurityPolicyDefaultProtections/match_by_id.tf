provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_security_policy_default_protections" "test" {
  config_id              = 43253
  security_policy_name   = "PLEB Cloned Test for Launchpad 15"
  security_policy_prefix = "PLEB"
}
