provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_security_policy_protections" "test" {
  config_id          = 12345
  security_policy_id = "test_policy"
}
