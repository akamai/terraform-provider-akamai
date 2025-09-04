provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

data "akamai_apr_protected_operations" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
}