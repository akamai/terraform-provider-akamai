provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_appsec_clientside_protection_configuration" "test" {
  security_policy_id = "test_policy_id"
}