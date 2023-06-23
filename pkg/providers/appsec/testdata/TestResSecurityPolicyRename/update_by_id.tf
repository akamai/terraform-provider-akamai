provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_security_policy_rename" "test" {
  config_id            = 43253
  security_policy_name = "Cloned Test for Launchpad 15"
  security_policy_id   = "PLE_114049"
}

